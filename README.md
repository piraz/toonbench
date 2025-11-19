# toontest

`toontest` is a Go benchmarking project from [piraz](https://github.com/piraz)
designed to compare serialization and deserialization performance across
different formats: Toon, Gotoon, JSON, and Protobuf.

## What does it do?

This repository provides microbenchmarks for marshaling (serializing) and
unmarshaling (deserializing) a moderate-sized payload. It evaluates:
- [Toon](https://github.com/toon-format/toon-go)
- [Gotoon](https://github.com/someuser/gotoon) (if installed)
- `encoding/json` (standard Go library)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)

The library benchmarks execution time, memory usage, and number of heap
allocations for each encoder/decoder.

## How to Run

First, clone this repository:

```sh
git clone https://github.com/piraz/toontest.git
cd toontest
```

Make sure you have your Go environment set up and dependencies installed.
You’ll also need `protoc` and its Go plugins if you want to regenerate protobuf
files.

### Benchmark

To run the full suite of benchmarks for all serialization formats:

```sh
make bench
```

This will output a table showing how many times each benchmark ran, the average
time per operation, memory used, and heap allocations.

### Quick (Short) Benchmark

For a faster, less statistically deep run (useful for development feedback):

```sh
make bench-short
```
This runs each benchmark only twice (see `-benchtime=2x` in the Makefile), so
results will be rough but quick.

## Adding New Benchmarks

Benchmarks are implemented in standard Go `*_test.go` files under `cmd/tut/`.  
To add another format, simply create a new `Benchmark...` function.

## Example Output

```
>>> make bench
go test -bench=. -benchmem ./cmd/tut
goos: linux
goarch: amd64
pkg: github.com/piraz/toontest/cmd/tut
cpu: AMD Ryzen 5 1600 Six-Core Processor
BenchmarkToonMarshal-12                7         149557411 ns/op        48315452 B/op     999951 allocs/op
BenchmarkToonUnmarshal-12              6         172380953 ns/op        69058568 B/op    1248966 allocs/op
BenchmarkGotoonMarshal-12              3         337262098 ns/op        72150930 B/op    1099958 allocs/op
BenchmarkJsonMarshal-12               50          25097259 ns/op         5963423 B/op          5 allocs/op
BenchmarkJsonUnmarshal-12              8         129124964 ns/op         5171892 B/op     200010 allocs/op
BenchmarkProtoMarshal-12              80          13132078 ns/op         2654402 B/op          2 allocs/op
BenchmarkProtoUnmarshal-12            25          45326514 ns/op        14892208 B/op     300028 allocs/op
PASS
ok      github.com/piraz/toontest/cmd/tut       7.981s
```

## Requirements

- Go 1.20 or higher
- (For Protobuf) `protoc` and [protoc-gen-go](https://pkg.go.dev/mod/google.golang.org/protobuf/cmd/protoc-gen-go)

## License

MIT
