![Benchmarks](https://github.com/sdfsdhgjkbmnmxc/ip-counter-go/workflows/Benchmarks/badge.svg)

# IP Address Counter

Counts unique IPv4 addresses in text files. Works like `sort -u | wc -l ip_addresses.txt` but optimized for performance and memory usage.

## Installation

```bash
go install github.com/sdfsdhgjkbmnmxc/ip-counter-go/cmd/ip-counter@latest
```

## Usage

```bash
# Default (Auto - automatic strategy selection)
ip-counter ip_addresses.txt

# Specify method explicitly
ip-counter -method=Map ip_addresses.txt
ip-counter -method=Bitmap ip_addresses.txt
ip-counter -method=ParallelBitmap ip_addresses.txt
```

### Available Methods

Use `-method` flag to select implementation:

- **Auto** (default): Automatically selects best strategy based on file size
  - Uses **Map** for < 28M addresses (memory-efficient)
  - Switches to **ParallelBitmap** for larger files (parallel processing)
  - Threshold determined experimentally (see `internal/u32/u32_test.go:TestMemoryCrossover`)
- **Map**: Hash map, memory-efficient for smaller datasets
- **Bitmap**: Sequential bitmap, fixed 512 MB memory
- **ParallelBitmap**: Parallel bitmap with atomic operations, 15-18% faster on large files (10M+ addresses)
- **Naive**: String-based map without IP parsing (research/comparison)

## Performance

Benchmark on 100M unique IP addresses (Apple M1 Pro, 4 workers):

![Performance Benchmark](https://quickchart.io/chart?c={type:'bar',data:{labels:['ParallelBitmap','Auto','Bitmap','Map','Naive'],datasets:[{label:'Time (seconds)',data:[7.2,7.2,12.1,21.3,46.6],backgroundColor:['%2322c55e','%2322c55e','%2360a5fa','%2360a5fa','%23ef4444']}]},options:{indexAxis:'y',plugins:{title:{display:true,text:'100M IP Addresses (lower is better)'},legend:{display:false}},scales:{x:{title:{display:true,text:'Time (seconds)'}}}}}})

**ParallelBitmap** achieves 40% speedup over sequential Bitmap through lock-free atomic operations.

## Optimizations

For maximum performance:

- **Memory-mapped I/O**: Reads files via `mmap` instead of standard I/O (stdin not supported)
- **Custom IPv4 parser**: Specialized parser converts addresses to `uint32`, significantly faster than standard library
- **Parallel processing with atomic operations**: Lock-free bitmap updates using CAS on uint64, 15-18% faster on large files (10M+ addresses)
- **Simple input format**: Expects valid IPv4 addresses, one per line (`\n` separated)

Attempted optimizations that didn't deliver:

- **Roaring bitmap** with various segment sizes: Added complexity without performance gains. For randomly distributed IPs, created too many small segments with excessive allocations, making it slower than simple bitmap. Simple strategy switching (map → bitmap based on dataset size) proved more effective
- **Bitmap element size** (uint8 vs uint32 vs uint64): No measurable difference between implementations
