package actuation

import (
	"context"
	"fmt"
	"log"

	"github.com/aovidi/routefast-linux/internal/reasoning"
)

type Executor struct{}

func NewExecutor() *Executor { return &Executor{} }

func (e *Executor) Apply(_ context.Context, d reasoning.Decision) error {
	switch d.Action {
	case "EXPLAIN_ONLY":
		log.Printf("no-op: %s", d.Reason)
		return nil
	case "BLOCK_SOURCE":
		log.Printf("would apply nftables block for %s", d.Target)
		return nil
	case "SHIFT_SERVICE_CLASS_TO_BEARER":
		log.Printf("would adjust policy routing for target %s with constraints=%v", d.Target, d.Constraints)
		return nil
	default:
		return fmt.Errorf("unsupported action %s", d.Action)
	}
}
