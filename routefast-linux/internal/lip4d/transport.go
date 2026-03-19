package lip4d

import (
    "bytes"
    "context"
    "crypto/ed25519"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    HTTP *http.Client
    NodeID string
    Role string
    KeyID string
    PrivKey ed25519.PrivateKey
}
func NewClient(nodeID, role, keyID string, priv ed25519.PrivateKey) *Client {
    return &Client{HTTP: &http.Client{Timeout: 3*time.Second}, NodeID: nodeID, Role: role, KeyID: keyID, PrivKey: priv}
}
func (c *Client) Post(ctx context.Context, url string, env Envelope) (*http.Response, error) {
    env.Sender = Sender{NodeID: c.NodeID, Role: c.Role, KeyID: c.KeyID}
    if err := Sign(&env, c.PrivKey); err != nil { return nil, err }
    b, _ := json.Marshal(env)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
    if err != nil { return nil, err }
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.HTTP.Do(req)
    if err != nil { return nil, err }
    if resp.StatusCode >= 300 { return resp, fmt.Errorf("peer returned status %s", resp.Status) }
    return resp, nil
}
