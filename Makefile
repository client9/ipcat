
all: generate

generate:
	go run ./cmd/ipcat/main.go
	./makestats.py < datacenters.csv > README.md

aws:
	go run ./cmd/ipcat/main.go -aws
	./makestats.py < datacenters.csv > README.md

README.md: makestats.py datacenters.csv
	./makestats.py < datacenters.csv > README.md

clean:
	rm -f *~
