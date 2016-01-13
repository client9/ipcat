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

	if err := set.AddRange("10.0.0.0", "11.0.0.0", dcName, dcURL); err == nil {
		t.Errorf("Allowed adding something larger than class A network")
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
