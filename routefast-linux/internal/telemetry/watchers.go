package telemetry

import (
    "bufio"
    "context"
    "encoding/json"
    "os"
    "os/exec"
    "strings"
    "time"
)

func WatchJSONLFile(ctx context.Context, source, path, eventType string, out chan<- Event) error {
    f, err := os.Open(path)
    if err != nil { return err }
    defer f.Close()
    s := bufio.NewScanner(f)
    for s.Scan() {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        line := strings.TrimSpace(s.Text())
        if line == "" { continue }
        attrs := map[string]any{"raw": line}
        _ = json.Unmarshal([]byte(line), &attrs)
        out <- Event{Source: source, Type: eventType, Timestamp: time.Now().UTC(), Attrs: attrs}
    }
    return s.Err()
}

func WatchSuricataEVE(ctx context.Context, path string, out chan<- Event) error {
    return WatchJSONLFile(ctx, "suricata", path, "suricata_alert", out)
}

func WatchEBPFExporter(ctx context.Context, path string, out chan<- Event) error {
    return WatchJSONLFile(ctx, "ebpf", path, "pps_anomaly", out)
}

func WatchFRRLog(ctx context.Context, path string, out chan<- Event) error {
    f, err := os.Open(path)
    if err != nil { return err }
    defer f.Close()
    s := bufio.NewScanner(f)
    for s.Scan() {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        line := s.Text()
        evType := "frr_log"
        if strings.Contains(line, "withdraw") || strings.Contains(line, "Withdraw") {
            evType = "bgp_withdraw"
        }
        if strings.Contains(strings.ToLower(line), "down") {
            evType = "link_down"
        }
        out <- Event{Source: "frr", Type: evType, Timestamp: time.Now().UTC(), Attrs: map[string]any{"line": line}}
    }
    return s.Err()
}

func PollKernelLinks(ctx context.Context, interval time.Duration, out chan<- Event) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            cmd := exec.CommandContext(ctx, "ip", "-json", "link", "show")
            raw, err := cmd.Output()
            attrs := map[string]any{"raw": string(raw)}
            evType := "link_state"
            if err != nil {
                attrs["error"] = err.Error()
                evType = "link_poll_error"
            }
            out <- Event{Source: "kernel", Type: evType, Timestamp: time.Now().UTC(), Attrs: attrs}
        }
    }
}
