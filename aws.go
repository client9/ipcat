package ipcat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// AWSPrefix is AWS prefix in their IP ranges file
type AWSPrefix struct {
	IPPrefix string `json:"ip_prefix"`
	Region   string `json:"region"`
	Service  string `json:"service"`
}

// AWS is main record for AWS IP info
type AWS struct {
	SyncToken  string      `json:"syncToken"`
	CreateDate string      `json:"createDate"`
	Prefixes   []AWSPrefix `json:"prefixes"`
}

func DownloadAWS() ([]byte, error) {
	const (
		awsDownload = "https://ip-ranges.amazonaws.com/ip-ranges.json"
	)
	resp, err := http.Get(awsDownload)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to download AWS ranges: status code %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return body, nil

}

// UpdateAWS
func UpdateAWS(ipmap *IntervalSet, body []byte) error {
	const (
		awsName = "Amazon AWS"
		awsURL  = "http://www.amazon.com/aws/"
	)

	aws := AWS{}
	err := json.Unmarshal(body, &aws)
	if err != nil {
		return err
	}

	// delete all existing records
	ipmap.DeleteByName(awsName)

	// and add back
	for _, rec := range aws.Prefixes {
		if rec.Service == "EC2" {
			err := ipmap.AddCIDR(rec.IPPrefix, awsName, awsURL)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
