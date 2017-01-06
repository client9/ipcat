package ipcat

import (
	"os/exec"
	"strings"
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

// DownloadAppEngine downloads and returns raw bytes of the Google App Engine ip
// range list
func DownloadAppEngine() ([]string, error) {
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

// UpdateAppEngine takes a raw data, parses it and updates the ipmap
func UpdateAppEngine(ipmap *IntervalSet, ranges []string) error {
	const (
		dcName = "Google App Engine"
		dcURL  = "https://cloud.google.com/appengine"
	)

	for _, ipRange := range ranges {
		err := ipmap.AddCIDR(ipRange, dcName, dcURL)
		if err != nil {
			return err
		}
	}

	return nil
}
