package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/client9/ipcat"
)

func main() {
	updateAWS := flag.Bool("aws", false, "update AWS records")
	datafile := flag.String("csvfile", "datacenters.csv", "read/write from this file")
	statsfile := flag.String("statsfile", "datacenters-stats.csv", "write statistics to this file")
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

	fileout, err := os.OpenFile(*datafile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Unable to open file to write: %s", err)
	}
	err = set.ExportCSV(fileout)
	if err != nil {
		log.Fatalf("Unable to export: %s", err)
	}
	fileout.Close()
}
