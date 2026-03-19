package main

import (
	"context"
	"log"

	"github.com/aovidi/routefast-linux/internal/actuation"
	"github.com/aovidi/routefast-linux/internal/reasoning"
)

func main() {
	exec := actuation.NewExecutor()
	decision := reasoning.Decision{Action: "SHIFT_SERVICE_CLASS_TO_BEARER", Target: "critical-control", Constraints: map[string]interface{}{"prefer_media": []string{"radio", "lte"}}, Reason: "demo"}
	if err := exec.Apply(context.Background(), decision); err != nil {
		log.Fatal(err)
	}
}
