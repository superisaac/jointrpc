package server

// exposition format of prometheus

import (
	//"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
	"strings"
	//"net"
	//datadir "github.com/superisaac/jointrpc/datadir"
	http "net/http"
)

type MetricsCollector struct {
	//router *rpcrouter.Router
	rootCtx context.Context
}

func NewMetricsCollector(rootCtx context.Context) *MetricsCollector {
	return &MetricsCollector{rootCtx: rootCtx}
}

func (self *MetricsCollector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// only support POST
	if r.Method != "GET" && r.Method != "HEAD" {
		w.WriteHeader(405)
		w.Write([]byte("405 Method not allowed"))
		return
	}

	router := rpcrouter.RouterFromContext(self.rootCtx)

	if router.Config.Metrics.BearerToken != "" && r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", router.Config.Metrics.BearerToken) {
		w.WriteHeader(401)
		w.Write([]byte("Authorization failed"))
		return
	}

	lines, err := self.Collect()
	if err != nil {
		log.Warnf("HTTP Error collect metrics, %s", err.Error())
		w.WriteHeader(500)
		w.Write([]byte("500 server error"))
		return
	}

	data := []byte(strings.Join(lines, "\n"))
	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
	w.Write([]byte("\n"))
}

func (self *MetricsCollector) Collect() ([]string, error) {
	msgId := misc.NewUuid()
	emptyArr := make([]interface{}, 0)
	reqmsg := jsonrpc.NewRequestMessage(msgId, "metrics.collect", emptyArr, nil)
	router := rpcrouter.RouterFromContext(self.rootCtx)
	resmsg, err := router.CallOrNotify(reqmsg, rpcrouter.WithBroadcast(true))
	if err != nil {
		return nil, nil
	}
	if resmsg.IsError() {
		return nil, errors.New(fmt.Sprintf("call metrifs.collect error %#v", resmsg.MustError()))
	}
	res := resmsg.MustResult()
	resArr, ok := res.([]interface{})
	fmt.Printf("sssssss %#v\n", resArr)
	if !ok {
		return nil, nil
	}
	var lines []string
	for _, a := range resArr {
		b, ok := a.(map[string]interface{})
		if !ok {
			continue
		}
		fmt.Printf("hhhhh %#v\n", b)
		msgProto := simplejson.New()
		msgProto.SetPath(nil, b)
		msgItem, err := jsonrpc.Parse(msgProto)
		if err != nil {
			log.Warnf("metrics error %s %+v", err.Error(), b)
			continue
		}

		if msgItem.IsResult() {
			result := msgItem.MustResult()
			childLines := self.buildMetricsLines(result)
			fmt.Printf("metrics child lines %#v, %#v\n", result, childLines)
			lines = append(lines, childLines...)
		} else {
			// TODO: log error
		}
	}
	return lines, nil
} // Collect()

func (self MetricsCollector) buildMetricsLines(result interface{}) []string {
	if resStr, ok := result.(string); ok {
		return strings.Split(resStr, "\n")
	}
	if resArr, ok := result.([]interface{}); ok {
		var arr []string
		for _, b := range resArr {
			if sb, ok := b.(string); ok {
				arr = append(arr, sb)
			} else {
				// not string
				return []string{}
			}
		}
		return arr
	}
	return []string{}
}