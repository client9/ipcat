package ipcat

import (
	"net"
	"strings"
)

func lookupSPFRecord(name string, f func(dir string) error) error {
	txt, err := net.LookupTXT(name)
	if err != nil {
		return err
	}

	for _, rec := range txt {
		spf := strings.TrimPrefix(rec, "v=spf1 ")
		if spf == rec {
			continue
		}

		for _, dir := range strings.Split(spf, " ") {
			if err := f(dir); err != nil {
				return err
			}
		}
	}

	return nil
}

// DownloadAppEngine downloads and returns raw bytes of the Google App Engine ip
// range list
func DownloadAppEngine() ([]string, error) {
	ranges := []string{}
	domainList := []string{"_cloud-netblocks.googleusercontent.com"}

	for len(domainList) > 0 {
		// Dequeue a domain from this list
		var domain string
		domain, domainList = domainList[0], domainList[1:]
		err := lookupSPFRecord(domain, func(dir string) error {
			// Enqueue domain from this record
			if inc := strings.TrimPrefix(dir, "include:"); dir != inc {
				domainList = append(domainList, inc)
			}
			// Add IPv4 range
			if ip4 := strings.TrimPrefix(dir, "ip4:"); dir != ip4 {
				ranges = append(ranges, ip4)
			}
			// Add IPv6 range
			if ip6 := strings.TrimPrefix(dir, "ip6:"); dir != ip6 {
				ranges = append(ranges, ip6)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return ranges, nil
}

// UpdateAppEngine takes a raw data, parses it and updates the ipmap
func UpdateAppEngine(ipmap *IntervalSet, ranges []string) error {
	const (
		dcName = "Google App Engine"
		dcURL  = "https://cloud.google.com/appengine"
	)

	// delete all existing records
	ipmap.DeleteByName(dcName)

	for _, ipRange := range ranges {
		err := ipmap.AddCIDR(ipRange, dcName, dcURL)
		if err != nil {
			return err
		}
	}

	return nil
}
