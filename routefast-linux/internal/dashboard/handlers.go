package dashboard

import (
    "encoding/json"
    "net/http"
)
type DemoRunner interface { RunDemoCorroboration() (any, error) }
type PollRunner interface { PollPeersNow() error }
type Handlers struct { State *State; Demo DemoRunner; Poll PollRunner }
func NewHandlers(state *State, demo DemoRunner, poll PollRunner) *Handlers { return &Handlers{State: state, Demo: demo, Poll: poll} }
func (h *Handlers) Register(mux *http.ServeMux) {
    mux.HandleFunc("/dashboard", h.handleDashboard)
    mux.HandleFunc("/api/status", h.handleStatus)
    mux.HandleFunc("/api/demo/corroborate", h.handleDemoCorroborate)
    mux.HandleFunc("/api/poll", h.handlePoll)
}
func (h *Handlers) handleDashboard(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "text/html; charset=utf-8"); _, _ = w.Write([]byte(PageHTML)) }
func (h *Handlers) handleStatus(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "application/json"); _ = json.NewEncoder(w).Encode(h.State.Snapshot()) }
func (h *Handlers) handleDemoCorroborate(w http.ResponseWriter, r *http.Request) { out, err := h.Demo.RunDemoCorroboration(); if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }; _ = json.NewEncoder(w).Encode(out) }
func (h *Handlers) handlePoll(w http.ResponseWriter, r *http.Request) { if err := h.Poll.PollPeersNow(); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }; _ = json.NewEncoder(w).Encode(map[string]any{"ok": true}) }
