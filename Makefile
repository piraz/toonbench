.PHONY: bench bench-short

go-bench:
	go test -bench=. -benchmem ./cmd/tut

go-bench-short:
	go test -bench=. -benchmem -benchtime=2x ./cmd/tut

js-bench:
	node --expose-gc src/bench.js	
