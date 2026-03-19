package lip4d

import (
    "errors"
    "fmt"
    "time"
)

func ValidateEnvelope(e Envelope) error {
    if e.LIP4DVersion == "" || e.MessageType == "" || e.IntentID == "" { return errors.New("missing required lip4d headers") }
    if e.Sender.NodeID == "" || e.Sender.KeyID == "" { return errors.New("missing sender identity") }
    if e.Authorization.RequestedCapability == "" { return errors.New("missing requested capability") }
    if e.Reason.EvidenceDigest == "" { return errors.New("missing evidence digest") }
    if e.Time.TTLSeconds <= 0 || e.Time.TTLSeconds > 3600 { return fmt.Errorf("invalid ttl: %d", e.Time.TTLSeconds) }
    now := time.Now().UTC()
    if e.Time.NotAfter.Before(now) { return errors.New("intent expired") }
    if e.Signature == "" { return errors.New("missing signature") }
    return nil
}
