package ipcat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAWS(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `
		{
			"syncToken": "0123456789",
			"createDate": "2016-11-30-23-19-08",
			"prefixes": [
				{
				  "ip_prefix": "216.182.224.0/20",
				  "region": "us-east-1",
				  "service": "AMAZON"
				},
				{
				  "ip_prefix": "13.54.0.0/15",
				  "region": "ap-southeast-2",
				  "service": "EC2"
				}
			]
		}`)

	}))
	awsDownload = ts.URL
	defer ts.Close()
	b, err := DownloadAWS()
	if err != nil {
		t.Fatalf("DownloadAWS() error: %v", err)
	}
	ipset := NewIntervalSet(100)
	err = UpdateAWS(ipset, b)
	if err != nil {
		t.Fatalf("UpdateAWS error: %v", err)
	}
	// UpdateAWS should only add EC2 IP ranges
	rec, err := ipset.Contains("216.182.224.0")
	if err != nil {
		t.Fatalf("ipset.Contains(%q) error: %v", "216.182.224.0", err)
	}
	if rec != nil {
		t.Errorf("ipset.Contains(%q) rec = %v, want nil", "216.182.224.0", rec)
	}
	rec, err = ipset.Contains("13.54.0.1")
	if err != nil {
		t.Fatalf("ipset.Contains(%q) error: %v", "13.54.0.0", err)
	}
	if rec == nil {
		t.Errorf("ipset.Contains(%q) rec = nil, want exists", "13.54.0.0")
	}
}
