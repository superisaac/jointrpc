package server

import (
	//"fmt"
	"bytes"
	"context"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/superisaac/jointrpc/msgutil"
	"github.com/superisaac/jointrpc/rpcrouter"
	"github.com/superisaac/jsonz"
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
		jsonz.ErrorResponse(w, r, errors.New("method not allowed"), 405, "Method not allowed")
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
		jsonz.ErrorResponse(w, r, err, 400, "Bad request")
		return
	}

	msg, err := jsonz.ParseBytes(buffer.Bytes())
	if err != nil {
		jsonz.ErrorResponse(w, r, err, 400, "Bad request")
		return
	}

	// FIXME: sennity check against TraceId
	msg.SetTraceId(r.Header.Get("X-Trace-Id"))
	if msg.TraceId() == "" {
		msg.SetTraceId(jsonz.NewUuid())
	}

	factory := rpcrouter.RouterFactoryFromContext(self.rootCtx)

	router := factory.Get(namespace)
	result, err := router.CallOrNotify(msg, namespace)
	if err != nil {
		jsonz.ErrorResponse(w, r, err, 500, "Server error")
		return
	}
	if result != nil {
		data, err1 := jsonz.MessageBytes(result)
		if err1 != nil {
			jsonz.ErrorResponse(w, r, err1, 500, "Server error")
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

// lives
func (self *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnf("ws upgrade failed %s", err)
		jsonz.ErrorResponse(w, r, err, 400, "Bad request")
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
			//router.Leave(conn)
			router.ChLeave <- rpcrouter.CmdLeave{Conn: conn}
		}
	}()
	adaptor := NewWSAdaptor(ws)
	err = WaitStream(ctx, adaptor, adaptor, conn)
	if err != nil {
		log.Errorf("serve websocket error %+v", err)
	}
}

// WSAdaptor
func NewWSAdaptor(ws *websocket.Conn) *WSAdaptor {
	adaptor := &WSAdaptor{
		ws:   ws,
		done: make(chan error, 10),
	}
	return adaptor
}

func (self WSAdaptor) SendMessage(context context.Context, msg jsonz.Message) error {
	return msgutil.WSSend(self.ws, msg)
}

func (self WSAdaptor) SendCmdMsg(context context.Context, cmdMsg rpcrouter.CmdMsg) error {
	return msgutil.WSSend(self.ws, cmdMsg.Msg)
}

func (self WSAdaptor) Done() chan error {
	return self.done
}

func (self WSAdaptor) Recv() (jsonz.Message, error) {
	msg, err := msgutil.WSRecv(self.ws)
	return msg, err
}
