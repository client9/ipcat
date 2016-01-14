ipcat: datasets for categorizing IP addresses.

[![Go Report Card](http://goreportcard.com/badge/client9/ipcat)](http://goreportcard.com/report/client9/ipcat) [![GoDoc](https://godoc.org/github.com/client9/ipcat?status.svg)](https://godoc.org/github.com/client9/ipcat) [![Coverage](http://gocover.io/_badge/github.com/client9/ipcat)](http://gocover.io/github.com/client9/ipcat)

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
<tr><th>IPs</th><td>57003088</td></tr>
<tr><th>Records</th><td>2828</td></tr>
<tr><th>ISPs</th><td>586</td></tr>
</table>

What is the file format?
-------------------------

Standard CSV with ip-start, ip-end (inclusive, in dot-notation),
name of provider, url of provider.  IP ranges are non-overlapping,
and in sorted order.

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


