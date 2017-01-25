package ipcat

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// http://www.microsoft.com/en-us/download/confirmation.aspx?id=41653
var (
	msAzure = "https://download.microsoft.com/download/0/1/8/018E208D-54F8-44CD-AA26-CD7BC9524A8C/PublicIPs_20170123.xml"
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

// DownloadAzure downloads and returns raw bytes of the MS Azure ip
// range list
func DownloadAzure() ([]byte, error) {
	resp, err := http.Get(msAzure)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to download Azure ranges: status code %s", resp.Status)
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
