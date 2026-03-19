package actuation

import (
    "context"
    "fmt"
    "os/exec"
    "strings"

    "github.com/aovidi/routefast-linux/internal/reasoning"
)

type Result struct {
    Applied     bool     `json:"applied"`
    Mode        string   `json:"mode"`
    Summary     string   `json:"summary"`
    Commands    []string `json:"commands"`
    RollbackETA int      `json:"rollback_eta"`
}

type Executor struct { DryRun bool }

func NewExecutor(dry bool) *Executor { return &Executor{DryRun: dry} }

func (e *Executor) Plan(d reasoning.Decision) ([]string, error) {
    switch d.Action {
    case "EXPLAIN_ONLY":
        return []string{"echo explain-only"}, nil
    case "BLOCK_SOURCE":
        if d.Target == "" { return nil, fmt.Errorf("missing target") }
        return []string{fmt.Sprintf("nft add rule inet routefast input ip saddr %s drop comment \"routefast ttl=%d\"", d.Target, d.RollbackTTLSeconds)}, nil
    case "SHIFT_SERVICE_CLASS_TO_BEARER":
        media, _ := d.Constraints["prefer_media"].([]interface{})
        var pref []string
        for _, m := range media { pref = append(pref, fmt.Sprint(m)) }
        if len(pref) == 0 {
            if xs, ok := d.Constraints["prefer_media"].([]string); ok { pref = xs }
        }
        chosen := "radio"
        if len(pref) > 0 { chosen = pref[0] }
        return []string{
            fmt.Sprintf("ip rule add fwmark 0x1 table 101"),
            fmt.Sprintf("ip route replace default dev %s0 table 101", chosen),
        }, nil
    default:
        return nil, fmt.Errorf("unsupported action %s", d.Action)
    }
}

func (e *Executor) Apply(ctx context.Context, d reasoning.Decision) (Result, error) {
    cmds, err := e.Plan(d)
    if err != nil { return Result{}, err }
    res := Result{Applied: true, RollbackETA: d.RollbackTTLSeconds, Commands: cmds}
    if e.DryRun {
        res.Mode = "dry-run"
        res.Summary = strings.Join(cmds, " && ")
        return res, nil
    }
    for _, c := range cmds {
        cmd := exec.CommandContext(ctx, "bash", "-lc", c)
        out, err := cmd.CombinedOutput()
        if err != nil { return Result{}, fmt.Errorf("%s: %v: %s", c, err, string(out)) }
    }
    res.Mode = "live"
    res.Summary = "commands executed"
    return res, nil
}
