# toonbench

`toonbench` is a benchmarking project by [piraz](https://github.com/piraz) that
compares serialization and deserialization performance across multiple formats
and runtimes: Go and JavaScript.

This repository contains:
- Go microbenchmarks measuring ns/op, memory (B/op) and allocations
  (allocs/op).
- A Node.js microbenchmark script (src/bench.js) that compares JS codecs (Toon,
  protobufjs, google‑protobuf and JSON).
- Makefile helpers to generate protobuf artifacts and run benchmarks.
- A shared protobuf schema (`user.proto`) used by Go and JS implementations.

Table of contents
- What it measures
- Requirements
- Generate protobuf artifacts
- Run the benchmarks (Go and JS)
- Bench design notes
- Troubleshooting: ESM/CJS imports
- Interpreting results
- Contributing
- License

---

## What it measures

This project compares the performance of different serialization formats and
implementations in both Go and JavaScript:

- Toon (Go and JS)  
- Gotoon (Go, if present)  
- encoding/json (Go standard library)  
- Protocol Buffers:
  - Go: protoc + protoc-gen-go
  - JS: protobufjs (pbjs) and google-protobuf (protoc JS)

Benchmarks report:
- throughput (ops/sec)
- time per operation (ns/op)
- memory used (B/op) and heap allocations (allocs/op) in Go
- elapsed wall-clock and approximate memory delta in JS (when run with
  --expose-gc)

** Important note about TOON **
- TOON is an emerging, specialized format optimized for sending data to large
  language models (LLMs). It is designed to reduce token usage and improve
  payload efficiency for LLM inputs.
- TOON is not intended as a general-purpose interoperability format between
  services, servers, or heterogeneous systems. It is not (yet) a drop‑in
  replacement for formats like JSON or Protocol Buffers when you need robust
  cross-language, cross-service interoperability.
- Recommendation: use TOON when your goal is compact, token‑efficient payloads
  for LLM consumption. For inter-service interoperability, persistence, or
      public APIs, prefer established formats (JSON, Protobuf, etc.).

---

## Requirements

- Go 1.20+ for Go benchmarks
- Node.js 16+ for the JavaScript benchmark
- npm / npx (for pbjs)
- protoc (Protocol Buffers compiler)
  - `protoc-gen-go` for Go generation
- make (Makefile targets)

---

## Generate protobuf artifacts

From the repository root:

```sh
make prepare
```

`make prepare` runs the typical commands used here:

- Generate protobufjs CommonJS module (static):
```
npx pbjs -t static-module -w commonjs -o ./src/compiled_protos.cjs user.proto
```

- Generate google-protobuf JS (CommonJS style):
```
protoc --js_out=import_style=commonjs,binary:./src user.proto
```

- Generate Go protobuf code:
```
protoc --go_out=. user.proto
```

Notes:
- If your project uses `"type": "module"` in package.json, prefer generating
  CommonJS artifacts as `.cjs` and require them from ESM code using
  `createRequire`.
- If you prefer ESM output from `pbjs`, generate with `-w es6` and import using
  the explicit file extension: `import ... from './compiled_protos.js'`.

---

## Run the benchmarks

### Go benchmarks

Full suite (statistically meaningful):

```sh
make go-bench
# runs: go test -bench=. -benchmem ./gobench
```

Short (fast feedback):

```sh
make go-bench-short
# runs: go test -bench=. -benchmem -benchtime=2x ./gobench
```

These produce standard `go test` benchmark output with `ns/op`, `B/op`, and
`allocs/op`.

### JavaScript benchmark

Run the Node.js benchmark script with garbage-collection exposed for more
reliable memory readings:

```sh
make js-bench
# or directly:
node --expose-gc src/bench.js
```

`src/bench.js` uses Benchmark.js and prints per-test metrics including ops/sec,
ns/op, approximate executed iterations and elapsed(ms). At the end it prints
total wall-clock time for all tests.

---

## Bench design notes (important)

- Payloads are prepared once, outside the measured function, to avoid counting
  payload-creation as part of the codec measurement.
- For protobuf-based implementations we prepare instances and/or binary buffers
  before the measured loop so encode/decode are measured in isolation.
- In JS we call `global.gc()` before each test to reduce noise and capture heap
  usage deltas (requires `--expose-gc`).
- The JS bench reports both throughput-derived ns/op and the wall-clock elapsed
  time of each run to give perspective.

---

## Troubleshooting: ESM vs CommonJS (Node)

Common issues when running `node src/bench.js`:

- Error: ERR_MODULE_NOT_FOUND for `./compiled_protos` or local modules
  - In ESM, include file extensions in local imports:
    ```js
    import { test } from './test.js';
    ```
  - If you generated CommonJS (`-w commonjs`) target, use `.cjs` and load it
    from ESM code via `createRequire`:
    ```js
    import { createRequire } from 'module';
    const require = createRequire(import.meta.url);
    const compiled = require('./compiled_protos.cjs');
    const toonbench = compiled.toonbench ?? compiled.default ?? compiled;
    ```
  - If you generated ESM (`-w es6`) import with the explicit extension:
    ```js
    import * as toonbench from './compiled_protos.js';
    ```
- If `PayloadP` or `UserP` are not found after require/import, log
  `Object.keys(...)` of the required module and adapt the access
  (`proto.toonbench?.PayloadP || proto.PayloadP`).

---

## Interpreting results (guidelines)

- JSON.stringify/parse can be very fast in Node because they are implemented in
  V8 (C++).
- protobufjs (pbjs) often has competitive encode/decode for plain objects.
- google-protobuf may have more JS-side overhead for encoding due to its
  Message API.
- Go protobufs usually show low allocations for marshal/unmarshal compared with
  some pure-JS implementations.

Use profiler tools when you need deeper insight:
- Node: `node --prof`, `node --inspect`, Clinic.js (clinic flame)
- Go: standard pprof and bench tooling

---

## Example outputs

Example JS benchmark fragment:

```
JSON.stringify                  ops/sec: 1025.81        ns/op: 974,843  executed (approx): 54    elapsed(ms): 1234
protobufjs encode.finish        ops/sec: 567.26         ns/op: 1,762,849 executed (approx): 30    elapsed(ms): 4321
Total benchmark wall-clock time: 5555 ms (5.555 s)
```

Example Go benchmark fragment:

```
BenchmarkProtoMarshal-12             80          13132078 ns/op         2654402 B/op          2 allocs/op
BenchmarkJsonMarshal-12             50          25097259 ns/op         5963423 B/op          5 allocs/op
```

---

## Adding benchmarks / contributing

- JS: edit `src/bench.js` and add a new entry in the `tests` array with `{
  name, fn }`. Ensure the payload is prepared outside the test function.
- Go: add a `BenchmarkXxx` function in the `_test.go` under `gobench` (or
  appropriate package).
- When adding formats, include instructions to build any extra artifacts and
  add a short note in the README explaining how to reproduce results.

Contributions are welcome — fork, add tests, and open a PR.

---

## Makefile help & banner

- `make` or `make help` prints documented targets that include a `##
  description` on the same line as the target.
- Several Makefile targets print an ASCII banner (Pirate banner) before running
  benches.
- If `make` prints an empty help, ensure the `##` comments exist and the
  Makefile is named `Makefile`.

---

## License

MIT
