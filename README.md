**ipcat**: datasets for categorizing IP addresses.

Archived in 2023. Please fork and edit as you wish. It's MIT now.  Onward -- nickg

---

This is a list of IPv4 addresses that correspond to datacenters,
co-location centers, shared and virtual webhosting providers.  In
other words, ip addresses that end web consumers should not be using.

Statistics
------------------------

Check out the new [datacenter stats](/datacenters-stats.csv)

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
an easy way of keeping track of changes.

How is this generated?
-------------------------

Manually from users like you, and automatically via proprietary
discovery algorithms.

Who made this?
-------------------------

Nick Galbreath.  See more at http://www.client9.com/

