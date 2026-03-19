package httpx

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

var Client = &http.Client{Timeout: 5 * time.Second}

func PostJSON(ctx context.Context, url string, reqBody any, respBody any) error {
    body, err := json.Marshal(reqBody)
    if err != nil { return err }
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
    if err != nil { return err }
    req.Header.Set("Content-Type", "application/json")
    resp, err := Client.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    data, _ := io.ReadAll(resp.Body)
    if resp.StatusCode >= 300 {
        return fmt.Errorf("%s: %s", resp.Status, string(data))
    }
    if respBody != nil && len(data) > 0 {
        return json.Unmarshal(data, respBody)
    }
    return nil
}
