package main

import (
	"flag"
	"log"
	"os"

	"github.com/client9/ipcat"
)

func main() {
	updateAWS := flag.Bool("aws", false, "update AWS records")
	flag.Parse()

	set := ipcat.IntervalSet{}
	err := set.ImportCSV(os.Stdin)
	if err != nil {
		log.Fatalf("Unable to import: %s", err)
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

	err = set.ExportCSV(os.Stdout)
	if err != nil {
		log.Fatalf("Unable to export: %s", err)
	}
}
