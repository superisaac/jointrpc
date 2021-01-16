package service

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type IService interface {
	Name() string
	CanRun(ctx context.Context) bool
	Start(ctx context.Context) error
}

func TryStartService(ctx context.Context, srv IService) {
	if srv.CanRun(ctx) {
		go func() {
			log.Infof("service %s starts", srv.Name())
			err := srv.Start(ctx)
			if err != nil {
				panic(err)
			}
		}()
	}
}
