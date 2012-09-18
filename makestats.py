#!/usr/bin/env python
"""
validates file and produces the README.md file with stats, or pops an exception
"""

import sys
import socket
import struct
from collections import Counter

def ip2int(ip):
    return struct.unpack('>I',socket.inet_aton(ip))[0]

rows = 0
counts = Counter()
total = 0
lastip0 = 0
lastip1 = 0
for line in sys.stdin:
    line = line.strip()
    if len(line) == 0 or line[0] == '#':
        continue
    rows += 1
    parts = line.split(',')
    if len(parts) != 4:
        raise Exception("Line %d has more than 4 entries: %s" % (row, line))

    (dots0,dots1,name,url) = parts
    ip0 = ip2int(dots0)
    ip1 = ip2int(dots1)
    if ip0 > ip1:
        raise Exception("Line %d has starting IP > ending IP: %s" % (row, line))

    if ip0 <= lastip1:
        raise Exception("Line %d is not sorted: %s" % (row, line))

    # we are correct
    lastip0 = ip0
    lastip1 = ip1

    sz = ip1 - ip0
    total += sz
    counts[name] += sz

print("""ipcat: datasets for categorizing IP addresses.

The first release "datacenters.csv" is focusing
on IPv4 address that correspond to datacenters, co-location centers,
shared and virtual webhosting providers.  In other words, ip addresses
that end web consumers should not be using.

Licensing -- GPL v3
------------------------

The data is licensed under GPL v3, see COPYING for details.

Relaxations and commericial licensing are gladly available by request.
The use of GPL is to prevent commercial data providers from scooping up
this data without compensation or attribution.

This may be changed to another less restrictive license later.

Statistics
------------------------

<table>
<tr><th>IPs</th><td>%d</td></tr>
<tr><th>Records</th><td>%d</td></tr>
<tr><th>ISPs</th><td>%d</td></tr>
</table>

What is the file format?
-------------------------

Standard CSV with ip-start, ip-end (inclusive, in dot-notation), name of provider, url
of provider.  IP ranges are non-overlapping, and in sorted order.

Why is hosting provider XXX is missing?
---------------------------------------

It might not be.  Many providers are resellers of another and will be
included under another name or ip range.

Also, as of 16-Oct-2011, many locations from Africa, Latin
America, Korea and Japan are missing.

Or, it might just be missing.  Please let us know!

Why GitHub + CSV?
-------------------------

The goal of the file format and the use of github was designed to make
it really easy for other to send patches or additions.  It also provides
and easy way of keeping track of changes.

How is this generated?
-------------------------

Manually from users like you, and automatically via proprietary
discovery algorithms.

Who made this?
-------------------------

Nick Galbreath.  See more at http://www.client9.com/

""" % (total, rows, len(counts)))



