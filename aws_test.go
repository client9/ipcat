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
				},
				{
				  "ipv6_prefix": "2a05:d016::/36",
				  "region": "eu-north-1",
				  "service": "AMAZON"
				},
				{
				  "ipv6_prefix": "2a05:d018::/36",
				  "region": "eu-west-1",
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

	tests := []struct {
		IP       string
		Contains bool
	}{
		{IP: "216.182.224.0", Contains: true},
		{IP: "2a05:d016::20", Contains: true},

		// Only AMAZON services will be contained (superset of all)
		{IP: "13.54.0.1", Contains: false},
		{IP: "2a05:d018::20", Contains: false},
	}

	for _, test := range tests {
		record, err := ipset.Contains(test.IP)
		if err != nil {
			t.Fatalf("ipset.Contains(%q) error: %v", test.IP, err)
		}
		if (record != nil) != test.Contains {
			t.Errorf("ipset.Contains(%q) rec = %v, want rec != nil to be %v",
				test.IP, record, test.Contains)
		}
	}
}
