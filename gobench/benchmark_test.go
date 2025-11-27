package gobench

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"testing"

	"github.com/alpkeskin/gotoon"
	"github.com/piraz/toonbench"
	"github.com/toon-format/toon-go"
	"google.golang.org/protobuf/proto"
)

const numUsers = 5_000

var roles = []string{"admin", "user", "moderator", "superuser"}

func init() {
	fmt.Println("GO TOON serialization microbenchmark")
	fmt.Printf("NUM_USERS: %d\n", numUsers)
	fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
}

func prepareToonPayload(numUsers int) Payload {
	users := make([]User, 0, numUsers)
	for i := range numUsers {
		users = append(users, User{
			Id:   i + 1,
			Name: fmt.Sprintf("User %d", i+1),
			Role: roles[rand.Intn(len(roles))],
		})
	}
	return Payload{Users: users}
}

var samplePaylod = prepareToonPayload(numUsers)

func BenchmarkToonMarshal(b *testing.B) {
	for b.Loop() {
		_, err := toon.Marshal(samplePaylod, toon.WithLengthMarkers(true))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkToonUnmarshal(b *testing.B) {
	data, _ := toon.Marshal(samplePaylod, toon.WithLengthMarkers(true))
	var out Payload
	b.ResetTimer()
	for b.Loop() {
		if err := toon.Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGotoonMarshal(b *testing.B) {
	for b.Loop() {
		_, err := gotoon.Encode(samplePaylod)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonMarshal(b *testing.B) {
	for b.Loop() {
		_, err := json.Marshal(samplePaylod)
		if err != nil {
			b.Fatal(err)
		}
	}
}

var data, _ = json.Marshal(samplePaylod)

func BenchmarkJsonUnmarshal(b *testing.B) {
	var out Payload
	b.ResetTimer()
	for b.Loop() {
		if err := json.Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonUnmarshalMap(b *testing.B) {
	var out map[string]any
	b.ResetTimer()
	for b.Loop() {
		if err := json.Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}

func prepareProtoPayload(n int) *toonbench.PayloadP {

	users := make([]*toonbench.UserP, n)
	for i := range n {
		users[i] = &toonbench.UserP{
			Id:   int32(i + 1),
			Name: fmt.Sprintf("User %d", i+1),
			Role: roles[rand.Intn(len(roles))],
		}
	}
	return &toonbench.PayloadP{Users: users}
}

func BenchmarkProtoMarshal(b *testing.B) {
	var payload = prepareProtoPayload(numUsers)
	b.ResetTimer()
	for b.Loop() {
		_, err := proto.Marshal(payload)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProtoUnmarshal(b *testing.B) {
	payload := prepareProtoPayload(numUsers)
	data, _ := proto.Marshal(payload)
	var out toonbench.PayloadP
	b.ResetTimer()
	for b.Loop() {
		if err := proto.Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}
