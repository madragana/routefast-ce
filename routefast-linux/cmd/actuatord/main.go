package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/aovidi/routefast-linux/internal/actuation"
    "github.com/aovidi/routefast-linux/internal/reasoning"
)

func main() {
    dryRun := os.Getenv("ROUTEFAST_DRY_RUN") != "false"
    exec := actuation.NewExecutor(dryRun)
    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
    mux.HandleFunc("/apply", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
        var d reasoning.Decision
        if err := json.NewDecoder(r.Body).Decode(&d); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        res, err := exec.Apply(r.Context(), d)
        if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(res)
    })
    log.Printf("actuatord listening on :8092 dry_run=%v", dryRun)
    log.Fatal(http.ListenAndServe(":8092", mux))
}
