package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/alpkeskin/gotoon"
	"github.com/toon-format/toon-go"
	"google.golang.org/protobuf/proto"
)

var roles = []string{"admin", "user", "moderator", "superuser"}

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

var samplePaylod = prepareToonPayload(100_000)

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

func prepareProtoPayload(n int) *PayloadP {
	users := make([]*UserP, n)
	for i := range n {
		users[i] = &UserP{
			Id:   int32(i + 1),
			Name: fmt.Sprintf("User %d", i+1),
			Role: roles[rand.Intn(len(roles))],
		}
	}
	return &PayloadP{Users: users}
}

func BenchmarkProtoMarshal(b *testing.B) {
	var payload = prepareProtoPayload(100_000)
	b.ResetTimer()
	for b.Loop() {
		_, err := proto.Marshal(payload)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProtoUnmarshal(b *testing.B) {
	payload := prepareProtoPayload(100_000)
	data, _ := proto.Marshal(payload)
	var out PayloadP
	b.ResetTimer()
	for b.Loop() {
		if err := proto.Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}
