package policy

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	NodeID              string
	AllowedActions      []string
	HighImpactActions   []string
	MaxImpactedPrefixes int
	MaxTTLSeconds       int
	RequireQuorum       bool
}

func Load(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	cfg := Config{}
	var currentList string
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "- ") {
			item := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			switch currentList {
			case "allowed_actions":
				cfg.AllowedActions = append(cfg.AllowedActions, item)
			case "high_impact_actions":
				cfg.HighImpactActions = append(cfg.HighImpactActions, item)
			}
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		currentList = ""
		switch key {
		case "node_id":
			cfg.NodeID = val
		case "allowed_actions", "high_impact_actions":
			currentList = key
		case "max_impacted_prefixes":
			cfg.MaxImpactedPrefixes, _ = strconv.Atoi(val)
		case "max_ttl_seconds":
			cfg.MaxTTLSeconds, _ = strconv.Atoi(val)
		case "require_quorum":
			cfg.RequireQuorum = strings.EqualFold(val, "true")
		}
	}
	if err := s.Err(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func Authorize(cfg Config, action string, ttl int) error {
	allowed := false
	for _, a := range cfg.AllowedActions {
		if a == action {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("action %s not allowed", action)
	}
	if ttl > cfg.MaxTTLSeconds {
		return fmt.Errorf("ttl %d exceeds max %d", ttl, cfg.MaxTTLSeconds)
	}
	return nil
}
