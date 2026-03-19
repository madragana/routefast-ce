package lip4d

import "sync"

type ReplayCache struct { mu sync.Mutex; seen map[string]int64 }
func NewReplayCache() *ReplayCache { return &ReplayCache{seen: map[string]int64{}} }
func (r *ReplayCache) Seen(sender string, seq int64) bool {
    r.mu.Lock(); defer r.mu.Unlock()
    last, ok := r.seen[sender]
    if ok && seq <= last { return true }
    r.seen[sender] = seq
    return false
}
