package server

import (
	"bytes"
	"context"
	"errors"
	uuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
	"net"
	//datadir "github.com/superisaac/jointrpc/datadir"
	http "net/http"
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

	router := rpcrouter.RouterFromContext(rootCtx)
	log.Infof("start http proxy at %s", bind)
	mux := http.NewServeMux()
	mux.Handle("/", NewJSONRPCHTTPServer(router))

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

//type handlerFunc func(w http.ResponseWriter, r *http.Request)

type JSONRPCHTTPServer struct {
	router *rpcrouter.Router
}

func NewJSONRPCHTTPServer(router *rpcrouter.Router) *JSONRPCHTTPServer {
	return &JSONRPCHTTPServer{router: router}
}

func (self *JSONRPCHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// only support POST
	if r.Method != "POST" {
		jsonrpc.ErrorResponse(w, r, errors.New("method not allowed"), 405, "Method not allowed")
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
		msg.SetTraceId(uuid.New().String())
	}

	msgvec := rpcrouter.MsgVec{Msg: msg}
	result, err := self.router.CallOrNotify(msgvec, false)
	if err != nil {
		jsonrpc.ErrorResponse(w, r, err, 500, "Server error")
		return
	}
	if result != nil {
		data, err1 := result.Bytes()
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
