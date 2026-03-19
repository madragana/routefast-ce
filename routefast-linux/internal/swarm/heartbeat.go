package swarm

import (
    "context"
    "crypto/ed25519"
    "encoding/json"
    "log"
    "time"
    "github.com/aovidi/routefast-linux/internal/lip4d"
)
type HeartbeatSender struct { NodeID, Role, KeyID string; Priv ed25519.PrivateKey; Peers PeersConfig }
func NewHeartbeatSender(nodeID, role, keyID string, priv ed25519.PrivateKey, peers PeersConfig) *HeartbeatSender { return &HeartbeatSender{NodeID: nodeID, Role: role, KeyID: keyID, Priv: priv, Peers: peers} }
func (h *HeartbeatSender) Start(ctx context.Context, interval time.Duration) {
    client := lip4d.NewClient(h.NodeID, h.Role, h.KeyID, h.Priv); t := time.NewTicker(interval); defer t.Stop()
    sendOnce := func() { for _, p := range h.Peers.Peers { env := lip4d.Envelope{LIP4DVersion:"0.5", MessageType:"HEARTBEAT", MessageID:"hb-"+h.NodeID+"-"+time.Now().UTC().Format("20060102150405.000000000"), Body:map[string]any{"health":"ok"}, Time: lip4d.TimeWindow{CreatedAt: time.Now().UTC().Format(time.RFC3339), NotAfter: time.Now().UTC().Add(15*time.Second).Format(time.RFC3339), Sequence: time.Now().UnixNano()}}; resp, err := client.Post(ctx, p.Address+"/v1/heartbeat", env); if err != nil { log.Printf("heartbeat to %s failed: %v", p.NodeID, err); continue }; var out map[string]any; _ = json.NewDecoder(resp.Body).Decode(&out); _ = resp.Body.Close() } }
    sendOnce(); for { select { case <-ctx.Done(): return; case <-t.C: sendOnce() } }
}
