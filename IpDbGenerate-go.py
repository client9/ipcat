#!/usr/bin/env python

#
# Generates a C-struct
#  Doesn't generate the lookup bsearch algorithm
#

import pprint
from operator import itemgetter
from socket import inet_aton
from struct import unpack
from urllib import urlopen

pp = pprint.PrettyPrinter(indent=4, width=50)
iplist = {}

#fetch remote datacenter list and convert to searchable datastructure
external_list = 'https://raw.github.com/client9/ipcat/master/datacenters.csv'
fp = urlopen(external_list)
for line in fp:
    line = line.strip()
    if not line or line[0] == '#':
        continue

    parts = line.split(",")

    if parts[0] == "start":
        continue

    ipint = unpack("!L", inet_aton(parts[0]))[0],
    newrow = {
        '_ip0': parts[0],
        '_ip1': parts[1],
        'owner': parts[3],
    }
    iplist[ipint] = newrow

#return the list of entries, sorted by the lowest ip in the range
iplist = [v for (k,v) in sorted(iplist.iteritems(), key=itemgetter(0))]

#autogenerate the class to perform lookups
print """
package ipdb

import "net"
import "bytes"
import "math"

type IpRangeOwner struct {
  Ip0 net.IP
  Ip1 net.IP
  Owner string
}
"""

print "const ipDbLen = {0};".format(len(iplist));

print """
func Find(ipStr string) bool { 
    ip := net.ParseIP(ipStr)
    if ip != nil {
        high := ipDbLen - 1
        low := 0
        for high >= low {
            probe := int(math.Floor(float64(high+low)/2))
            ipRange := ipDb[probe]

            // ip0 is greater
            if bytes.Compare(ipRange.Ip0, ip) == 1 {
                high = probe - 1
            } else if bytes.Compare(ipRange.Ip1, ip) == -1 {
                low = probe + 1
            } else {
                return true
            }
        }
    }
    return false
}

var ipDb = []IpRangeOwner{
"""

for val in iplist:
    print 'IpRangeOwner{{ net.ParseIP("{0}"), net.ParseIP("{1}"),"{2}" }},'.format(val['_ip0'], val['_ip1'], val['owner'])

print "}"
