package ipcat

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	torDownload = "https://panwdbl.appspot.com/lists/ettor.txt"
	torRegexp   = regexp.MustCompile(`^([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})(?:-([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))?$`)
)

// UpdateTor downloads and parses the Tor exit node list and updates the interval set
func UpdateTor(ipmap *IntervalSet) error {
	const (
		torName = "Tor Project"
		torURL  = "https://www.torproject.org/"
	)

	// Download list
	resp, err := http.Get(torDownload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to download Tor ranges: status code %s", resp.Status)
	}

	// Delete all existing records
	ipmap.DeleteByName(torName)

	// Add them back
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// A line with # serves as a comment.
		if commentStart := strings.IndexByte(line, '#'); commentStart >= 0 {
			line = line[:commentStart]
		}

		// A blank line matches no files,
		// so it can serve as a separator for readability.
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		ips := torRegexp.FindStringSubmatch(line)
		switch len(ips) {
		case 2:
			// Found singular IP
			ipmap.AddRange(ips[1], ips[1], torName, torURL)
		case 3:
			// Found range
			ipmap.AddRange(ips[1], ips[2], torName, torURL)
		}
	}

	return nil
}
