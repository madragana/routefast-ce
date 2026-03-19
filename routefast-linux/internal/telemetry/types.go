package telemetry

import "time"

type Event struct {
    Source    string                 `json:"source"`
    Type      string                 `json:"type"`
    Timestamp time.Time              `json:"timestamp"`
    Attrs     map[string]interface{} `json:"attrs"`
}

type EventBatch struct {
    Events []Event `json:"events"`
}
