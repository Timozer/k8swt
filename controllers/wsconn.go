package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type WsMsg struct {
	MsgType int
	Data    []byte
}

type WsConn struct {
	conn    *websocket.Conn
	ctx     *gin.Context
	inChan  chan *WsMsg
	outChan chan *WsMsg
}

var (
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewWsConn(ctx *gin.Context) (*WsConn, error) {
	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, err
	}

	ret := &WsConn{
		conn:    conn,
		ctx:     ctx,
		inChan:  make(chan *WsMsg, 1000),
		outChan: make(chan *WsMsg, 1000),
	}

	go ret.ReadLoop()
	go ret.WriteLoop()

	return ret, nil
}

// just for debug
func (w *WsConn) getLogger() *zerolog.Logger {
	return getLogger(w.ctx)
}

func (w *WsConn) ReadLoop() {
	defer w.conn.Close()

	logger := getLogger(w.ctx)

	for {
		select {
		case <-w.ctx.Done():
			logger.Debug().Msg("request done")
			return
		default:
			msgType, data, err := w.conn.ReadMessage()
			if err != nil {
				logger.Error().Err(err).Msg("read message fail")
				return
			}
			logger.Debug().Int("WsMsgType", msgType).Str("Data", string(data)).Msg("ReadMsg")
			w.inChan <- &WsMsg{MsgType: msgType, Data: data}
		}
	}
}

func (w *WsConn) WriteLoop() {
	defer w.conn.Close()

	logger := getLogger(w.ctx)

	for {
		select {
		case <-w.ctx.Done():
			logger.Debug().Msg("request done")
			return
		case msg := <-w.outChan:
			logger.Debug().Int("WsMsgType", msg.MsgType).Str("Data", string(msg.Data)).Msg("WriteMsg")
			err := w.conn.WriteMessage(msg.MsgType, msg.Data)
			if err != nil {
				logger.Error().Err(err).Msg("write message fail")
				return
			}
		}
	}
}

func (w *WsConn) Write(msgType int, data []byte) error {
	w.outChan <- &WsMsg{msgType, data}
	return nil
}

func (w *WsConn) Read() (int, []byte, error) {
	msg := <-w.inChan
	return msg.MsgType, msg.Data, nil
}

func (w *WsConn) Close() {
	w.conn.Close()
}
