package server

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/rpcrouter"
	http "net/http"
)

func StartHTTPd(http_bind string, router *rpcrouter.Router) {
	log.Infof("start http proxy at %s", http_bind)
	http.HandleFunc("/", tubeHandler(router))
	//http.HandleFunc("/", HandleHome)
	log.Fatal(
		http.ListenAndServe(http_bind, nil))
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

func tubeHandler(router *rpcrouter.Router) handlerFunc {
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

		result, err := router.SingleCall(msg, false)
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
