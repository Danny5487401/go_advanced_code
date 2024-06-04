package main

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkGlobalRand(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Intn(100)
		}
	})
}

func BenchmarkCustomRand(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		rd := rand.New(rand.NewSource(time.Now().Unix()))
		for pb.Next() {
			rd.Intn(100)
		}
	})
}

/*
BenchmarkGlobalRand
BenchmarkGlobalRand-8   	 8989150	       131.1 ns/op
BenchmarkCustomRand
BenchmarkCustomRand-8   	1000000000	         0.8160 ns/op
*/
