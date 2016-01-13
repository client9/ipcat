
all: generate

generate:
	go run golang/ipset.go < datacenters.csv > tmp.txt
	cp tmp.txt datacenters.csv
	./makestats.py < datacenters.csv > README.md

README.md: makestats.py datacenters.csv
	./makestats.py < datacenters.csv > README.md

clean:
	rm -f *~
