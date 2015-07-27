
README.md: makestats.py datacenters.csv
	./makestats.py < datacenters.csv > README.md

clean:
	rm -f *~
