package ipcat

import (
	"testing"
)

func TestSetting(t *testing.T) {
	const (
		dcName = "foobar"
		dcURL  = "http://junk"
	)

	set := NewIntervalSet(100)
	if err := set.AddRange("10.0.0.0", "10.0.1.0", dcName, dcURL); err != nil {
		// Smoke test
		t.Fatalf("Adding failed: %s", err)
	}
	if set.Len() != 1 {
		// Smoke test
		// major pass-by-val, not pass-by-ref problem
		t.Fatalf("Added entry but size %d != 1", set.Len())
	}

	if err := set.AddRange("1.0.0.0", "1.255.255.255", dcName, dcURL); err != nil {
		t.Errorf("Unable to add full class A network: %s", err)
	}

	if err := set.AddRange("10.1.0.0", "10.0.0.0", dcName, dcURL); err == nil {
		t.Errorf("Allowed backwards interval")
	}

	if err := set.AddRange("2.0.0.0", "2.0.0.0", dcName, dcURL); err != nil {
		t.Errorf("Unable to add single IP: %s", err)
	}

	if err := set.AddRange("Busted", "2.0.0.0", dcName, dcURL); err == nil {
		t.Errorf("Allows busted input1")
	}

	if err := set.AddRange("2.0.0.0", "Busted", dcName, dcURL); err == nil {
		t.Errorf("Allows busted input2")
	}
}

var ipStringTests = []string{
	"4.4.4.4",
	"192.168.1.1",
	"fe80:cd00::211e:729c",
}

func TestIPString(t *testing.T) {
	for _, tt := range ipStringTests {
		ip := IPParse(tt)
		if got := IPString(ip); got != tt {
			t.Errorf("IPString(%d) = %q, want %q", ip, got, tt)
		}
	}
}

var cidr2RangeTests = []struct {
	cidr  string
	start string
	end   string
}{
	{"10.0.0.0/8", "10.0.0.0", "10.255.255.255"},
	{"192.168.0.0/24", "192.168.0.0", "192.168.0.255"},
}

func TestCIDR2Range(t *testing.T) {
	for _, tt := range cidr2RangeTests {
		start, end, err := CIDR2Range(tt.cidr)
		if err != nil {
			t.Fatal(err)
		}
		if start != tt.start {
			t.Errorf("CIDR2Range(%q) start = %q, want %q", tt.cidr, start, tt.start)
		}
		if end != tt.end {
			t.Errorf("CIDR2Range(%q) end = %q, want %q", tt.cidr, end, tt.start)
		}
	}
}

func TestIntervalSet(t *testing.T) {
	ipset := NewIntervalSet(256)
	ipset.AddCIDR("1.1.1.1/24", "Test Range", "Test URL")
	rec, err := ipset.Contains("1.1.1.2")
	if err != nil {
		t.Fatalf("ipset.Contains(%q) error: %s", "1.1.1.2", err)
	}
	if rec == nil {
		t.Errorf("ipset.Contains(%q) record is nil", "1.1.1.2")
	}
	if _, err := ipset.Contains("0"); err == nil {
		t.Errorf("ipset.Contains(%q) = %v, want error", "1.1.1.2", err)
	}

	ipset.AddCIDR("2.2.2.2/24", "Another Test Range", "Another Test URL")
	ipset.DeleteByName("Test Range")
	rec, err = ipset.Contains("1.1.1.2")
	if err != nil {
		t.Fatalf("ipset.Contains(%q) error: %s", "1.1.1.2", err)
	}
	if rec != nil {
		t.Errorf("ipset.Contains(%q) record is not nil after DeleteByName: %v", "1.1.1.2", rec)
	}
}
