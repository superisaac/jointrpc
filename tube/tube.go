package tube

import (
	"context"
)

var (
	tube *TubeT
)

func Tube() *TubeT {
	if tube == nil {
		tube = new(TubeT).Init()
	}
	return tube
}

func (self *TubeT) Init() *TubeT {
	self.Router = NewRouter()
	return self
}

func (self *TubeT) Start(ctx context.Context) {
	self.Router.Start(ctx)
	//go self.Router.Start()
	//self.ServiceManager.Start()
}
