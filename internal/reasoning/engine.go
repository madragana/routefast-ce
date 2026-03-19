package reasoning

import (
	"context"
	"time"

	"github.com/aovidi/routefast-linux/internal/telemetry"
)

type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

func (e *Engine) Decide(_ context.Context, batch []telemetry.Event) (Decision, error) {
	decision := Decision{
		Action:                    "EXPLAIN_ONLY",
		Confidence:                0.40,
		Reason:                    "insufficient evidence for mitigation",
		Evidence:                  []string{},
		ImpactEstimate:            "low",
		RollbackTTLSeconds:        0,
		PeerCorroborationRequired: false,
		CreatedAt:                 time.Now().UTC(),
	}

	var sawAttack bool
	var target string
	for _, ev := range batch {
		decision.Evidence = append(decision.Evidence, ev.Source+":"+ev.Type)
		if ev.Type == "attack_suspected" {
			sawAttack = true
			if ip, ok := ev.Attrs["src_ip"].(string); ok {
				target = ip
			}
		}
		if ev.Type == "bearer_failover_candidate" {
			decision = Decision{
				Action:             "SHIFT_SERVICE_CLASS_TO_BEARER",
				Target:             "critical-control",
				Constraints:        map[string]interface{}{"prefer_media": []string{"radio", "lte"}},
				Confidence:         0.82,
				Reason:             "primary bearer degraded and fallback bearer probe succeeded",
				Evidence:           decision.Evidence,
				ImpactEstimate:     "medium",
				RollbackTTLSeconds: 600,
				CreatedAt:          time.Now().UTC(),
			}
		}
	}

	if sawAttack {
		decision = Decision{
			Action:                    "BLOCK_SOURCE",
			Target:                    target,
			Confidence:                0.91,
			Reason:                    "correlated attack telemetry from multiple sources",
			Evidence:                  decision.Evidence,
			ImpactEstimate:            "medium",
			RollbackTTLSeconds:        300,
			PeerCorroborationRequired: true,
			CreatedAt:                 time.Now().UTC(),
		}
	}

	return decision, nil
}
