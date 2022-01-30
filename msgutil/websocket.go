package msgutil

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/superisaac/jsonz"
)

func WSSend(ws *websocket.Conn, msg jsonz.Message) error {
	msgBytes, err := jsonz.MessageBytes(msg)
	if err != nil {
		return err
	}
	return ws.WriteMessage(websocket.TextMessage, msgBytes)
}

func WSRecv(ws *websocket.Conn) (jsonz.Message, error) {
	for {
		messageType, msgBytes, err := ws.ReadMessage()
		if err != nil {
			return nil, errors.Wrap(err, "ws.ReadMessage()")
		}
		if messageType != websocket.TextMessage {
			log.Infof("message type %d is not text, wait for next", messageType)
			continue
		}
		msg, err := jsonz.ParseBytes(msgBytes)
		return msg, err
	}
}
