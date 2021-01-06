package misc

import (
	"context"
)

type ContextBinder struct {
	ctx context.Context
}

func NewBinder(ctx context.Context) *ContextBinder {
	return &ContextBinder{ctx: ctx}
}

func (self *ContextBinder) Bind(attr string, value interface{}) *ContextBinder {
	self.ctx = context.WithValue(self.ctx, attr, value)
	return self
}

func (self ContextBinder) Context() context.Context {
	return self.ctx
}
