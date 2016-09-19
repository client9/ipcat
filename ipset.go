package ipcat

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
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

// CIDR2Range converts a CIDR to a dotted IP address pair, or empty strings and error
//  generic.. does not care if ipv4 or ipv6
func CIDR2Range(c string) (string, string, error) {
	left, ipnet, err := net.ParseCIDR(c)
	if err != nil {
		return "", "", err
	}
	left4 := left.To4()
	if left4 == nil {
		return "", "", nil
	}
	right := net.IPv4(0, 0, 0, 0).To4()
	right[0] = left4[0] | ^ipnet.Mask[0]
	right[1] = left4[1] | ^ipnet.Mask[1]
	right[2] = left4[2] | ^ipnet.Mask[2]
	right[3] = left4[3] | ^ipnet.Mask[3]

	return left4.String(), right.To4().String(), nil
}

// ToDots converts a uint32 to a IPv4 Dotted notation
func ToDots(val uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		val>>24,
		(val>>16)&0xFF,
		(val>>8)&0xFF,
		val&0xFF)
}

// Interval is a closed interval [a,b] of an IPv4 range
type Interval struct {
	Left      uint32
	Right     uint32
	LeftDots  string
	RightDots string
	Name      string
	URL       string
}

type intervallist []Interval

// Len satisfies the sort.Sortable interface
func (ipset intervallist) Len() int {
	return len(ipset)
}

// Less satisfies the sort.Sortable interface
func (ipset intervallist) Less(i, j int) bool {
	return ipset[i].Left < ipset[j].Left
}

// Swap satisfies the sort.Sortable interface
func (ipset intervallist) Swap(i, j int) {
	ipset[i], ipset[j] = ipset[j], ipset[i]
}

// IntervalSet is a mapping of an IP range (the closed interval)
//  to additional data
type IntervalSet struct {
	btree  intervallist
	sorted bool
}

// NewIntervalSet creates a new set with a capacity
func NewIntervalSet(capacity int) *IntervalSet {
	return &IntervalSet{
		btree: make([]Interval, 0, capacity),
	}
}

// ImportCSV imports data from a CSV file
func (ipset *IntervalSet) ImportCSV(in io.Reader) error {
	ipset.btree = nil
	ipset.sorted = false
	line := 0
	r := csv.NewReader(in)
	for {
		line++
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) != 4 {
			return fmt.Errorf("line %d: expected 4 records but got %d %v", line, len(record), record)
		}
		if err = ipset.AddRange(record[0], record[1], record[2], record[3]); err != nil {
			return err
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

func (ipset *IntervalSet) sort() error {
	if ipset.sorted {
		return nil
	}
	sort.Sort(ipset.btree)

	last := Interval{}
	// check validity -- probably worth ripping out
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
				return fmt.Errorf("Overlapping regions %v vs. %v", last, val)
			}
		}
		last = val
	}
	ipset.sorted = true

	// now merge adjacent items
	newtree := make([]Interval, 0, len(ipset.btree))
	last = Interval{}
	for pos, val := range ipset.btree {
		if pos == 0 {
			newtree = append(newtree, val)
			last = val
			continue
		}
		if last.Right+1 == val.Left && last.Name == val.Name {
			last.Right = val.Right
			newtree[len(newtree)-1] = last
			continue
		}
		newtree = append(newtree, val)
		last = val
	}
	ipset.btree = newtree
	return nil
}

// AddCIDR adds an entry based on a CIDR range
func (ipset *IntervalSet) AddCIDR(cidr, name, url string) error {
	dotsleft, dotsright, err := CIDR2Range(cidr)
	if err != nil {
		return err
	}
	return ipset.AddRange(dotsleft, dotsright, name, url)
}

// AddRange adds an entry based on an IP range
func (ipset *IntervalSet) AddRange(dotsleft, dotsright, name, url string) error {
	left := dots2uint32(dotsleft)
	if left == 0 && dotsleft != "0.0.0.0" {
		return fmt.Errorf("Unable to convert %s", dotsleft)
	}
	right := dots2uint32(dotsright)
	if right == 0 && dotsright != "0.0.0.0" {
		return fmt.Errorf("Unable to convert %s", dotsright)
	}
	if left > right {
		return fmt.Errorf("%s > %s", dotsleft, dotsright)
	}
	if right-left >= uint32(1)<<24 {
		return fmt.Errorf("Range too big for [%s %s] %s %s", dotsleft, dotsright, name, url)
	}
	ipset.sorted = false
	ipset.btree = append(ipset.btree,
		Interval{
			Left:      left,
			Right:     right,
			LeftDots:  dotsleft,
			RightDots: dotsright,
			Name:      name,
			URL:       url,
		},
	)
	return nil
}

// DeleteByName deletes all entries with the given name
func (ipset *IntervalSet) DeleteByName(name string) {
	newlist := intervallist{}
	for _, entry := range ipset.btree {
		if entry.Name != name {
			newlist = append(newlist, entry)
		}
	}
	ipset.btree = newlist
}

// Len returns the number of elements in the set
func (ipset IntervalSet) Len() int {
	return ipset.btree.Len()
}

// Contains returns the internal record if the IP address is in some
// interval else nil or error.  It returns a pointer to the internal
// record, so be careful.
func (ipset IntervalSet) Contains(dots string) (*Interval, error) {
	if !ipset.sorted {
		err := ipset.sort()
		if err != nil {
			return nil, err
		}
	}

	val := dots2uint32(dots)
	if val == 0 && dots != "0.0.0.0" {
		return nil, fmt.Errorf("Invalid input: %q", dots)
	}
	i := sort.Search(len(ipset.btree), func(i int) bool {
		return ipset.btree[i].Left >= val
	})

	// lots of cases in the lookup here.
	// if exactly equals, then compare with [i]
	if i < ipset.Len() && ipset.btree[i].Left == val && val <= ipset.btree[i].Right {
		return &ipset.btree[i], nil
	}

	// ok then it's the record before
	i--
	if i >= 0 && ipset.btree[i].Left < val && val <= ipset.btree[i].Right {
		return &ipset.btree[i], nil
	}
	return nil, nil
}

// NameSize is a tuple mapping name with a size
type NameSize struct {
	Name string
	Size int
}

// NameSizeList is a list of NameSize
type NameSizeList []NameSize

// Len satisfies the sort.Sortable interface
func (list NameSizeList) Len() int {
	return len(list)
}

// Less satisfies the sort.Sortable interface
// THIS IS DESCENDING SORT, the sign is flipped
//  MORE = FIRST
func (list NameSizeList) Less(i, j int) bool {
	return list[i].Size > list[j].Size
}

// Swap satisfies the sort.Sortable interface
func (list NameSizeList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

// RankBySize returns a list ISP and how many IPs they have
//  From this it's easy to compute
//    * Lastest providers
//    * Number of providers
//    * Total number IPs address
//
func (ipset IntervalSet) RankBySize() NameSizeList {
	counts := make(map[string]int, ipset.Len())
	for _, val := range ipset.btree {
		counts[val.Name] += int(val.Right-val.Left) + 1
	}
	rank := make(NameSizeList, 0, len(counts))
	for k, v := range counts {
		rank = append(rank, NameSize{k, v})
	}
	sort.Sort(rank)
	return rank
}
