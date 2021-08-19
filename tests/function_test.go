package tests

import (
	"context"
	"decouple/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type FnCase struct {
	fn   interface{}
	name string
	fail bool
	ctx  bool
	arg  bool
	err  bool
	ret  bool
}

type Res struct {
}
type Req struct {
}

func TestParseFunction(t *testing.T) {
	table := []FnCase{
		{
			fn:   func() {},
			name: "empty",
		},
		{
			fn:   func(req Req) {},
			name: "req",
			arg:  true,
		},
		{
			fn:   func(ctx context.Context, req Req) {},
			name: "ctx req",
			arg:  true,
			ctx:  true,
		},
		{
			fn:   func(req Req, ctx context.Context) {},
			name: "ctx req",
			fail: true,
		},
		{
			fn:   func() Res { return Res{} },
			name: "res",
			ret:  true,
		},
		{
			fn:   func() (Res, error) { return Res{}, nil },
			name: "res err",
			ret:  true,
			err:  true,
		},
		{
			fn:   func() error { return nil },
			name: "err",
			err:  true,
		},
		{
			fn:   func() (Res, Res) { return Res{}, Res{} },
			name: "res res",
			fail: true,
		},
		{
			fn:   func() (Res, error, Res) { return Res{}, nil, Res{} },
			name: "res err res",
			fail: true,
		},
		{
			fn:   func(req *Req) {},
			name: "reqRef",
			fail: true,
		},
		{
			fn:   func(ctx context.Context, req *Req) {},
			name: "ctx reqRef",
			fail: true,
		},
		{
			fn:   func() *Res { return nil },
			name: "ctx reqRef",
			fail: true,
		},
	}

	for _, i := range table {
		t.Run(i.name, func(t *testing.T) {
			fn, err := types.ParseFunction(i.fn)
			if i.fail {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, fn.ContextPresent, i.ctx)
			require.Equal(t, fn.IsArgumentPresent(), i.arg)
			require.Equal(t, fn.ErrorPresent, i.err)
			require.Equal(t, fn.IsReturnPresent(), i.ret)
		})
	}
}
