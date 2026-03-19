package lip4d

import (
    "crypto/ed25519"
    "encoding/base64"
    "encoding/json"
)

type signableEnvelope struct {
    LIP4DVersion string `json:"lip4d_version"`
    MessageType string `json:"message_type"`
    MessageID string `json:"message_id"`
    IntentID string `json:"intent_id,omitempty"`
    Sender Sender `json:"sender"`
    Body interface{} `json:"body"`
    Time TimeWindow `json:"time"`
}
func CanonicalBytes(env Envelope) ([]byte, error) {
    return json.Marshal(signableEnvelope{LIP4DVersion: env.LIP4DVersion, MessageType: env.MessageType, MessageID: env.MessageID, IntentID: env.IntentID, Sender: env.Sender, Body: env.Body, Time: env.Time})
}
func Sign(env *Envelope, priv ed25519.PrivateKey) error {
    b, err := CanonicalBytes(*env)
    if err != nil { return err }
    env.Signature = base64.RawURLEncoding.EncodeToString(ed25519.Sign(priv, b))
    return nil
}
