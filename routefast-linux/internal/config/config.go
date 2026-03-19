package config

import (
    "bufio"
    "os"
    "strconv"
    "strings"
)

type Policy struct {
    NodeID              string
    Domain              string
    AllowedActions      []string
    HighImpactActions   []string
    MaxImpactedPrefixes int
    MaxTTLSeconds       int
    RequireQuorum       bool
    ActuatorURL         string
    LIP4DURL            string
    ListenAddr          string
    AuditPath           string
    DryRun              bool
    KeySeed             string
}

type Peer struct {
    NodeID    string
    URL       string
    PublicKey string
    Role      string
}

type Bearer struct {
    Name             string
    MediaType        string
    AdminState       string
    OperState        string
    Priority         int
    JamResistance    float64
    TrustScore       float64
    ThroughputMbps   int
    PreferredFor     []string
}

func LoadPolicy(path string) (Policy, error) {
    lines, err := readLines(path)
    if err != nil {
        return Policy{}, err
    }
    var cfg Policy
    var list string
    for _, line := range lines {
        if strings.HasPrefix(line, "- ") {
            item := strings.TrimSpace(strings.TrimPrefix(line, "- "))
            switch list {
            case "allowed_actions":
                cfg.AllowedActions = append(cfg.AllowedActions, item)
            case "high_impact_actions":
                cfg.HighImpactActions = append(cfg.HighImpactActions, item)
            }
            continue
        }
        key, val, ok := splitKV(line)
        if !ok {
            continue
        }
        list = ""
        switch key {
        case "node_id":
            cfg.NodeID = val
        case "domain":
            cfg.Domain = val
        case "allowed_actions", "high_impact_actions":
            list = key
        case "max_impacted_prefixes":
            cfg.MaxImpactedPrefixes, _ = strconv.Atoi(val)
        case "max_ttl_seconds":
            cfg.MaxTTLSeconds, _ = strconv.Atoi(val)
        case "require_quorum":
            cfg.RequireQuorum = strings.EqualFold(val, "true")
        case "actuator_url":
            cfg.ActuatorURL = val
        case "lip4d_url":
            cfg.LIP4DURL = val
        case "listen_addr":
            cfg.ListenAddr = val
        case "audit_path":
            cfg.AuditPath = val
        case "dry_run":
            cfg.DryRun = strings.EqualFold(val, "true")
        case "key_seed":
            cfg.KeySeed = val
        }
    }
    if cfg.ListenAddr == "" { cfg.ListenAddr = ":8090" }
    if cfg.ActuatorURL == "" { cfg.ActuatorURL = "http://127.0.0.1:8092" }
    if cfg.LIP4DURL == "" { cfg.LIP4DURL = "http://127.0.0.1:8091" }
    if cfg.AuditPath == "" { cfg.AuditPath = "routefast-audit.jsonl" }
    if cfg.Domain == "" { cfg.Domain = "routefast.local" }
    return cfg, nil
}

func LoadPeers(path string) ([]Peer, error) {
    lines, err := readLines(path)
    if err != nil {
        return nil, err
    }
    var out []Peer
    var cur *Peer
    for _, line := range lines {
        if strings.HasPrefix(line, "- ") {
            if cur != nil { out = append(out, *cur) }
            cur = &Peer{}
            line = strings.TrimSpace(strings.TrimPrefix(line, "- "))
            if line != "" {
                key, val, ok := splitKV(line)
                if ok { setPeer(cur, key, val) }
            }
            continue
        }
        if cur == nil { continue }
        key, val, ok := splitKV(line)
        if ok { setPeer(cur, key, val) }
    }
    if cur != nil { out = append(out, *cur) }
    return out, nil
}

func LoadBearers(path string) ([]Bearer, error) {
    lines, err := readLines(path)
    if err != nil {
        return nil, err
    }
    var out []Bearer
    var cur *Bearer
    var list string
    for _, line := range lines {
        if strings.HasPrefix(line, "- ") {
            if cur != nil { out = append(out, *cur) }
            cur = &Bearer{}
            list = ""
            line = strings.TrimSpace(strings.TrimPrefix(line, "- "))
            if line != "" {
                key, val, ok := splitKV(line)
                if ok { list = setBearer(cur, key, val) }
            }
            continue
        }
        if cur == nil { continue }
        if strings.HasPrefix(line, "  - ") || strings.HasPrefix(line, "- ") {
            item := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "- "))
            if list == "preferred_for" {
                cur.PreferredFor = append(cur.PreferredFor, item)
            }
            continue
        }
        key, val, ok := splitKV(line)
        if ok { list = setBearer(cur, key, val) }
    }
    if cur != nil { out = append(out, *cur) }
    return out, nil
}

func readLines(path string) ([]string, error) {
    f, err := os.Open(path)
    if err != nil { return nil, err }
    defer f.Close()
    var out []string
    s := bufio.NewScanner(f)
    for s.Scan() {
        line := strings.TrimRight(s.Text(), " \t\r\n")
        trimmed := strings.TrimSpace(line)
        if trimmed == "" || strings.HasPrefix(trimmed, "#") { continue }
        out = append(out, line)
    }
    return out, s.Err()
}

func splitKV(line string) (string,string,bool) {
    trimmed := strings.TrimSpace(line)
    parts := strings.SplitN(trimmed, ":", 2)
    if len(parts) != 2 { return "","",false }
    return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
}

func setPeer(p *Peer, key, val string) {
    switch key {
    case "node_id": p.NodeID = val
    case "url": p.URL = val
    case "public_key": p.PublicKey = val
    case "role": p.Role = val
    }
}

func setBearer(b *Bearer, key, val string) string {
    switch key {
    case "name": b.Name = val
    case "media_type": b.MediaType = val
    case "admin_state": b.AdminState = val
    case "oper_state": b.OperState = val
    case "priority": b.Priority, _ = strconv.Atoi(val)
    case "jam_resistance": b.JamResistance, _ = strconv.ParseFloat(val, 64)
    case "trust_score": b.TrustScore, _ = strconv.ParseFloat(val, 64)
    case "throughput_mbps": b.ThroughputMbps, _ = strconv.Atoi(val)
    case "preferred_for": return "preferred_for"
    }
    return ""
}
