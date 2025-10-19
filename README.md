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
- **Parallel processing with atomic operations**: Lock-free bitmap updates using CAS on uint64, 15-18% faster on large files (10M+ addresses)
- **Simple input format**: Expects valid IPv4 addresses, one per line (`\n` separated)

Attempted optimizations that didn't deliver:

- **Roaring bitmap** with various segment sizes: Added complexity without performance gains. For randomly distributed IPs, created too many small segments with excessive allocations, making it slower than simple bitmap. Simple strategy switching (map → bitmap based on dataset size) proved more effective
- **Bitmap element size** (uint8 vs uint32 vs uint64): No measurable difference between implementations
