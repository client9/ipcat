#!/usr/bin/env python

"""Generate the ipdb for Golang."""

import pprint
from operator import itemgetter
from socket import inet_aton
from struct import unpack
from urllib import urlopen
import sys

pp = pprint.PrettyPrinter(indent=4, width=50)
iplist = {}


if len(sys.argv) > 1:
    fp = open(sys.argv[1], 'r')
else:
    # fetch remote datacenter list and convert to searchable datastructure
    external_list = \
        'https://raw.github.com/client9/ipcat/master/datacenters.csv'
    fp = urlopen(external_list)

fp.readline()
for line in fp:
    line = line.strip()
    if not line or line[0] == '#':
        continue

    parts = line.split(",")
    newrow = (
        unpack("!L", inet_aton(parts[0]))[0],
        unpack("!L", inet_aton(parts[1]))[0],
        parts[3],
    )
    iplist[newrow[0]] = newrow

# return the list of entries, sorted by the lowest ip in the range
iplist = tuple(v for (k, v) in sorted(iplist.iteritems(),
                                      key=itemgetter(0)))

values = ""
for ip in iplist:
    values += '\n\tIpRange{From: %d, To: %d, Service: `%s`},' % (ip[0], ip[1], ip[2])

# autogenerate the class to perform lookups
print """// ipcatmodel.go
package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

type IpRange struct {
	From    uint32
	To      uint32
	Service string
}

type IpRanges []IpRange

var ipranges = IpRanges{%s
}
var iprLen = len(ipranges)

func Ip2long(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}

func IpFind(ipaddress string) *IpRange {

	if ip, err := Ip2long(ipaddress); err == nil {
		low := 0
		high := iprLen - 1
		for high >= low {
			probe := (high + low) / 2
			ir := ipranges[probe]
			if ir.From > ip {
				high = probe - 1
			} else if ir.To < ip {
				low = probe + 1
			} else {
				return &ir
			}
		}
	}
	return nil
}

func main() {
	fmt.Println(time.Now())
	fmt.Printf("ipranges: %%+v\\n", IpFind("5.10.64.0"))
	fmt.Println(time.Now())
}
""" % values
