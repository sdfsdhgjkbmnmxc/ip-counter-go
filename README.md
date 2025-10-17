![Benchmarks](https://github.com/sdfsdhgjkbmnmxc/ip-counter-go/workflows/Benchmarks/badge.svg)

# IP Address Counter

Counts unique IPv4 addresses in text files. Works like `sort -u | wc -l ip_addresses.txt` but optimized for performance and memory usage.

## Installation

```bash
go install github.com/sdfsdhgjkbmnmxc/ip-counter-go/cmd/ip-counter@latest
```

## Usage

```bash
# Default (ComboSet - automatic strategy selection)
ip-counter ip_addresses.txt

# Specify method explicitly
ip-counter -method=MapSet ip_addresses.txt
ip-counter -method=BitmapSet ip_addresses.txt
```

### Available Methods

Use `-method` flag to select implementation:

- **ComboSet** (default): Automatically switches between MapSet and BitmapSet based on file size
  - Uses **MapSet** for < 28M addresses (memory-efficient)
  - Switches to **BitmapSet** for larger files (caps memory at 512 MB)
  - Threshold determined experimentally (see `internal/u32/u32_test.go:TestMemoryCrossover`)
- **MapSet**: Hash map, memory-efficient for smaller datasets
- **BitmapSet**: Bitmap, fast on large datasets with fixed 512 MB memory
- **NaiveCounter**: String-based map without IP parsing (research/comparison)

## Optimizations

For maximum performance:

- **Memory-mapped I/O**: Reads files via `mmap` instead of standard I/O (stdin not supported)
- **Custom IPv4 parser**: Specialized parser converts addresses to `uint32`, significantly faster than standard library
- **Simple input format**: Expects valid IPv4 addresses, one per line (`\n` separated)

Benchmarks:

```
cpu: AMD EPYC 7763 64-Core Processor                
BenchmarkCounter/ips_1k.txt/MapSet-4  	  134082	     71154 ns/op	   19424 B/op	      11 allocs/op
BenchmarkCounter/ips_1k.txt/BitmapSet-4         	     368	  25763915 ns/op	536871299 B/op	       6 allocs/op
BenchmarkCounter/ips_1k.txt/ComboSet-4          	  133092	     71275 ns/op	   19424 B/op	      11 allocs/op
BenchmarkCounter/ips_1k.txt/NaiveCounter-4      	   68467	    140081 ns/op	  128985 B/op	    1024 allocs/op
BenchmarkCounter/ips_10k.txt/MapSet-4           	   15548	    618999 ns/op	  152593 B/op	      39 allocs/op
BenchmarkCounter/ips_10k.txt/BitmapSet-4        	     345	  26894345 ns/op	536871295 B/op	       6 allocs/op
BenchmarkCounter/ips_10k.txt/ComboSet-4         	   15343	    624670 ns/op	  152593 B/op	      39 allocs/op
BenchmarkCounter/ips_10k.txt/NaiveCounter-4     	    7123	   1335448 ns/op	 1037489 B/op	   10083 allocs/op
BenchmarkCounter/ips_100k.txt/MapSet-4          	    1309	   7235669 ns/op	 1218188 B/op	     263 allocs/op
BenchmarkCounter/ips_100k.txt/BitmapSet-4       	     255	  37688670 ns/op	536871292 B/op	       6 allocs/op
BenchmarkCounter/ips_100k.txt/ComboSet-4        	    1336	   7161134 ns/op	 1218137 B/op	     263 allocs/op
BenchmarkCounter/ips_100k.txt/NaiveCounter-4    	     600	  15975669 ns/op	 8593775 B/op	  100534 allocs/op
BenchmarkCounter/ips_1m.txt/MapSet-4            	      78	 107245563 ns/op	19483062 B/op	    4103 allocs/op
BenchmarkCounter/ips_1m.txt/BitmapSet-4         	      60	 145867462 ns/op	536871291 B/op	       6 allocs/op
BenchmarkCounter/ips_1m.txt/ComboSet-4          	      88	 105277180 ns/op	19483056 B/op	    4103 allocs/op
BenchmarkCounter/ips_1m.txt/NaiveCounter-4      	      28	 314847159 ns/op	127540705 B/op	 1008195 allocs/op
BenchmarkCounter/ips_10m.txt/MapSet-4           	       5	1825795596 ns/op	155845027 B/op	   32775 allocs/op
BenchmarkCounter/ips_10m.txt/BitmapSet-4        	       7	1257546789 ns/op	536871288 B/op	       6 allocs/op
BenchmarkCounter/ips_10m.txt/ComboSet-4         	       5	1800368941 ns/op	155845065 B/op	   32775 allocs/op
BenchmarkCounter/ips_10m.txt/NaiveCounter-4     	       2	4309922776 ns/op	1054699056 B/op	10065565 allocs/op
BenchmarkCounter/ips_100m.txt/MapSet-4          	       1	21315049036 ns/op	1246757264 B/op	  262151 allocs/op
BenchmarkCounter/ips_100m.txt/BitmapSet-4       	       1	12165125514 ns/op	536871288 B/op	       6 allocs/op
BenchmarkCounter/ips_100m.txt/ComboSet-4        	       1	12081513940 ns/op	536871288 B/op	       6 allocs/op
BenchmarkCounter/ips_100m.txt/NaiveCounter-4    	       1	46508692824 ns/op	8757519072 B/op	100524322 allocs/op
```
