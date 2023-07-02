package controllers

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"k8s.io/client-go/tools/remotecommand"
)

type Event struct {
	Type string
	Data interface{}
}

const (
	EVT_RESIZE = "resize"
)

type Size struct {
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

type StreamHandler struct {
	wsConn     *WsConn
	resizeChan chan *remotecommand.TerminalSize
	logger     *zerolog.Logger
}

func NewStreamHandler(wsConn *WsConn) *StreamHandler {
	return &StreamHandler{
		wsConn:     wsConn,
		resizeChan: make(chan *remotecommand.TerminalSize, 10),
		logger:     wsConn.getLogger(),
	}
}

func (s *StreamHandler) Next() *remotecommand.TerminalSize {
	size := <-s.resizeChan
	s.logger.Debug().Any("size", size).Msg("stream handler next() called")
	return size
}

func (s *StreamHandler) Read(p []byte) (int, error) {
	s.logger.Debug().Msg("stream handler read() called")
	_, data, err := s.wsConn.Read()
	if err != nil {
		return 0, err
	}
	evt := Event{}
	err = json.Unmarshal(data, &evt)
	if err != nil {
		copy(p, data)
		return len(data), nil
	}
	switch evt.Type {
	case EVT_RESIZE:
		size := Size{}
		tmpData, err := json.Marshal(evt.Data)
		if err != nil {
			return 0, err
		}
		err = json.Unmarshal(tmpData, &size)
		if err != nil {
			return 0, err
		}
		s.resizeChan <- &remotecommand.TerminalSize{Width: size.Cols, Height: size.Rows}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported event type %s", evt.Type)
	}
}

func (s *StreamHandler) Write(p []byte) (int, error) {
	s.logger.Debug().Msg("stream handler write() called")
	var err error
	data := make([]byte, len(p))
	copy(data, p)
	if utf8.ValidString(string(data)) {
		err = s.wsConn.Write(websocket.TextMessage, data)
	} else {
		err = s.wsConn.Write(websocket.BinaryMessage, data)
	}
	return len(data), err
}
