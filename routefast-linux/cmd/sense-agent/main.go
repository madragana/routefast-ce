package main

import (
    "bytes"
    "context"
    "encoding/json"
    "flag"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/aovidi/routefast-linux/internal/telemetry"
)

func postEvent(url string, ev telemetry.Event) error {
    buf, _ := json.Marshal(ev)
    resp, err := http.Post(url, "application/json", bytes.NewReader(buf))
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 { return nil }
    return nil
}

func main() {
    mode := flag.String("mode", "watch", "watch or demo")
    routefastURL := flag.String("routefast-url", "http://127.0.0.1:8090/event", "routefast event endpoint")
    suricataPath := flag.String("suricata", "/var/log/suricata/eve.json", "suricata eve path")
    frrPath := flag.String("frr", "/var/log/frr/frr.log", "frr log path")
    ebpfPath := flag.String("ebpf", "./ebpf-export.jsonl", "ebpf jsonl path")
    flag.Parse()

    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    events := make(chan telemetry.Event, 256)
    go func() { for ev := range events { _ = postEvent(*routefastURL, ev) } }()

    if *mode == "demo" {
        demo := []telemetry.Event{{Source:"suricata",Type:"suricata_alert",Timestamp:time.Now().UTC(),Attrs:map[string]any{"src_ip":"198.51.100.44"}}, {Source:"ebpf",Type:"pps_anomaly",Timestamp:time.Now().UTC(),Attrs:map[string]any{"src_ip":"198.51.100.44","pps":420000}}, {Source:"frr",Type:"bgp_withdraw",Timestamp:time.Now().UTC(),Attrs:map[string]any{"peer":"peer-3"}}}
        for _, ev := range demo { _ = postEvent(*routefastURL, ev) }
        log.Println("demo events sent")
        return
    }

    go telemetry.PollKernelLinks(ctx, 10*time.Second, events)
    if _, err := os.Stat(*suricataPath); err == nil { go func() { _ = telemetry.WatchSuricataEVE(ctx, *suricataPath, events) }() }
    if _, err := os.Stat(*frrPath); err == nil { go func() { _ = telemetry.WatchFRRLog(ctx, *frrPath, events) }() }
    if _, err := os.Stat(*ebpfPath); err == nil { go func() { _ = telemetry.WatchEBPFExporter(ctx, *ebpfPath, events) }() }

    <-ctx.Done()
    log.Println("sense-agent stopping")
}
