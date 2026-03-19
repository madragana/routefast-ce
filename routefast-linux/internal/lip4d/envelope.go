package lip4d

type Sender struct {
    NodeID string `json:"node_id"`
    Role string `json:"role,omitempty"`
    KeyID string `json:"key_id,omitempty"`
}
type TimeWindow struct {
    CreatedAt string `json:"created_at"`
    NotAfter string `json:"not_after"`
    Sequence int64 `json:"sequence"`
}
type Envelope struct {
    LIP4DVersion string `json:"lip4d_version"`
    MessageType string `json:"message_type"`
    MessageID string `json:"message_id"`
    IntentID string `json:"intent_id,omitempty"`
    Sender Sender `json:"sender"`
    Body interface{} `json:"body"`
    Time TimeWindow `json:"time"`
    Signature string `json:"signature"`
}
