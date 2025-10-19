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

- **Naive**: String-based map without IP parsing (research/comparison)
- **Map**: Hash map, memory-efficient for smaller datasets (< 28M addresses, see `internal/u32/u32_test.go:TestMemoryCrossover`)
- **Bitmap**: Sequential bitmap, fixed 512 MB memory
- **ParallelBitmap**: Parallel bitmap with atomic operations, 3x faster than sequential Bitmap
- **Auto** (default): Automatically selects between Map and ParallelBitmap based on file size

## Optimizations

For maximum performance:

- **Memory-mapped I/O**: Reads files via `mmap` instead of standard I/O (stdin not supported)
- **Custom IPv4 parser**: Specialized parser converts addresses to `uint32`, significantly faster than standard library
- **Parallel processing with atomic operations**: Lock-free bitmap updates using CAS on uint64. Each worker maintains local counters to avoid contention, achieving near-linear scaling
- **Simple input format**: Expects valid IPv4 addresses, one per line (`\n` separated)

Attempted optimizations that didn't deliver:

- **Roaring bitmap** with various segment sizes: Added complexity without performance gains. For randomly distributed IPs, created too many small segments with excessive allocations, making it slower than simple bitmap. Simple strategy switching (map → bitmap based on dataset size) proved more effective
- **Bitmap element size** (uint8 vs uint32 vs uint64): No measurable difference between implementations

### Performance

Benchmark results (AMD EPYC 7763, 4 workers) - Throughput in K IP/sec:

| Method | 1K | 10K | 100K | 1M | 10M | 100M |
|--------|-----|-----|------|-----|------|------|
| Naive | 6,791 | 7,278 | 6,183 | 3,601 | 2,437 | 2,149 |
| Map | **12,345** | **14,440** | **12,735** | 10,815 | 5,609 | 4,646 |
| Bitmap | 37 | 347 | 2,567 | 6,637 | 7,840 | 8,068 |
| ParallelBitmap | 37 | 363 | 3,212 | **15,127** | **23,434** | **25,444** |
