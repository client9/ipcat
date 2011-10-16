
ipcat is a dataset for categorizing IP addresses.

The first release "datacetners.csv" is focusing
on IPv4 address that correspond to a datacenters, co-location centers,
shared and virtual webhosting providers.  In other words, ip addresses
that end web consumers should not be using.

Licensing -- GPL v3
------------------------

This data is licensed under GPL v3, see COPYING for details.
Relaxations and commericial licensing are gladly available by request.
The use of GPL is to prevent commercial data providers from scoping up
this data without compensation or attribution.

This may be changed to another (less restrictive) license later.


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

Magic!
