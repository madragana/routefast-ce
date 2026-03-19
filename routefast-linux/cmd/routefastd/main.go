package main

import (
    "context"
    "crypto/ed25519"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"
    rfcrypto "github.com/aovidi/routefast-linux/internal/crypto"
    "github.com/aovidi/routefast-linux/internal/dashboard"
    "github.com/aovidi/routefast-linux/internal/lip4d"
    "github.com/aovidi/routefast-linux/internal/policy"
    "github.com/aovidi/routefast-linux/internal/swarm"
)
type nodeConfig struct { Node struct { ID, Domain, ListenAddr, Role, Region, PrivateKeyFile string `yaml:",omitempty"` } `yaml:"node"` }
func getenv(k, d string) string { v := os.Getenv(k); if v == "" { return d }; return v }
func mustLoadPrivateKey(path string) ed25519.PrivateKey { priv, err := rfcrypto.LoadPrivateKey(path); if err != nil { log.Fatal(err) }; return priv }
type App struct { Node nodeConfig; Peers swarm.PeersConfig; Quorum swarm.QuorumConfig; Client *lip4d.Client; State *dashboard.State; Poller *dashboard.PeerPoller }
func (a *App) PollPeersNow() error { return a.Poller.PollPeersNow() }
func (a *App) RunDemoCorroboration() (any, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second); defer cancel()
    intentID := "intent-demo-1"; peerConf := []float64{}; accepts, rejects, uncertain := 0,0,0
    for _, p := range a.Peers.Peers {
        env := lip4d.Envelope{LIP4DVersion:"0.5", MessageType:"REQUEST_CORROBORATION", MessageID:"msg-"+p.NodeID+"-"+time.Now().UTC().Format("150405"), IntentID:intentID, Body: map[string]any{"query": map[string]any{"subject":"uplink_disruption","target":"203.0.113.0/24","window_seconds":30}, "local_summary": map[string]any{"confidence":0.82,"evidence":[]string{"bgp_withdraw_peer3","loss_spike_fiber0"}}}, Time: lip4d.TimeWindow{CreatedAt: time.Now().UTC().Format(time.RFC3339), NotAfter: time.Now().UTC().Add(30*time.Second).Format(time.RFC3339), Sequence: time.Now().UnixNano()}}
        resp, err := a.Client.Post(ctx, p.Address+"/v1/corroborate", env); if err != nil { a.State.AddEvent("peer_error", p.NodeID+": "+err.Error()); continue }
        var out map[string]any; _ = json.NewDecoder(resp.Body).Decode(&out); _ = resp.Body.Close(); a.State.UpdatePeerHeartbeat(p.NodeID, time.Now().UTC())
        verdict, _ := out["verdict"].(string); conf, _ := out["confidence"].(float64)
        switch verdict { case "agree": accepts++; peerConf = append(peerConf, conf); case "disagree": rejects++; default: uncertain++ }
    }
    rule := a.Quorum.Quorum["class_b"]; result := swarm.Evaluate(rule, 0.82, accepts, rejects, uncertain, len(a.Peers.Peers), peerConf); result.IntentID = intentID; result.ActionClass = "class_b"; a.State.SetLastQuorum(result); return result, nil
}
func main() {
    var peers swarm.PeersConfig; var node nodeConfig; var quorum swarm.QuorumConfig
    if err := policy.LoadYAML(getenv("PEERS_FILE", "./configs/peers-rt-a.yaml"), &peers); err != nil { log.Fatal(err) }
    if err := policy.LoadYAML(getenv("NODE_FILE", "./configs/node-rt-a.yaml"), &node); err != nil { log.Fatal(err) }
    if err := policy.LoadYAML(getenv("QUORUM_FILE", "./configs/quorum.yaml"), &quorum); err != nil { log.Fatal(err) }
    priv := mustLoadPrivateKey(node.Node.PrivateKeyFile)
    client := lip4d.NewClient(node.Node.ID, node.Node.Role, "ed25519:"+node.Node.ID+":1", priv)
    state := dashboard.NewState(node.Node.ID, node.Node.Role); state.SetPeers(peers.Peers); state.AddEvent("startup", "routefastd started")
    poller := dashboard.NewPeerPoller(peers, state)
    app := &App{Node: node, Peers: peers, Quorum: quorum, Client: client, State: state, Poller: poller}
    ctx, cancel := context.WithCancel(context.Background()); defer cancel(); go poller.Start(ctx, 5*time.Second)
    mux := http.NewServeMux(); dashboard.NewHandlers(state, app, app).Register(mux)
    addr := getenv("ROUTEFASTD_ADDR", ":8080")
    log.Printf("routefastd dashboard listening on %s", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}
