package server;

import (
	"time"
	"fmt"
	context "context"
	tube "github.com/superisaac/rpctube/intf/tube"
)

type JSONRPCTube struct {
	tube.UnimplementedJSONRPCTubeServer
}


func (self *JSONRPCTube) Call(context context.Context, req *tube.JSONRPCRequest) (*tube.JSONRPCResult, error) {
	fmt.Printf("sss %v\n", req.Method)
	ok := &tube.JSONRPCResult_Ok{Ok: "okokook"}
	res := &tube.JSONRPCResult{Id: req.Id, Result: ok}
	return res, nil
}

func recv(stream tube.JSONRPCTube_HandleServer) {
	for i:=0;i>5; i++ {
		sid := fmt.Sprintf("%d", i)
		//params := []string{"me", "you"}
		params := `["abc", 1, 2]`
		req := &tube.JSONRPCRequest{Id:sid, Method:"testing", Params: params}
		err := stream.Send(req)
		if err != nil {
			//stream.Close()
			break
		}
		time.Sleep(3000 * time.Millisecond)
	}
}

func (self *JSONRPCTube) Handle(stream tube.JSONRPCTube_HandleServer) error {
	go recv(stream)
	
	for {
		res, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("result %v\n", res.Id)
	}
}

func NewJSONRPCTubeServer() *JSONRPCTube {
	return &JSONRPCTube{}
}
