package service

import (
	"context"
)

type IService interface {
	CanRun(ctx context.Context) bool
	Start(ctx context.Context) error
}

func TryStartService(ctx context.Context, srv IService) {
	if srv.CanRun(ctx) {
		go srv.Start(ctx)
	}
}
