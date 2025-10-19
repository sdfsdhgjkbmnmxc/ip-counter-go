.PHONY: test bench testdata testdata-large

sample:
	go run cmd/ip-counter/main.go testdata/ips_10m.txt -method=ParallelBitmap

test: testdata
	go test -v ./internal/...

bench: testdata
	go test -bench=BenchmarkCounter -benchtime=8s ./internal/...

testdata:
	mkdir -p testdata
	[ -f testdata/ips_1k.txt ]   || go run ./cmd/generate-testdata 1000 > testdata/ips_1k.txt
	[ -f testdata/ips_10k.txt ]  || go run ./cmd/generate-testdata 10000 > testdata/ips_10k.txt
	[ -f testdata/ips_100k.txt ] || go run ./cmd/generate-testdata 100000 > testdata/ips_100k.txt
	[ -f testdata/ips_1m.txt ]   || go run ./cmd/generate-testdata 1000000 > testdata/ips_1m.txt
	[ -f testdata/ips_10m.txt ]  || go run ./cmd/generate-testdata 10000000 > testdata/ips_10m.txt
	[ -f testdata/ips_100m.txt ] || go run ./cmd/generate-testdata 100000000 > testdata/ips_100m.txt

testdata-large:
	mkdir -p testdata
	echo "Downloading large test file..."
	wget -q https://ecwid-vgv-storage.s3.eu-central-1.amazonaws.com/ip_addresses.zip
	echo "Extracting..."
	unzip -q ip_addresses.zip -d testdata/
	rm ip_addresses.zip
	ls -lh testdata/
