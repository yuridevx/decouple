package pkg

import "context"

type CallOptions struct {
	Context context.Context
	Target  interface{}
}

type CallOption func(opts *CallOptions)

func WithContext(ctx context.Context) CallOption {
	return func(opts *CallOptions) {
		opts.Context = ctx
	}
}

func CopyTo(target interface{}) CallOption {
	return func(opts *CallOptions) {
		opts.Target = target
	}
}

func NewCallOptions(options []CallOption) *CallOptions {
	opts := &CallOptions{}
	for _, o := range options {
		o(opts)
	}
	if opts.Context == nil {
		opts.Context = context.Background()
	}
	return opts
}
