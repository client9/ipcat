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
    newrow = {
        '_ip0': unpack("!L", inet_aton(parts[0]))[0],
        '_ip1': unpack("!L", inet_aton(parts[1]))[0],
        'owner': parts[3],
    }
    iplist[newrow['_ip0']] = newrow

#return the list of entries, sorted by the lowest ip in the range
iplist = [v for (k,v) in sorted(iplist.iteritems(), key=itemgetter(0))]

#autogenerate the class to perform lookups
print """

#include <stdint.h>

typedef struct ip_range_owner {
  uint32_t ip0;
  uint32_t ip1;
  const char* owner;
} ip_range_owner_t;

ip_range_owner_t iplist[] = {
"""

for val in iplist:
    print '{{ {0},{1},"{2}" }},'.format(val['_ip0'], val['_ip1'], val['owner'])

print "};"

print "const int iplist_len = {0};".format(len(iplist));

print "int main() { return 0; }"
