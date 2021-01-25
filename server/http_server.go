package server

import (
	"bytes"
	"context"
	"errors"
	uuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	//datadir "github.com/superisaac/jointrpc/datadir"
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

	log.Infof("start http proxy at %s", bind)
	mux := http.NewServeMux()
	mux.Handle("/metrics", NewMetricsCollector(rootCtx))
	mux.Handle("/", NewJSONRPCHTTPServer(rootCtx))

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
	//router *rpcrouter.Router
	rootCtx context.Context
}

func NewJSONRPCHTTPServer(rootCtx context.Context) *JSONRPCHTTPServer {
	return &JSONRPCHTTPServer{rootCtx: rootCtx}
}

func (self *JSONRPCHTTPServer) Authorize(r *http.Request) bool {
	// basic auth
	router := rpcrouter.RouterFromContext(self.rootCtx)
	cfg := router.Config
	if len(cfg.Authorizations) >= 1 {
		if username, password, ok := r.BasicAuth(); ok {
			for _, bauth := range cfg.Authorizations {
				if bauth.Username == username && bauth.Password == password {
					return true
				}
			}
		}
		return false
	}
	return true
}

func (self *JSONRPCHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// only support POST
	if r.Method != "POST" {
		jsonrpc.ErrorResponse(w, r, errors.New("method not allowed"), 405, "Method not allowed")
		return
	}

	if !self.Authorize(r) {
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
		msg.SetTraceId(uuid.New().String())
	}

	router := rpcrouter.RouterFromContext(self.rootCtx)
	result, err := router.CallOrNotify(msg)
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
