package ipcat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAzure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `
<?xml version="1.0" encoding="utf-8"?>
<AzurePublicIpAddresses xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <Region Name="europewest">
    <IpRange Subnet="40.112.124.0/24" />
    <IpRange Subnet="65.52.128.0/19" />
  </Region>
  <Region Name="useast">
    <IpRange Subnet="23.96.0.0/18" />
    <IpRange Subnet="23.96.64.0/28" />
  </Region>
</AzurePublicIpAddresses>
`)

	}))

	defer ts.Close()
	msazure = ts.URL
	b, err := DownloadAzure()
	if err != nil {
		t.Fatalf("DownloadAzure() error: %v", err)
	}

	ipset := NewIntervalSet(100)
	err = UpdateAzure(ipset, b)
	if err != nil {
		t.Fatalf("UpdateAzure error: %v", err)
	}
	rec, err := ipset.Contains("23.96.0.0")
	if err != nil {
		t.Fatalf("ipset.Contains(%q) error: %v", "23.96.0.0", err)
	}
	if rec == nil {
		t.Errorf("ipset.Contains(%q) rec = nil, want exists", "23.96.0.0")
	}
}
