package actuation

import (
    "fmt"

    "github.com/aovidi/routefast-linux/internal/reasoning"
)

func FRRCommands(d reasoning.Decision) ([]string, error) {
    switch d.Action {
    case "LOWER_LOCAL_PREF":
        return []string{
            "vtysh -c 'configure terminal' -c 'route-map ROUTEFAST permit 10' -c 'set local-preference 50'",
        }, nil
    case "REROUTE_PREFIX", "SHIFT_SERVICE_CLASS_TO_BEARER":
        target := d.Target
        if target == "" { target = "203.0.113.0/24" }
        return []string{
            fmt.Sprintf("vtysh -c 'configure terminal' -c 'ip prefix-list ROUTEFAST seq 10 permit %s'", target),
            "vtysh -c 'configure terminal' -c 'route-map ROUTEFAST permit 10' -c 'match ip address prefix-list ROUTEFAST' -c 'set local-preference 200'",
        }, nil
    default:
        return nil, nil
    }
}
