package reasoning

import "time"

type Decision struct {
	Action                    string                 `json:"action"`
	Target                    string                 `json:"target,omitempty"`
	Constraints               map[string]interface{} `json:"constraints,omitempty"`
	Confidence                float64                `json:"confidence"`
	Reason                    string                 `json:"reason"`
	Evidence                  []string               `json:"evidence"`
	ImpactEstimate            string                 `json:"impact_estimate,omitempty"`
	RollbackTTLSeconds        int                    `json:"rollback_ttl_seconds"`
	PeerCorroborationRequired bool                   `json:"peer_corroboration_required"`
	CreatedAt                 time.Time              `json:"created_at"`
}
