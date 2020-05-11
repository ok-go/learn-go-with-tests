package ctx_aware_reader

import (
	"context"
	"io"
)

type readerCtx struct {
	ctx      context.Context
	delegate io.Reader
}

func NewCancellableReader(ctx context.Context, r io.Reader) io.Reader {
	return &readerCtx{
		ctx:      ctx,
		delegate: r,
	}
}

func (r *readerCtx) Read(p []byte) (n int, err error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.delegate.Read(p)
}
