package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/client9/ipcat"
)

func main() {
	lookup := flag.String("l", "", "lookup an IP address")
	updateAWS := flag.Bool("aws", false, "update AWS records")
	updateAzure := flag.Bool("azure", false, "update Azure records")
	datafile := flag.String("csvfile", "datacenters.csv", "read/write from this file")
	statsfile := flag.String("statsfile", "datacenters-stats.csv", "write statistics to this file")
	addCIDR := flag.String("addcidr", "", "add this CIDR range to the data file [CIDR,name,url]")
	flag.Parse()

	filein, err := os.Open(*datafile)
	if err != nil {
		log.Fatalf("Unable to read %s: %s", *datafile, err)
	}
	set := ipcat.IntervalSet{}
	err = set.ImportCSV(filein)
	if err != nil {
		log.Fatalf("Unable to import: %s", err)
	}
	filein.Close()
	log.Printf("Loaded %d entries", set.Len())

	if *lookup != "" {
		rec, err := set.Contains(*lookup)
		if err != nil {
			log.Fatalf("Unable to find %s: %s", *lookup, err)
		}
		if rec == nil {
			log.Fatalf("Not found: %s", *lookup)
		}
		fmt.Printf("[%s:%s] %s %s\n", rec.LeftDots, rec.RightDots, rec.Name, rec.URL)
		return
	}

	if *updateAWS {
		body, err := ipcat.DownloadAWS()
		if err != nil {
			log.Fatalf("Unable to download AWS rules: %s", err)
		}
		err = ipcat.UpdateAWS(&set, body)
		if err != nil {
			log.Fatalf("Unable to parse AWS rules: %s", err)
		}
	}

	if *updateAzure {
		body, err := ipcat.DownloadAzure()
		if err != nil {
			log.Fatalf("Unable to download Azure rules: %s", err)
		}
		err = ipcat.UpdateAzure(&set, body)
		if err != nil {
			log.Fatalf("Unable to parse Azure rules: %s", err)
		}
	}

	if *statsfile != "" {
		fileout, err := os.OpenFile(*statsfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalf("Unable to open file to write: %s", err)
		}
		list := set.RankBySize()
		fileout.WriteString("Datacenter Name, Total IPs\n")
		for _, val := range list {
			fileout.WriteString(fmt.Sprintf("%s,%d\n", val.Name, val.Size))
		}
		fileout.Close()
	}

	if *addCIDR != "" {
		t := strings.Split(*addCIDR, ",")
		if len(t) != 3 {
			log.Fatal("range must be in format: CIDR,name,url")
		}
		err := set.AddCIDR(t[0], t[1], t[2])
		if err != nil {
			log.Fatalf("Could not add range: %v", err)
		}
		log.Println("Range added successfully")
	}

	fileout, err := os.OpenFile(*datafile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Unable to open file to write: %s", err)
	}
	err = set.ExportCSV(fileout)
	if err != nil {
		log.Fatalf("Unable to export: %s", err)
	}
	fileout.Close()
}
