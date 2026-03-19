package lip4d

import "time"

type Envelope struct {
	LIP4DVersion string         `json:"lip4d_version"`
	MessageType  string         `json:"message_type"`
	IntentID     string         `json:"intent_id"`
	Sender       Sender         `json:"sender"`
	Purpose      Purpose        `json:"purpose"`
	Context      map[string]any `json:"context"`
	Reason       Reason         `json:"reason"`
	Authorization Authorization `json:"authorization"`
	Time         Timing         `json:"time"`
	Attestation  map[string]any `json:"attestation,omitempty"`
	Signature    string         `json:"signature"`
}

type Sender struct {
	NodeID string `json:"node_id"`
	Domain string `json:"domain"`
	KeyID  string `json:"key_id"`
}

type Purpose struct {
	Verb   string `json:"verb"`
	Object string `json:"object"`
	Target string `json:"target"`
}

type Reason struct {
	Summary        string   `json:"summary"`
	Evidence       []string `json:"evidence"`
	EvidenceDigest string   `json:"evidence_digest"`
}

type Authorization struct {
	Role                string `json:"role"`
	ROEProfile          string `json:"roe_profile"`
	RequestedCapability string `json:"requested_capability"`
}

type Timing struct {
	CreatedAt  time.Time `json:"created_at"`
	NotAfter   time.Time `json:"not_after"`
	TTLSeconds int       `json:"ttl_seconds"`
	Sequence   int64     `json:"sequence"`
}
