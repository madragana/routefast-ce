package lip4d

import (
    "crypto/ed25519"
    "encoding/base64"
    "errors"
    "time"

    rfcrypto "github.com/aovidi/routefast-linux/internal/crypto"
)

func VerifyAgainstKeyString(env Envelope, pubKey string) error {
    pub, err := rfcrypto.ParsePublicKey(pubKey)
    if err != nil { return err }
    b, err := CanonicalBytes(env)
    if err != nil { return err }
    sig, err := base64.RawURLEncoding.DecodeString(env.Signature)
    if err != nil { return err }
    if !ed25519.Verify(pub, b, sig) { return errors.New("invalid signature") }
    return nil
}
func ValidateTimeWindow(tw TimeWindow, maxSkew time.Duration) error {
    created, err := time.Parse(time.RFC3339, tw.CreatedAt)
    if err != nil { return err }
    notAfter, err := time.Parse(time.RFC3339, tw.NotAfter)
    if err != nil { return err }
    now := time.Now().UTC()
    if created.After(now.Add(maxSkew)) { return errors.New("created_at too far in future") }
    if now.After(notAfter) { return errors.New("message expired") }
    return nil
}
