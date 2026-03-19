# RouteFast v0.5 Demo

## Nodes
- rt-a
- rt-b
- rt-c

## Run

```bash
PORT=9443 PEERS_FILE=./configs/peers-rt-a.yaml NODE_FILE=./configs/node-rt-a.yaml go run ./cmd/lip4dd
PORT=9444 PEERS_FILE=./configs/peers-rt-b.yaml NODE_FILE=./configs/node-rt-b.yaml go run ./cmd/lip4dd
PORT=9445 PEERS_FILE=./configs/peers-rt-c.yaml NODE_FILE=./configs/node-rt-c.yaml go run ./cmd/lip4dd
PEERS_FILE=./configs/peers-rt-a.yaml NODE_FILE=./configs/node-rt-a.yaml QUORUM_FILE=./configs/quorum.yaml ROUTEFASTD_ADDR=:8080 go run ./cmd/routefastd
```
