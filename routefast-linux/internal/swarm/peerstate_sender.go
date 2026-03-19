package swarm

import (
    "context"
    "crypto/ed25519"
    "encoding/json"
    "log"
    "time"
    "github.com/aovidi/routefast-linux/internal/lip4d"
)
type PeerStateSender struct { NodeID, Role, KeyID string; Priv ed25519.PrivateKey; Peers PeersConfig }
func NewPeerStateSender(nodeID, role, keyID string, priv ed25519.PrivateKey, peers PeersConfig) *PeerStateSender { return &PeerStateSender{NodeID: nodeID, Role: role, KeyID: keyID, Priv: priv, Peers: peers} }
func (p *PeerStateSender) Start(ctx context.Context, interval time.Duration) {
    client := lip4d.NewClient(p.NodeID, p.Role, p.KeyID, p.Priv); t := time.NewTicker(interval); defer t.Stop()
    sendOnce := func() { for _, peer := range p.Peers.Peers { env := lip4d.Envelope{LIP4DVersion:"0.5", MessageType:"GOSSIP_STATE", MessageID:"ps-"+p.NodeID+"-"+time.Now().UTC().Format("20060102150405.000000000"), Body: map[string]any{"node_id": p.NodeID, "timestamp": time.Now().UTC().Format(time.RFC3339), "attack_score": 0.2, "confidence": 0.88, "active_intents": 0, "bearers": []map[string]any{{"name":"fiber0","health":0.91},{"name":"lte0","health":0.72},{"name":"radio0","health":0.95}}}, Time: lip4d.TimeWindow{CreatedAt: time.Now().UTC().Format(time.RFC3339), NotAfter: time.Now().UTC().Add(20*time.Second).Format(time.RFC3339), Sequence: time.Now().UnixNano()}}; resp, err := client.Post(ctx, peer.Address+"/v1/peerstate", env); if err != nil { log.Printf("peerstate to %s failed: %v", peer.NodeID, err); continue }; var out map[string]any; _ = json.NewDecoder(resp.Body).Decode(&out); _ = resp.Body.Close() } }
    sendOnce(); for { select { case <-ctx.Done(): return; case <-t.C: sendOnce() } }
}
