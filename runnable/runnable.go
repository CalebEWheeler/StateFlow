package runnables

import "context"

type LIFORunnables []Runnable

func (runnables LIFORunnables) Start(ctx context.Context) {
	for _, r := range runnables {
		r.Start(ctx)
	}
}

func (runnables LIFORunnables) Stop(ctx context.Context) {
	for i := len(runnables) - 1; i >= 0; i-- {
		runnables[i].Stop(ctx)
	}
}

type Runnable interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}
