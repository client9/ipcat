package ipcat

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	cloudflareDownload = "https://www.cloudflare.com/ips-v4"
)

// DownloadCloudflare downloads the latest Cloudflare IP ranges list
func DownloadCloudflare() ([]byte, error) {
	resp, err := http.Get(cloudflareDownload)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to download Cloudflare ranges: status code %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return bytes.TrimSpace(body), nil
}

// UpdateCloudflare parses the Cloudflare IP text file and updates the interval set
func UpdateCloudflare(ipmap *IntervalSet, body []byte) error {
	const (
		cloudflareName = "Cloudflare Inc"
		cloudflareURL  = "https://www.cloudflare.com/"
	)

	// delete all existing records
	ipmap.DeleteByName(cloudflareName)

	// and add back
	for _, cidr := range bytes.Split(body, []byte("\n")) {
		err := ipmap.AddCIDR(string(cidr), cloudflareName, cloudflareURL)
		if err != nil {
			return err
		}
	}

	return nil
}
