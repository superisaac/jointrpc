package server

import (
	"bytes"
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
	//datadir "github.com/superisaac/jointrpc/datadir"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/msgutil"
	"github.com/superisaac/jointrpc/rpcrouter"
)

type HTTPOption struct {
	keyFile  string
	certFile string
}

type HTTPOptionFunc func(opt *HTTPOption)

func WithTLS(certFile string, keyFile string) HTTPOptionFunc {
	return func(opt *HTTPOption) {
		opt.certFile = certFile
		opt.keyFile = keyFile
	}
}

func StartHTTPServer(rootCtx context.Context, bind string, opts ...HTTPOptionFunc) {
	httpOption := &HTTPOption{}
	for _, opt := range opts {
		opt(httpOption)
	}

	log.Infof("start http proxy at %s", bind)
	mux := http.NewServeMux()
	mux.Handle("/metrics", NewMetricsCollector(rootCtx))
	mux.Handle("/ws", NewWSServer(rootCtx))
	mux.Handle("/", NewHTTPServer(rootCtx))

	server := &http.Server{Addr: bind, Handler: mux}
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		panic(err)
	}

	serverCtx, cancelServer := context.WithCancel(rootCtx)
	defer cancelServer()

	go func() {
		for {
			<-serverCtx.Done()
			log.Debugf("http server %s stops", bind)
			listener.Close()
			return
		}
	}()

	if httpOption.certFile != "" && httpOption.keyFile != "" {
		err = server.ServeTLS(listener, httpOption.certFile, httpOption.keyFile)
	} else {
		err = server.Serve(listener)
	}
	if err != nil {
		log.Println("HTTP Server Error - ", err)
		//panic(err)
	}
}

// JSONRPC HTTP Server
type HTTPServer struct {
	//router *rpcrouter.Router
	rootCtx context.Context
}

func NewHTTPServer(rootCtx context.Context) *HTTPServer {
	return &HTTPServer{rootCtx: rootCtx}
}

func (self *HTTPServer) Authorize(r *http.Request) (bool, string) {
	// basic auth
	factory := rpcrouter.RouterFactoryFromContext(self.rootCtx)
	cfg := factory.Config
	if len(cfg.Authorizations) >= 1 {
		if username, password, ok := r.BasicAuth(); ok {
			for _, bauth := range cfg.Authorizations {
				if bauth.Authorize(username, password, r.RemoteAddr) {
					ns := bauth.Namespace
					if ns == "" {
						ns = "default"
					}
					return true, ns
				}
			}
		}
		return false, ""
	}
	return true, "default"
}

func (self *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// only support POST
	if r.Method != "POST" {
		jsonrpc.ErrorResponse(w, r, errors.New("method not allowed"), 405, "Method not allowed")
		return
	}

	ok, namespace := self.Authorize(r)
	if !ok {
		log.Warnf("http auth failed %d", 401)
		w.WriteHeader(401)
		w.Write([]byte("auth failed"))
		return
	}

	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		jsonrpc.ErrorResponse(w, r, err, 400, "Bad request")
		return
	}

	msg, err := jsonrpc.ParseBytes(buffer.Bytes())
	if err != nil {
		jsonrpc.ErrorResponse(w, r, err, 400, "Bad request")
		return
	}

	// FIXME: sennity check against TraceId
	msg.SetTraceId(r.Header.Get("X-Trace-Id"))
	if msg.TraceId() == "" {
		msg.SetTraceId(misc.NewUuid())
	}

	factory := rpcrouter.RouterFactoryFromContext(self.rootCtx)

	result, err := factory.Get(namespace).CallOrNotify(msg, namespace)
	if err != nil {
		jsonrpc.ErrorResponse(w, r, err, 500, "Server error")
		return
	}
	if result != nil {
		data, err1 := jsonrpc.GetMessageBytes(result)
		if err1 != nil {
			jsonrpc.ErrorResponse(w, r, err1, 500, "Server error")
			return
		}
		w.Header().Add("X-Trace-Id", result.TraceId())
		w.Write(data)
	} else {
		data := []byte("{}")
		w.Write(data)
	}
} // ServeHTTP

// Websocket servo
type WSServer struct {
	//router *rpcrouter.Router
	rootCtx context.Context
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10240,
	WriteBufferSize: 10240,
}

func NewWSServer(rootCtx context.Context) *WSServer {
	return &WSServer{rootCtx: rootCtx}
}

func (self *WSServer) Authorize(r *http.Request) (bool, string) {
	// basic auth
	factory := rpcrouter.RouterFactoryFromContext(self.rootCtx)
	cfg := factory.Config
	if len(cfg.Authorizations) >= 1 {
		if username, password, ok := r.BasicAuth(); ok {
			for _, bauth := range cfg.Authorizations {
				if bauth.Authorize(username, password, r.RemoteAddr) {
					ns := bauth.Namespace
					if ns == "" {
						ns = "default"
					}
					return true, ns
				}
			}
		}
		return false, ""
	}
	return true, "default"
}

func relayDownWSMessages(context context.Context, ws *websocket.Conn, conn *rpcrouter.ConnT, chResult chan dispatch.ResultT) {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("recovered ERROR %+v", r)
		}
	}()

	for {
		select {
		case <-context.Done():
			log.Debugf("context done")
			return
		case rest, ok := <-chResult:
			if !ok {
				log.Debugf("conn handler channel closed")
				return
			}
			msgutil.WSSend(ws, rest.ResMsg)
		case msgvec, ok := <-conn.RecvChannel:
			if !ok {
				log.Debugf("recv channel closed")
				return
			}
			msgutil.WSSend(ws, msgvec.Msg)
		case state, ok := <-conn.StateChannel():
			if !ok {
				log.Debugf("state channel closed")
				return
			}
			stateJson := make(map[string]interface{})
			err := mapstructure.Decode(state, &stateJson)
			if err != nil {
				panic(err)
			}
			ntf := jsonrpc.NewNotifyMessage("_state.changed", []interface{}{stateJson})
			msgutil.WSSend(ws, ntf)
		}
	} // and for loop
}

func (self *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnf("ws upgrade failed %s", err)
		jsonrpc.ErrorResponse(w, r, err, 400, "Bad request")
		return
	}

	defer ws.Close()

	conn := rpcrouter.NewConn()

	factory := rpcrouter.RouterFactoryFromContext(self.rootCtx)
	ctx, cancel := context.WithCancel(self.rootCtx)

	defer func() {
		cancel()
		if conn.Joined() {
			router := factory.Get(conn.Namespace)
			router.Leave(conn)
		}
	}()

	chResult := make(chan dispatch.ResultT, misc.DefaultChanSize())
	go relayDownWSMessages(ctx, ws, conn, chResult)

	for {
		msg, err := msgutil.WSRecv(ws)
		if err != nil {
			//jsonrpc.ErrorResponse(w, r, err, 400, "Bad request")
			log.Warnf("bad request %s", err)
			return
		}

		if msg.TraceId() == "" {
			msg.SetTraceId(misc.NewUuid())
		}

		msgvec := rpcrouter.MsgVec{
			Msg:        msg,
			Namespace:  conn.Namespace,
			FromConnId: conn.ConnId,
		}

		streamDisp := GetStreamDispatcher()
		instRes := streamDisp.HandleMessage(ctx, msgvec, chResult, conn)

		if instRes != nil {
			err := msgutil.WSSend(ws, instRes)
			if err != nil {
				break
			}
			if instRes.IsError() {
				break
			}
		}
	} // end of for
} // ServeHTTP
