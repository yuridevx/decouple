package tests

import (
	"context"
	"decouple"
	"decouple/local"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEngine(t *testing.T) {
	container := decouple.NewContainer()
	engine := local.NewEngine(container)
	parentCtx := context.Background()

	reply := func(ctx context.Context, req MyIn) MyOut {
		require.Equal(t, ctx, parentCtx)
		return MyOut{Message: req.Message}
	}
	container.Request(reply)

	req := MyIn{Message: "Hi"}

	t.Run("simple request", func(t *testing.T) {
		res, _ := engine.Request(req, decouple.WithContext(parentCtx))
		myRes := res.(MyOut)
		require.Equal(t, req.Message, myRes.Message)
	})

	t.Run("response copy", func(t *testing.T) {
		res := MyOut{}
		_, _ = engine.Request(req, decouple.CopyTo(&res), decouple.WithContext(parentCtx))
		require.Equal(t, res.Message, req.Message)
	})

	t.Run("no request type", func(t *testing.T) {
		require.Panics(t, func() {
			container.Request(func() Res { return Res{} })
		})
	})

	t.Run("no response type", func(t *testing.T) {
		require.Panics(t, func() {
			container.Request(func(_ Req) {})
		})
	})
}

type MyIn struct {
	Message string
}

type MyOut struct {
	Message string
}
