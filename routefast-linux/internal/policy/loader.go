package policy

import (
    "os"
    "gopkg.in/yaml.v3"
)

func LoadYAML[T any](path string, out *T) error {
    b, err := os.ReadFile(path)
    if err != nil { return err }
    return yaml.Unmarshal(b, out)
}
