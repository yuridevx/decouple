package decouple

import "context"

type CallOptions struct {
	Context context.Context
	Target  interface{}
	Fork    bool
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

func Fork() CallOption {
	return func(opts *CallOptions) {
		opts.Fork = true
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
