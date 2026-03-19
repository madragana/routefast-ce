package actuation

import "github.com/aovidi/routefast-linux/internal/reasoning"

func VerifyPlan(d reasoning.Decision, commands []string) map[string]any {
    return map[string]any{
        "action": d.Action,
        "command_count": len(commands),
        "target": d.Target,
        "rollback_ttl_seconds": d.RollbackTTLSeconds,
    }
}
