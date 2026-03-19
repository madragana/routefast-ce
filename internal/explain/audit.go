package explain

import (
	"encoding/json"
	"os"
	"time"

	"github.com/aovidi/routefast-linux/internal/reasoning"
)

type Record struct {
	Timestamp   time.Time          `json:"timestamp"`
	Decision    reasoning.Decision `json:"decision"`
	Policy      string             `json:"policy"`
	Result      string             `json:"result"`
}

type Audit struct {
	path string
}

func NewAudit(path string) *Audit { return &Audit{path: path} }

func (a *Audit) Write(decision reasoning.Decision, policyResult, result string) error {
	rec := Record{Timestamp: time.Now().UTC(), Decision: decision, Policy: policyResult, Result: result}
	f, err := os.OpenFile(a.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(rec)
}
