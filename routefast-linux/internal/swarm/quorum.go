package swarm

type Rule struct {
    MinCorroborations int `yaml:"min_corroborations"`
    MinAccepts int `yaml:"min_accepts"`
    MinPeersContacted int `yaml:"min_peers_contacted"`
    TimeoutMS int `yaml:"timeout_ms"`
    DowngradeOnTimeout bool `yaml:"downgrade_on_timeout"`
}
type QuorumConfig struct { Quorum map[string]Rule `yaml:"quorum"` }
type Result struct {
    IntentID string `json:"intent_id"`
    ActionClass string `json:"action_class"`
    LocalConfidence float64 `json:"local_confidence"`
    ResponsesReceived int `json:"responses_received"`
    Accepts int `json:"accepts"`
    Rejects int `json:"rejects"`
    Uncertain int `json:"uncertain"`
    MergedConfidence float64 `json:"merged_confidence"`
    Decision string `json:"decision"`
}
func Evaluate(rule Rule, local float64, accepts, rejects, uncertain, peersContacted int, peerConf []float64) Result {
    merged := local
    if len(peerConf) > 0 {
        maxp, sum := peerConf[0], 0.0
        for _, p := range peerConf { sum += p; if p > maxp { maxp = p } }
        avg := sum / float64(len(peerConf))
        merged = 0.50*local + 0.30*maxp + 0.20*avg
    }
    d := "denied"
    if peersContacted >= rule.MinPeersContacted && accepts >= rule.MinAccepts { d = "approved" }
    return Result{LocalConfidence: local, ResponsesReceived: accepts+rejects+uncertain, Accepts: accepts, Rejects: rejects, Uncertain: uncertain, MergedConfidence: merged, Decision: d}
}
