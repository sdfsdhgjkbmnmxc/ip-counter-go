.PHONY: build test benchmark testdata clean

build:
	go build -o bin/ip-counter ./cmd/ip-counter
	go build -o bin/generate-testdata ./cmd/generate-testdata

testdata: build
	@mkdir -p testdata
	@[ -f testdata/ips_100.txt ] || ./bin/generate-testdata 100 > testdata/ips_100.txt
	@[ -f testdata/ips_1k.txt ] || ./bin/generate-testdata 1000 > testdata/ips_1k.txt
	@[ -f testdata/ips_10k.txt ] || ./bin/generate-testdata 10000 > testdata/ips_10k.txt

test:
	go test -v ./...

benchmark: testdata
	go test -bench=. -benchmem ./internal/...

clean:
	rm -rf bin testdata
