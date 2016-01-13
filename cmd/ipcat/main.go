package main

import (
	"github.com/client9/ipcat"
)

func main() {
	set := IntervalSet{}
	err := set.ImportCSV(os.Stdin)
	if err != nil {
		log.Fatalf("Unable to import: %s", err)
	}
	err = set.ExportCSV(os.Stdout)
	if err != nil {
		log.Fatalf("Unable to export: %s", err)
	}
}
