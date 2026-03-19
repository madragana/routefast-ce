package swarm

import "time"

type Peer struct {
    NodeID string `yaml:"node_id" json:"node_id"`
    Address string `yaml:"address" json:"address"`
    Role string `yaml:"role" json:"role"`
    PubKey string `yaml:"pubkey" json:"pubkey"`
    Region string `yaml:"region" json:"region"`
    LastHeartbeat time.Time `json:"last_heartbeat"`
}
type PeersConfig struct { Peers []Peer `yaml:"peers"` }
func (c PeersConfig) ByNodeID(id string) *Peer {
    for i := range c.Peers { if c.Peers[i].NodeID == id { return &c.Peers[i] } }
    return nil
}
