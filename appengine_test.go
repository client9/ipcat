package ipcat

import (
	"os/exec"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func answers(prefix string, lines []string) []string {
	var answers []string
	for i, line := range lines {
		if strings.HasPrefix(line, ";; ANSWER SECTION:") {
			for _, v := range strings.Fields(lines[i+1]) {
				if strings.HasPrefix(v, prefix) {
					nb := strings.TrimPrefix(v, prefix)
					answers = append(answers, nb)
				}
			}
		}
	}

	return answers
}

func testDownloadAppEngine() ([]string, error) {
	out, err := exec.Command("dig", "-t", "TXT", "_cloud-netblocks.googleusercontent.com", "@ns1.google.com").Output()
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(string(out), "\n")
	netblocks := answers("include:", lines)

	var ranges []string
	for _, nb := range netblocks {
		out, err := exec.Command("dig", "-t", "TXT", nb, "@ns1.google.com").Output()
		if err != nil {
			return []string{}, err
		}
		lines := strings.Split(string(out), "\n")
		rs := answers("ip4:", lines)
		ranges = append(ranges, rs...)
	}

	return ranges, nil
}

func TestDownloadAppEngine(t *testing.T) {
	r1, err := DownloadAppEngine()
	if err != nil {
		t.Fatal(err)
	}

	r2, err := testDownloadAppEngine()
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(r1)
	sort.Strings(r2)

	if !reflect.DeepEqual(r1, r2) {
		t.Error("result differs")
		t.Logf("expected: %q", r2)
		t.Logf("got: %q", r1)
	}
}

func BenchmarkDownloadAppEngine(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := DownloadAppEngine()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTestDownloadAppEngine(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := testDownloadAppEngine()
		if err != nil {
			b.Fatal(err)
		}
	}
}
