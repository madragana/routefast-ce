package policy

import (
    "fmt"
    "slices"

    "github.com/aovidi/routefast-linux/internal/config"
)

func Authorize(cfg config.Policy, action string, ttl int) error {
    if !slices.Contains(cfg.AllowedActions, action) {
        return fmt.Errorf("action %s not allowed", action)
    }
    if ttl < 0 || ttl > cfg.MaxTTLSeconds {
        return fmt.Errorf("ttl %d exceeds max %d", ttl, cfg.MaxTTLSeconds)
    }
    return nil
}

func IsHighImpact(cfg config.Policy, action string) bool {
    return slices.Contains(cfg.HighImpactActions, action)
}
