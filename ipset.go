package ipcat

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"math/big"
	"net"
	"sort"
	"strings"
)

// IPParse converts a string IP address to a byte slice, or nil on error.
func IPParse(dots string) []byte {
	return net.ParseIP(dots)
}

// IPString converts a byte slice representing an IP address to a string.
func IPString(ip []byte) string {
	return net.IP(ip).String()
}

// IPIncrementEquals returns true if the first IP + 1 equals the second.
func IPIncrementEquals(bytes, bytesinc []byte) bool {
	if len(bytes) != len(bytesinc) {
		return false
	}

	// Iterate backwards
	var carry byte = 1
	for i := len(bytes) - 1; i >= 0; i-- {
		inc := bytes[i] + carry
		if inc != bytesinc[i] {
			return false
		}
		if inc == 0 {
			carry = 1
		} else {
			carry = 0
		}
	}

	return true
}

// CIDR2Range converts a CIDR to a dotted IP address pair, or empty strings and error
//
// Generic.. does not care if ipv4 or ipv6 (for sure this time)
func CIDR2Range(c string) (string, string, error) {
	// Parse CIDR notation
	addr, network, err := net.ParseCIDR(c)
	if err != nil {
		return "", "", err
	}

	// Create new bounds addresses
	left := make(net.IP, 16)
	right := make(net.IP, 16)

	// Pad mask to 16 bytes with 1 bits
	mask := make([]byte, 16)
	for i := 0; i < 16; i++ {
		mask[i] = 0xff
	}
	copy(mask[16-len(network.Mask):], network.Mask)

	// Mask address for left and right bounds
	for i := range left {
		left[i] = addr[i] & mask[i]
		right[i] = addr[i] | ^mask[i]
	}

	return left.String(), right.String(), nil
}

// Interval is a closed interval [a,b] of an IP range
type Interval struct {
	Left  [16]byte
	Right [16]byte
	Name  string
	URL   string
}

// Contains returns true if the IP address is found within the interval.
func (interval *Interval) Contains(ip []byte) bool {
	return bytes.Compare(ip, interval.Left[:]) >= 0 && bytes.Compare(ip, interval.Right[:]) <= 0
}

// Size returns the number of IP addresses that fit in the range
func (interval *Interval) Size() *big.Int {
	size := big.NewInt(1)
	size.Add(size, new(big.Int).SetBytes(interval.Right[:]))
	size.Sub(size, new(big.Int).SetBytes(interval.Left[:]))
	return size
}

type intervallist []Interval

// Len satisfies the sort.Sortable interface
func (ipset intervallist) Len() int {
	return len(ipset)
}

// Less satisfies the sort.Sortable interface
func (ipset intervallist) Less(i, j int) bool {
	return bytes.Compare(ipset[i].Left[:], ipset[j].Left[:]) < 0
}

// Swap satisfies the sort.Sortable interface
func (ipset intervallist) Swap(i, j int) {
	ipset[i], ipset[j] = ipset[j], ipset[i]
}

// IntervalSet is a mapping of an IP range (the closed interval)
// to additional data
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
		rec := []string{IPString(val.Left[:]), IPString(val.Right[:]), val.Name, val.URL}
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
		if bytes.Compare(val.Left[:], val.Right[:]) > 0 {
			return fmt.Errorf("left %d > right %d at pos %d",
				val.Left, val.Right, pos)
		}
		if pos > 0 {
			if bytes.Compare(val.Left[:], last.Right[:]) <= 0 || bytes.Compare(val.Right[:], last.Right[:]) <= 0 {
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
		if last.Name == val.Name && IPIncrementEquals(last.Right[:], val.Left[:]) {
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
	left := net.ParseIP(dotsleft)
	if left == nil {
		return fmt.Errorf("Unable to convert %s", dotsleft)
	}
	right := net.ParseIP(dotsright)
	if right == nil {
		return fmt.Errorf("Unable to convert %s", dotsright)
	}
	if bytes.Compare(left, right) > 0 {
		return fmt.Errorf("%s > %s", dotsleft, dotsright)
	}

	ipset.sorted = false
	ipset.btree = append(ipset.btree,
		Interval{
			Name: name,
			URL:  url,
		},
	)

	index := len(ipset.btree) - 1
	copy(ipset.btree[index].Left[:], left)
	copy(ipset.btree[index].Right[:], right)

	return nil
}

// DeleteByName deletes all entries with the given name
func (ipset *IntervalSet) DeleteByName(name string) {
	newlist := ipset.btree[:0]
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
	if err := ipset.sort(); err != nil {
		return nil, err
	}

	ip := net.ParseIP(dots)
	if ip == nil {
		return nil, fmt.Errorf("Invalid input: %q", dots)
	}

	len := ipset.Len()

	index := sort.Search(len, func(i int) bool {
		left := ipset.btree[i].Left[:]
		cmp := bytes.Compare(left, ip)
		return cmp >= 0
	})

	if index < len {
		// lots of cases in the lookup here.
		// if exactly equals, then compare with [i]
		interval := &ipset.btree[index]
		if interval.Contains(ip) {
			return interval, nil
		}
	}

	// ok then it's the record before
	if index > 0 {
		interval := &ipset.btree[index-1]
		if interval.Contains(ip) {
			return interval, nil
		}
	}

	return nil, nil
}

// NameSize is a tuple mapping name with a size
type NameSize struct {
	Name string
	Size *big.Int
}

// NameSizeList is a list of NameSize
type NameSizeList []NameSize

type lessFunc func(p1, p2 *NameSize) bool

// multiSorter implements the Sort interface, sorting the NameSizes within.
// from https://golang.org/pkg/sort/#example__sortMultiKeys
type multiSorter struct {
	nameSizes []NameSize
	less      []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to orderedBy.
func (ms *multiSorter) Sort(nameSizes []NameSize) {
	ms.nameSizes = nameSizes
	sort.Sort(ms)
}

// orderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func orderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.nameSizes)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.nameSizes[i], ms.nameSizes[j] = ms.nameSizes[j], ms.nameSizes[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that is either Less or
// !Less. Note that it can call the less functions twice per call. We
// could change the functions to return -1, 0, 1 and reduce the
// number of calls for greater efficiency: an exercise for the reader.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.nameSizes[i], &ms.nameSizes[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

// RankBySize returns a list ISP and how many IPs they have
// From this it's easy to compute:
//
// * Lastest providers
//
// * Number of providers
//
// * Total number IP address
//
func (ipset IntervalSet) RankBySize() NameSizeList {
	counts := make(map[string]*big.Int, ipset.Len())
	for _, val := range ipset.btree {
		count, ok := counts[val.Name]
		if !ok {
			count = big.NewInt(0)
			counts[val.Name] = count
		}

		count.Add(count, val.Size())
	}
	rank := make(NameSizeList, 0, len(counts))
	for k, v := range counts {
		rank = append(rank, NameSize{k, v})
	}

	size := func(left, right *NameSize) bool {
		return left.Size.Cmp(right.Size) > 0
	}

	name := func(l1, l2 *NameSize) bool {
		return strings.ToLower(l1.Name) < strings.ToLower(l2.Name)
	}

	orderedBy(size, name).Sort(rank)
	return rank
}
