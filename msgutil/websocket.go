package msgutil

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/superisaac/jsonrpc"
)

func WSSend(ws *websocket.Conn, msg jsonrpc.IMessage) error {
	msgBytes, err := jsonrpc.MessageBytes(msg)
	if err != nil {
		return err
	}
	return ws.WriteMessage(websocket.TextMessage, msgBytes)
}

func WSRecv(ws *websocket.Conn) (jsonrpc.IMessage, error) {
	for {
		messageType, msgBytes, err := ws.ReadMessage()
		if err != nil {
			return nil, errors.Wrap(err, "ws.ReadMessage()")
		}
		if messageType != websocket.TextMessage {
			log.Infof("message type %d is not text, wait for next", messageType)
			continue
		}
		msg, err := jsonrpc.ParseBytes(msgBytes)
		return msg, err
	}
}
