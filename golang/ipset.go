package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
)

// generic utility function
//    returns 0 if not valid
func dots2uint32(dots string) uint32 {
	ip := net.ParseIP(dots)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

// ToDots converts a uint32 to a IPv4 Dotted notation
func ToDots(val uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		val>>24,
		(val>>16)&0xFF,
		(val>>8)&0xFF,
		val&0xFF)
}

// IpRange is a closed interval
type interval struct {
	Left      uint32
	Right     uint32
	LeftDots  string
	RightDots string
	Name      string
	URL       string
}

type intervallist []interval

func (ipset intervallist) Len() int {
	return len(ipset)
}
func (ipset intervallist) Less(i, j int) bool {
	return ipset[i].Left < ipset[j].Left
}
func (ipset intervallist) Swap(i, j int) {
	ipset[i], ipset[j] = ipset[j], ipset[i]
}

type IntervalSet struct {
	btree  intervallist
	sorted bool
}

func NewIntervalSet(capacity int) *IntervalSet {
	return &IntervalSet{
		btree: make([]interval, 0, capacity),
	}
}

func (ipset *IntervalSet) ImportCSV(in io.Reader) error {
	ipset.btree = nil
	ipset.sorted = false
	r := csv.NewReader(in)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) != 4 {
			return fmt.Errorf("Expected 4 records but got %d", len(record))
		}
		if !ipset.AddRange(record[0], record[1], record[2], record[3]) {
			return fmt.Errorf("Unable to add record")
		}
	}
	return ipset.sort()
}

// ExportCSV export data to a CSV file
func (ipset *IntervalSet) ExportCSV(in io.Writer) error {
	if !ipset.sorted {
		err := ipset.sort()
		if err != nil {
			return err
		}
	}
	w := csv.NewWriter(in)
	for _, val := range ipset.btree {
		rec := []string{ToDots(val.Left), ToDots(val.Right), val.Name, val.URL}
		if err := w.Write(rec); err != nil {
			return err
		}
	}
	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func (ipset IntervalSet) sort() error {
	if ipset.sorted {
		return nil
	}
	sort.Sort(ipset.btree)
	last := interval{}
	for pos, val := range ipset.btree {
		if val.Left > val.Right {
			return fmt.Errorf("left %d > right %d at pos %d",
				val.Left, val.Right, pos)
		}
		if val.Right-val.Left > (uint32(255) << 24) {
			return fmt.Errorf("Interval too large: [%d,%d]",
				val.Left, val.Right)
		}
		if pos > 0 {
			if val.Left <= last.Right || val.Right <= last.Right {
				return fmt.Errorf("Overlaping regions %v vs. %v", last, val)
			}
		}
		last = val
	}
	ipset.sorted = true
	return nil
}

func (ipset *IntervalSet) AddRange(dotsleft, dotsright, name, url string) bool {
	left := dots2uint32(dotsleft)
	if left == 0 && dotsleft != "0.0.0.0" {
		return false
	}
	right := dots2uint32(dotsright)
	if right == 0 && dotsright != "0.0.0.0" {
		return false
	}
	if left > right {
		return false
	}
	if right-left >= uint32(1)<<24 {
		return false
	}
	ipset.sorted = false
	ipset.btree = append(ipset.btree,
		interval{Left: left, Right: right, LeftDots: dotsleft, RightDots: dotsright, Name: name, URL: url})
	return true
}

func (ipset IntervalSet) Len() int {
	return ipset.btree.Len()
}

func (ipset IntervalSet) Contains(dots string) (bool, error) {
	if !ipset.sorted {
		err := ipset.sort()
		if err != nil {
			return false, err
		}
	}

	val := dots2uint32(dots)
	if val == 0 && dots != "0.0.0.0" {
		return false, fmt.Errorf("Invalid input: %q", dots)
	}
	i := sort.Search(len(ipset.btree), func(i int) bool {
		return ipset.btree[i].Left >= val
	})
	if i < ipset.btree.Len() && ipset.btree[i].Left >= val && ipset.btree[i].Right <= val {
		return true, nil
	}
	return false, nil
}

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
