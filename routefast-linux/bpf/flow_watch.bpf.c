// SPDX-License-Identifier: GPL-2.0
// routefast-linux v0.4 starter eBPF skeleton
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

struct event_t {
    __u64 ts_ns;
    __u32 src_ip;
    __u32 dst_ip;
    __u32 packets;
};

struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 20);
} events SEC(".maps");

SEC("tracepoint/net/netif_receive_skb")
int flow_watch(void *ctx) {
    struct event_t *e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) return 0;
    e->ts_ns = bpf_ktime_get_ns();
    e->src_ip = 0;
    e->dst_ip = 0;
    e->packets = 1;
    bpf_ringbuf_submit(e, 0);
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
