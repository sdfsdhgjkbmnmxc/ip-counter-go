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

## Optimizations

For maximum performance:

- **Memory-mapped I/O**: Reads files via `mmap` instead of standard I/O (stdin not supported)
- **Custom IPv4 parser**: Specialized parser converts addresses to `uint32`, significantly faster than standard library
- **Parallel processing with atomic operations**: Lock-free bitmap updates using CAS on uint64, 15-18% faster on large files (10M+ addresses)
- **Simple input format**: Expects valid IPv4 addresses, one per line (`\n` separated)

Attempted optimizations that didn't deliver:

- **Roaring bitmap** with various segment sizes: Added complexity without performance gains. For randomly distributed IPs, created too many small segments with excessive allocations, making it slower than simple bitmap. Simple strategy switching (map → bitmap based on dataset size) proved more effective
- **Bitmap element size** (uint8 vs uint32 vs uint64): No measurable difference between implementations

### Performance

Benchmark results (AMD EPYC 7763, 4 workers) - Throughput in K IP/sec:

| Method | 1K | 10K | 100K | 1M | 10M | 100M |
|--------|-----|-----|------|-----|------|------|
| Naive | 6,950 | 7,440 | 6,240 | 3,440 | 2,370 | 2,150 |
| Map | **13,870** | **16,030** | **13,830** | 10,660 | 5,450 | 4,700 |
| Bitmap | 38 | 367 | 2,620 | 6,860 | 8,140 | 8,240 |
| ParallelBitmap | 38 | 370 | 3,030 | **10,140** | **13,360** | **13,830** |
