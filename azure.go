package ipcat

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// AzureIPRange is a MS Azure record
type AzureIPRange struct {
	Subnet string `xml:"Subnet,attr"`
}

// AzureRegion is a MS Region
type AzureRegion struct {
	Name    string         `xml:"Name,attr"`
	IPRange []AzureIPRange `xml:"IpRange"`
}

// AzurePublicIPAddresses is a listing of regions
type AzurePublicIPAddresses struct {
	AzureRegion []AzureRegion `xml:"Region"`
}

var retried bool

var findPublicIPsURL = func() (string, error) {
	downloadPage := "http://www.microsoft.com/en-us/download/confirmation.aspx?id=41653"

	resp, err := http.Get(downloadPage)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("url=https://download.microsoft.com/download/.*/PublicIPs_.*.xml")
	addr := re.Find(b)

	if string(addr) == "" {
		return "", errors.New("could not find PublicIPs address on download page")
	}

	return string(addr)[4:], nil
}

// DownloadAzure downloads and returns raw bytes of the MS Azure ip
// range list
func DownloadAzure() ([]byte, error) {
	url, err := findPublicIPsURL()
	if err != nil {
		return nil, fmt.Errorf("failed to find public IPs url during retry: %s", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable read body: %s", err)
	}
	resp.Body.Close()
	return body, nil
}

// UpdateAzure takes a raw data, parses it and updates the ipmap
func UpdateAzure(ipmap *IntervalSet, body []byte) error {
	const (
		dcName = "Microsoft Azure"
		dcURL  = "http://www.windowsazure.com/en-us/"
	)

	azure := AzurePublicIPAddresses{}
	err := xml.Unmarshal(body, &azure)
	if err != nil {
		return err
	}

	// delete all existing records
	ipmap.DeleteByName(dcName)

	for _, region := range azure.AzureRegion {
		for _, rng := range region.IPRange {
			err = ipmap.AddCIDR(rng.Subnet, dcName, dcURL)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
