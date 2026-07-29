package githubactions

import (
	"encoding/json"
	"fmt"
)

// ParsePRNumber extracts the pull request number from a GitHub Actions event
// payload (the JSON at $GITHUB_EVENT_PATH). It reads .pull_request.number and
// falls back to a top-level .number.
func ParsePRNumber(payload []byte) (int, error) {
	var ev struct {
		PullRequest struct {
			Number int `json:"number"`
		} `json:"pull_request"`
		Number int `json:"number"`
	}
	if err := json.Unmarshal(payload, &ev); err != nil {
		return 0, fmt.Errorf("parsing event payload: %w", err)
	}
	if ev.PullRequest.Number != 0 {
		return ev.PullRequest.Number, nil
	}
	if ev.Number != 0 {
		return ev.Number, nil
	}
	return 0, fmt.Errorf("no pull request number in event payload")
}
