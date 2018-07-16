#from subprocess import _args_from_interpreter_flags
#from multiprocessing.dummy import Lock
from threading import Lock
from operator import itemgetter
from datetime import datetime
from struct import unpack
from math import floor
from socket import inet_aton
from urllib import urlopen

class IpDatabase(object):
# some very hand-wavy benchmarking:
#	DB generation takes around 1.6 seconds on average
#	Find is fast and stable: hits and misses in under 1/20,000 second, every time
#	DB update is 50,000x slower than DB seek, so we will set sensible limits
#	on the cachelimit param so as to not constantly refresh and bottleneck seeks

	ext_list = "https://raw.github.com/client9/ipcat/master/datacenters.csv"
	def __init__(self):
		self.db = [] 
		self.refreshed_at = 0
		self.mutex = Lock()
	
	@classmethod
	def generate(cls):
		c = cls()
		c.update()
		return c

	def update(self):
		# larger mem footprint (215K) but less time in lock means shorter wait for competing reqs
		csv = urlopen(self.ext_list)		
		csv.readline()
		tdb = [] 
		for line in csv:
			line = line.strip()
			if not line or line[0] == '#':
				continue
			cells = line.split(",")
			row = (
				unpack("!L", inet_aton(cells[0]))[0],
				unpack("!L", inet_aton(cells[1]))[0],
				cells[2], cells[3],
			)
			tdb.append(row)
			tdb.sort(key=itemgetter(0))
                self.mutex.acquire()
                self.db = tdb
                self.refreshed_at = datetime.now()
                self.mutex.release()

	def needs_update(self, cachelimit):
		# compare cachelimit to elapsed_since_refresh and return true / false
		elapsed = datetime.now() - self.refreshed_at
		return True if elapsed.total_seconds() > cachelimit else False
	
	def find(self, ip):
		# convert ip to int and peform a binary search of db
		ip_int = unpack("!L", inet_aton(ip))[0]
		hlimit = len(self.db) - 1
		llimit = 0
                self.mutex.acquire()
		while hlimit >= llimit:
			midpoint = int(floor((hlimit + llimit) / 2))
			if self.db[midpoint][0] > ip_int:
				hlimit = midpoint - 1
			elif self.db[midpoint][1] < ip_int:
				llimit = midpoint + 1
			else:
                                self.mutex.release()
				return self.db[midpoint]
                self.mutex.release()
		return None

