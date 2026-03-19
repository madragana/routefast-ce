package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "time"
    rfcrypto "github.com/aovidi/routefast-linux/internal/crypto"
    "github.com/aovidi/routefast-linux/internal/lip4d"
    "github.com/aovidi/routefast-linux/internal/policy"
    "github.com/aovidi/routefast-linux/internal/swarm"
)
func getenv(k, def string) string { v := os.Getenv(k); if v == "" { return def }; return v }
func main() {
    port := getenv("PORT", "9443")
    peersFile := getenv("PEERS_FILE", "./configs/peers-rt-a.yaml")
    nodeFile := getenv("NODE_FILE", "./configs/node-rt-a.yaml")
    var peers swarm.PeersConfig; var node swarm.NodeConfig
    if err := policy.LoadYAML(peersFile, &peers); err != nil { log.Fatal(err) }
    if err := policy.LoadYAML(nodeFile, &node); err != nil { log.Fatal(err) }
    priv, err := rfcrypto.LoadPrivateKey(node.Node.PrivateKeyFile); if err != nil { log.Fatal(err) }
    mux := http.NewServeMux(); srv := lip4d.NewServer(peers); srv.Register(mux)
    ctx, cancel := context.WithCancel(context.Background()); defer cancel()
    keyID := "ed25519:" + node.Node.ID + ":1"
    go swarm.NewHeartbeatSender(node.Node.ID, node.Node.Role, keyID, priv, peers).Start(ctx, 5*time.Second)
    go swarm.NewPeerStateSender(node.Node.ID, node.Node.Role, keyID, priv, peers).Start(ctx, 8*time.Second)
    log.Printf("lip4dd node=%s role=%s listening on :%s", node.Node.ID, node.Node.Role, port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}
