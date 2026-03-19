package dashboard

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/aovidi/routefast-linux/internal/swarm"
)
type PeerPoller struct { HTTP *http.Client; Peers swarm.PeersConfig; State *State }
func NewPeerPoller(peers swarm.PeersConfig, state *State) *PeerPoller { return &PeerPoller{HTTP: &http.Client{Timeout: 2*time.Second}, Peers: peers, State: state} }
func (p *PeerPoller) Start(ctx context.Context, interval time.Duration) { t := time.NewTicker(interval); defer t.Stop(); _ = p.PollPeersNow(); for { select { case <-ctx.Done(): return; case <-t.C: _ = p.PollPeersNow() } } }
func (p *PeerPoller) PollPeersNow() error {
    failed := 0
    for _, peer := range p.Peers.Peers {
        resp, err := p.HTTP.Get(peer.Address + "/v1/status")
        if err != nil { p.State.AddEvent("poll_error", peer.NodeID+": "+err.Error()); failed++; continue }
        var out map[string]any
        if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { _ = resp.Body.Close(); failed++; continue }
        _ = resp.Body.Close()
        p.State.UpdatePeerHeartbeat(peer.NodeID, time.Now().UTC())
        p.State.SetRemotePeerStatus(peer.NodeID, out)
    }
    if failed == len(p.Peers.Peers) && len(p.Peers.Peers) > 0 { return fmt.Errorf("all peer polls failed") }
    return nil
}
