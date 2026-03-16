package fetcher

import (
	"context"
	"testing"

	"SuperBizAgent/internal/model"
)

func TestGitHubFetcher_Fetch(t *testing.T) {
	fetcher := NewGitHubFetcher()
	ctx := context.Background()

	testCases := []struct {
		name string
		url  string
	}{
		{"electron", "https://github.com/electron/electron"},
		{"rails", "https://github.com/rails/rails"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sub := &model.Subscription{
				ID:        1,
				SourceURL: tc.url,
			}

			events, err := fetcher.Fetch(ctx, sub)
			if err != nil {
				t.Errorf("Fetch failed for %s: %v", tc.name, err)
				return
			}

			t.Logf("[%s] Fetched %d events", tc.name, len(events))
			for i, event := range events {
				if i >= 2 {
					break
				}
				t.Logf("  - %s | %s", event.Title, event.Severity)
			}
		})
	}
}
