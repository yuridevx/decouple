package tests

import (
	"context"
	"decouple/pkg"
	"decouple/pkg/local"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

var ErrFromSubs = errors.New("hi")
var ErrOtherFromSubs = errors.New("hi other")

func TestSubscribe(t *testing.T) {
	container := pkg.NewContainer()
	engine := local.NewEngine(container)

	SubOne := func(_ context.Context, r MyIn) (MyOut, error) {
		return MyOut{Message: r.Message}, nil
	}

	SubTwo := func(_ MyIn) error {
		return ErrFromSubs
	}

	SubThree := func(_ MyIn) error {
		return ErrOtherFromSubs
	}

	container.Subscribe(SubOne)
	container.Subscribe(SubTwo)
	container.Subscribe(SubThree)

	req := MyIn{
		Message: "Holla",
	}

	t.Run("simple publish", func(t *testing.T) {
		res, err := engine.Broadcast(req)
		require.ErrorIs(t, err, ErrFromSubs)
		require.ErrorIs(t, err, ErrOtherFromSubs)
		require.Len(t, res, 1)
		require.Equal(t, req.Message, res[0].(MyOut).Message)
	})
}
