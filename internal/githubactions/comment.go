package githubactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// PostComment posts body as a new comment on the given pull request using the
// GitHub issues-comments API. A non-2xx response (including the response body,
// truncated) is returned as an error.
func PostComment(client *http.Client, env Environment, prNumber int, body string) error {
	url := fmt.Sprintf("%s/repos/%s/issues/%d/comments",
		strings.TrimRight(env.APIURL, "/"), env.Repository, prNumber)

	payload, err := json.Marshal(map[string]string{"body": body})
	if err != nil {
		return fmt.Errorf("encoding comment: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+env.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("posting PR comment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("posting PR comment: status %d: %s", resp.StatusCode, strings.TrimSpace(string(snippet)))
	}
	return nil
}
