package guiniol

import "golang.org/x/net/context"

var defaultCb = CallbackFunc(func(_ context.Context, _ *EventCtx) {})

type Callback interface {
	Next(context.Context, *EventCtx)
}

type CallbackFunc func(context.Context, *EventCtx)

func (cf CallbackFunc) Next(c context.Context, ctx *EventCtx) {
	cf(c, ctx)
}
