#!/usr/bin/env python
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
from socket import inet_aton
from struct import unpack
from math import floor
class IpDb(object):
    iplist = %s

    @staticmethod
    def find(ipstring):
        ip = unpack("!L", inet_aton(ipstring))[0]

        high = len(IpDb.iplist)-1
        low = 0
        while high >= low:
            probe = int(floor((high+low)/2))
            if IpDb.iplist[probe]['_ip0'] > ip:
                high = probe - 1
            elif IpDb.iplist[probe]['_ip1'] < ip:
                low = probe + 1
            else:
                return IpDb.iplist[probe]
        return None
""" % (pp.pformat(iplist), )