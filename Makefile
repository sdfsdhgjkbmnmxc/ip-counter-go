.PHONY: test bench testdata clean

testdata:
	@mkdir -p testdata
	@[ -f testdata/ips_1k.txt ] || go run ./cmd/generate-testdata 1000 > testdata/ips_1k.txt
	@[ -f testdata/ips_10k.txt ] || go run ./cmd/generate-testdata 10000 > testdata/ips_10k.txt
	@[ -f testdata/ips_100k.txt ] || go run ./cmd/generate-testdata 100000 > testdata/ips_100k.txt

test:
	go test -v ./...

bench: testdata
	go test -bench=. -benchmem ./internal/...

clean:
	rm -rf testdata
