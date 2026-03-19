package explain

import (
    "encoding/json"
    "os"
    "sync"
    "time"

    "github.com/aovidi/routefast-linux/internal/reasoning"
)

type Audit struct {
    path string
    mu   sync.Mutex
}

type Record struct {
    Time     time.Time          `json:"time"`
    Decision reasoning.Decision `json:"decision"`
    Status   string             `json:"status"`
    Detail   string             `json:"detail"`
}

func NewAudit(path string) *Audit { return &Audit{path: path} }

func (a *Audit) Write(d reasoning.Decision, status, detail string) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    f, err := os.OpenFile(a.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
    if err != nil { return err }
    defer f.Close()
    rec := Record{Time: time.Now().UTC(), Decision: d, Status: status, Detail: detail}
    enc := json.NewEncoder(f)
    return enc.Encode(rec)
}
