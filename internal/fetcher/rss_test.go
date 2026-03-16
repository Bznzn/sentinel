package fetcher

import (
	"context"
	"testing"

	"SuperBizAgent/internal/model"
)

func TestRSSFetch(t *testing.T) {
	fetcher := NewRSSFetcher()
	ctx := context.Background()

	testURLs := []string{
		//"https://nvd.nist.gov/feeds/xml/cve/misc/nvd-rss.xml",
		"https://www.sentinelone.com/labs/feed/",
		//"https://paper.seebug.org/rss/",
		//"https://www.freebuf.com/feed",
		//"https://www.4hou.com/feed",
		//"https://wechat2rss.xlab.app/feed/edf17c7a01a7152b7a1ca14133258dc4abc1ca61.xml",
	}

	for _, url := range testURLs {
		t.Run(url, func(t *testing.T) {
			sub := &model.Subscription{
				SourceURL: url,
			}

			events, err := fetcher.Fetch(ctx, sub)
			if err != nil {
				t.Fatalf("Failed to fetch RSS: %v", err)
			}

			t.Logf("Successfully fetched %d items from %s", len(events), url)

			if len(events) > 0 {
				t.Logf("First item: %s", events[0].Title)
			}
		})
	}
}
