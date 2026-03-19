package lip4d

import (
    "crypto/ed25519"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "fmt"
)

func SeedToPrivateKey(seed string) ed25519.PrivateKey {
    sum := sha256.Sum256([]byte(seed))
    return ed25519.NewKeyFromSeed(sum[:])
}

func PublicKeyString(pub ed25519.PublicKey) string {
    return base64.StdEncoding.EncodeToString(pub)
}

func EvidenceDigest(evidence []string) string {
    buf, _ := json.Marshal(evidence)
    sum := sha256.Sum256(buf)
    return "sha256:" + hex.EncodeToString(sum[:])
}

func SignEnvelope(e *Envelope, priv ed25519.PrivateKey) error {
    payload, err := canonicalPayload(*e)
    if err != nil { return err }
    e.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(priv, payload))
    return nil
}

func VerifyEnvelope(e Envelope, pub ed25519.PublicKey) error {
    sig, err := base64.StdEncoding.DecodeString(e.Signature)
    if err != nil { return err }
    payload, err := canonicalPayload(e)
    if err != nil { return err }
    if !ed25519.Verify(pub, payload, sig) { return fmt.Errorf("invalid signature") }
    return nil
}

func canonicalPayload(e Envelope) ([]byte, error) {
    e.Signature = ""
    return json.Marshal(e)
}
