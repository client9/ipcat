package ipcat

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	cloudflareDownload = []string{
		"https://www.cloudflare.com/ips-v4",
		"https://www.cloudflare.com/ips-v6",
	}
)

// DownloadCloudflare downloads the latest Cloudflare IP ranges list
func DownloadCloudflare() ([]byte, error) {
	readers := make([]io.Reader, 0, len(cloudflareDownload))

	defer func() {
		for _, reader := range readers {
			reader.(io.ReadCloser).Close()
		}
	}()

	for _, uri := range cloudflareDownload {
		resp, err := http.Get(uri)
		if err != nil {
			return nil, err
		}
		readers = append(readers, resp.Body)

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("Failed to download Cloudflare ranges: status code %s", resp.Status)
		}
	}

	body, err := ioutil.ReadAll(io.MultiReader(readers...))
	if err != nil {
		return nil, err
	}

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
