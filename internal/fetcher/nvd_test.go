package fetcher

import (
	"context"
	"testing"

	"SuperBizAgent/internal/model"
)

func TestNVDFetcher_Fetch(t *testing.T) {
	fetcher := NewNVDFetcher("")
	ctx := context.Background()

	sub := &model.Subscription{
		ID:        1,
		SourceURL: "https://services.nvd.nist.gov/rest/json/cves/2.0",
	}

	events, err := fetcher.Fetch(ctx, sub)
	if err != nil {
		t.Errorf("Fetch failed: %v", err)
		return
	}

	t.Logf("Fetched %d events", len(events))
	for i, event := range events {
		if i >= 3 {
			break
		}
		t.Logf("  - %s | %s | Vendor: %s | Product: %s | CVSS: %.1f",
			event.CVEID, event.Severity, event.AffectedVendor, event.AffectedProduct, event.CVSSScore)
	}
}

func TestParseCPE(t *testing.T) {
	tests := []struct {
		cpe     string
		vendor  string
		product string
	}{
		{"cpe:2.3:a:apache:struts:2.5.30", "Apache", "Struts"},
		{"cpe:2.3:a:vmware:spring_framework:5.3.17", "Vmware", "Spring Framework"},
		{"cpe:2.3:o:microsoft:windows_10:1809", "Microsoft", "Windows 10"},
		{"cpe:2.3:a:*:product:1.0", "", "Product"},
	}

	for _, tt := range tests {
		t.Run(tt.cpe, func(t *testing.T) {
			vendor, product := parseCPE(tt.cpe)
			if vendor != tt.vendor || product != tt.product {
				t.Errorf("parseCPE(%s) = (%s, %s), want (%s, %s)",
					tt.cpe, vendor, product, tt.vendor, tt.product)
			}
		})
	}
}
