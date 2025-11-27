//go:build parallel

package gobench

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/alpkeskin/gotoon"
	"github.com/piraz/toonbench"
	"github.com/toon-format/toon-go"
	"google.golang.org/protobuf/proto"
)

func init() {
	fmt.Println("GO TOON serialization microbenchmark")
	fmt.Println("Parallel Executions")
	fmt.Printf("NUM_USERS: %d\n", numUsers)
	fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
}

var samplePaylodParallel = prepareToonPayload(numUsers)

func BenchmarkToonMarshalParallel(b *testing.B) {
	// Prepara N payloads
	payloads := make([]Payload, b.N)
	for i := range payloads {
		payloads[i] = samplePaylodParallel
	}

	b.ResetTimer()
	var wg sync.WaitGroup
	for _, p := range payloads {
		wg.Add(1)
		go func(pl Payload) {
			defer wg.Done()
			_, _ = toon.Marshal(pl, toon.WithLengthMarkers(true))
		}(p)
	}
	wg.Wait()
}

func BenchmarkToonUnmarshalParallel(b *testing.B) {
	// Prepara um buffer serializado para todos
	data, _ := toon.Marshal(samplePaylodParallel, toon.WithLengthMarkers(true))
	buffers := make([][]byte, b.N)
	for i := range buffers {
		buffers[i] = make([]byte, len(data))
		copy(buffers[i], data)
	}

	b.ResetTimer()
	var wg sync.WaitGroup
	for _, buf := range buffers {
		wg.Add(1)
		go func(bf []byte) {
			defer wg.Done()
			var out Payload
			_ = toon.Unmarshal(bf, &out)
		}(buf)
	}
	wg.Wait()
}

func BenchmarkGotoonMarshalParallel(b *testing.B) {
	// Prepara N payloads
	payloads := make([]Payload, b.N)
	for i := range payloads {
		payloads[i] = samplePaylodParallel
	}

	b.ResetTimer()
	var wg sync.WaitGroup
	for _, p := range payloads {
		wg.Add(1)
		go func(pl Payload) {
			defer wg.Done()
			_, _ = gotoon.Encode(pl)
		}(p)
	}
	wg.Wait()
}

func BenchmarkJsonMarshalParallel(b *testing.B) {
	// Prepare a slice of payloads to be marshaled in parallel
	payloads := make([]Payload, b.N)
	for i := range payloads {
		payloads[i] = prepareToonPayload(numUsers)
	}

	b.ResetTimer()
	var wg sync.WaitGroup
	for _, p := range payloads {
		wg.Add(1)
		go func(pl Payload) {
			defer wg.Done()
			_, _ = json.Marshal(pl)
		}(p)
	}
	wg.Wait()
}

func BenchmarkJsonUnmarshalParallel(b *testing.B) {
	buffers := make([][]byte, b.N)
	for i := range buffers {
		buffers[i], _ = json.Marshal(prepareToonPayload(numUsers))
	}

	b.ResetTimer()
	var wg sync.WaitGroup
	for _, buf := range buffers {
		wg.Add(1)
		go func(b []byte) {
			defer wg.Done()
			var out Payload
			json.Unmarshal(b, &out)
		}(buf)
	}
	wg.Wait()
}

func BenchmarkProtoMarshalParallel(b *testing.B) {
	payloads := make([]*toonbench.PayloadP, b.N)
	for i := range payloads {
		payloads[i] = prepareProtoPayload(numUsers)
	}

	b.ResetTimer()
	var wg sync.WaitGroup
	for _, p := range payloads {
		wg.Add(1)
		go func(pl *toonbench.PayloadP) {
			defer wg.Done()
			_, _ = proto.Marshal(pl)
		}(p)
	}
	wg.Wait()
}

func BenchmarkProtoUnmarshalParallel(b *testing.B) {
	payload := prepareProtoPayload(numUsers)
	buf, _ := proto.Marshal(payload)
	buffers := make([][]byte, b.N)
	for i := range buffers {
		buffers[i] = make([]byte, len(buf))
		copy(buffers[i], buf)
	}

	b.ResetTimer()
	var wg sync.WaitGroup
	for _, data := range buffers {
		wg.Add(1)
		go func(d []byte) {
			defer wg.Done()
			var out toonbench.PayloadP
			_ = proto.Unmarshal(d, &out)
		}(data)
	}
	wg.Wait()
}
