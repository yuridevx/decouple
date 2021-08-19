package tests

import (
	"context"
	"decouple"
	"decouple/local"
	"testing"
)

func BenchmarkEngine(t *testing.B) {
	container := decouple.NewContainer()
	engine := local.NewEngine(container)
	container.Request(respondBench)

	myReq := MyBench{"hi"}
	myRes := MyBench{}

	t.Run("benchmark request", func(b *testing.B) {
		resRaw, _ := engine.Request(myReq)
		_ = resRaw.(MyBench)
	})

	t.Run("benchmark request copy", func(b *testing.B) {
		_, _ = engine.Request(myReq, decouple.CopyTo(&myRes))
	})
}

type MyBench struct {
	Message string
}

func respondBench(_ context.Context, req MyBench) MyBench {
	return req
}
