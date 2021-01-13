package server

import (
	"bytes"
	"errors"
	uuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/rpcrouter"
	http "net/http"
)

func StartHTTPd(http_bind string, router *rpcrouter.Router) {
	log.Infof("start http proxy at %s", http_bind)
	http.HandleFunc("/", jointrpcHandler(router))
	//http.HandleFunc("/", HandleHome)
	log.Fatal(
		http.ListenAndServe(http_bind, nil))
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

func jointrpcHandler(router *rpcrouter.Router) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		result, err := router.CallOrNotify(msgvec, false)
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
			w.Write(data)
		} else {
			data := []byte("{}")
			w.Write(data)
		}
	}
}
