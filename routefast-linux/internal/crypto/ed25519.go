package crypto

import (
    "crypto/ed25519"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "os"
    "strings"
)

func GenerateKeypair() (ed25519.PublicKey, ed25519.PrivateKey, error) { return ed25519.GenerateKey(rand.Reader) }
func LoadPrivateKey(path string) (ed25519.PrivateKey, error) {
    b, err := os.ReadFile(path)
    if err != nil { return nil, err }
    raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(b)))
    if err != nil { return nil, err }
    if len(raw) != ed25519.PrivateKeySize { return nil, errors.New("invalid private key size") }
    return ed25519.PrivateKey(raw), nil
}
func ParsePublicKey(s string) (ed25519.PublicKey, error) {
    s = strings.TrimPrefix(strings.TrimSpace(s), "ed25519:")
    raw, err := base64.StdEncoding.DecodeString(s)
    if err != nil { return nil, err }
    if len(raw) != ed25519.PublicKeySize { return nil, errors.New("invalid public key size") }
    return ed25519.PublicKey(raw), nil
}
