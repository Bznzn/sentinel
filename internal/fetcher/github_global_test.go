package fetcher

import (
	"context"
	"testing"

	"SuperBizAgent/internal/model"
)

func TestGitHubFetcher_FetchGlobalAdvisories(t *testing.T) {
	fetcher := NewGitHubFetcher()

	sub := &model.Subscription{
		ID:        1,
		Name:      "GitHub Global Advisories",
		SourceURL: "https://api.github.com/advisories",
	}

	ctx := context.Background()
	events, err := fetcher.Fetch(ctx, sub)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	t.Logf("Fetched %d global advisories", len(events))

	if len(events) == 0 {
		t.Fatal("Expected at least 1 advisory")
	}

	for i, e := range events {
		if i >= 5 {
			break
		}
		t.Logf("[%d] %s | Severity: %s | CVE: %s", i+1, e.Title, e.Severity, e.CVEID)
		t.Logf("    URL: %s", e.SourceURL)
	}
}
