package actuation

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "strings"

    "github.com/aovidi/routefast-linux/internal/reasoning"
)

type Result struct {
    Applied     bool                   `json:"applied"`
    Mode        string                 `json:"mode"`
    Summary     string                 `json:"summary"`
    Commands    []string               `json:"commands"`
    RollbackETA int                    `json:"rollback_eta"`
    Verification map[string]any        `json:"verification,omitempty"`
}

type Executor struct { DryRun bool }

func NewExecutor(dry bool) *Executor { return &Executor{DryRun: dry} }

func (e *Executor) Plan(d reasoning.Decision) ([]string, error) {
    switch d.Action {
    case "EXPLAIN_ONLY":
        return []string{"echo explain-only"}, nil
    case "BLOCK_SOURCE":
        if d.Target == "" { return nil, fmt.Errorf("missing target") }
        return []string{fmt.Sprintf("nft add rule inet routefast input ip saddr %s drop comment "routefast ttl=%d"", d.Target, d.RollbackTTLSeconds)}, nil
    case "SHIFT_SERVICE_CLASS_TO_BEARER", "REROUTE_PREFIX", "LOWER_LOCAL_PREF":
        cmds, err := FRRCommands(d)
        if err != nil { return nil, err }
        if len(cmds) > 0 { return cmds, nil }
        return []string{"echo no-op"}, nil
    default:
        return nil, fmt.Errorf("unsupported action %s", d.Action)
    }
}

func (e *Executor) Apply(ctx context.Context, d reasoning.Decision) (Result, error) {
    cmds, err := e.Plan(d)
    if err != nil { return Result{}, err }
    res := Result{Applied: true, RollbackETA: d.RollbackTTLSeconds, Commands: cmds, Verification: VerifyPlan(d, cmds)}
    live := os.Getenv("ROUTEFAST_EXECUTE") == "1"
    if e.DryRun || !live {
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
