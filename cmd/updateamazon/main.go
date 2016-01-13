package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

)

// Prefix is an AWS Prefix
type Prefix struct {
	IPPrefix string `json:"ip_prefix"`
	Region   string `json:"region"`
	Service  string `json:"service"`
}

// AWS is main record for AWS IP info
type AWS struct {
	SyncToken  string   `json:"syncToken"`
	CreateDate string   `json:"createDate"`
	Prefixes   []Prefix `json:"prefixes"`
}

// IPRange is a MS Azure record
type IPRange struct {
	Subnet string `xml:"Subnet,attr"`
}

// DCInfo contains the name and url of a datacenter.
type DCInfo struct {
	Name string
	URL  string
}

func main() {
	flag.Parse()
	filename := flag.Args()[0]
	ipmap := NewIntervalSet(5000)

	log.Printf("Opening original file: %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to open file: %s", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")

		// don't add existing AWS records
		if parts[2] == "Amazon AWS" {
			continue
		}
		err := ipset.AddRange(parts[0], parts[1], parts[2], parts[3])
		if err != nil {
			log.Fatalf("Could not add: %s", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Getting AWS ranges...")
	resp, err := http.Get("https://ip-ranges.amazonaws.com/ip-ranges.json")
	if err != nil {
		log.Fatalf("oops: %s", err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Failed to download AWS ranges: status code %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable read body: %s", err)
	}
	resp.Body.Close()

	aws := AWS{}
	json.Unmarshal(body, &aws)
	if err != nil {
		log.Fatalf("unable to deserialize: %s", aws)
	}
	for _, rec := range aws.Prefixes {
		if rec.Service == "EC2" {
			err := ipmap.AddCIDR(rec.IPPrefix, "Amazon AWS", "http://www.amazon.com/aws/")
			if err != nil {
				log.Fatalf("unable to parse cidr: %s", err)
			}
		}
	}

	err = set.ExportCSV(os.Stdout)
	if err != nil {
		log.Fatalf("Unable to export: %s", err)
	}
}
