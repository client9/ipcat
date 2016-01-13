package ipset

import (
	"testing"
)

func TestSetting(t *testing.T) {
	set := NewIntervalSet(100)
	if set.AddRange("10.0.0.0", "10.0.1.0") == false {
		// Smoke test
		t.Fatalf("Adding failed")
	}
	if set.Len() != 1 {
		// Smoke test
		// major pass-by-val, not pass-by-ref problem
		t.Fatalf("Added entry but size %d != 1", set.Len())
	}

	if set.AddRange("10.0.0.0", "11.0.0.0") == true {
		t.Errorf("Allowed adding something larger than class A network")
	}

	if set.AddRange("1.0.0.0", "1.255.255.255") == false {
		t.Errorf("Unable to add full class A network")
	}

	if set.AddRange("10.1.0.0", "10.0.0.0") == true {
		t.Errorf("Allowed backwards interval")
	}

	if set.AddRange("2.0.0.0", "2.0.0.0") == false {
		t.Errorf("Unable to add single IP")
	}

	if set.AddRange("Busted", "2.0.0.0") == true {
		t.Errorf("Allows busted input1")
	}

	if set.AddRange("2.0.0.0", "Busted") == true {
		t.Errorf("Allows busted input2")
	}
}

/*
func BenchmarkLookup(b *testing.B) {

	for n := 0; n < b.N; n++ {
		Find("0.255.255.255")
	}
}
*/
