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

type Sender struct { NodeID, Domain, KeyID string }
type Purpose struct { Verb, Object, Target string }
type Reason struct { Summary string; Evidence []string; EvidenceDigest string }
type Authorization struct { Role, ROEProfile, RequestedCapability string }
type Timing struct { CreatedAt, NotAfter time.Time; TTLSeconds int; Sequence int64 }

type ValidationResponse struct {
    Accepted bool   `json:"accepted"`
    Detail   string `json:"detail"`
}
