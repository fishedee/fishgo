package util

import (
	"errors"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	conn      *websocket.Conn
	writeLock sync.Mutex
}

// 创建一个连接
func NewWebsocket(resp http.ResponseWriter, req *http.Request, responseHeader http.Header, readBufSize, writeBufSize int) (*Websocket, error) {
	websck := &Websocket{}

	// 建立websocket连接
	ws, err := websocket.Upgrade(resp, req, responseHeader, readBufSize, writeBufSize)
	if _, ok := err.(websocket.HandshakeError); ok {
		return nil, errors.New("Not a websocket handshake:" + err.Error())
	} else if err != nil {
		return nil, errors.New("Cannot setup WebSocket connection:" + err.Error())
	}

	websck.conn = ws

	return websck, nil
}

// 读消息
func (this *Websocket) ReadMessage() (messageType int, p []byte, err error) {
	if this.conn == nil {
		return 0, []byte{}, errors.New("还未初始化连接！")
	}
	return this.conn.ReadMessage()
}

// 写消息
func (this *Websocket) WriteMessage(messageType int, data []byte) error {
	if this.conn == nil {
		return errors.New("还未初始化连接！")
	}
	this.writeLock.Lock()
	defer this.writeLock.Unlock()

	return this.conn.WriteMessage(messageType, data)
}

// 关闭连接
func (this *Websocket) Close() error {
	if this.conn == nil {
		return errors.New("还未初始化连接！")
	}
	return this.conn.Close()
}

// 获取连接关闭的处理方法
func (this *Websocket) CloseHandler() func(code int, text string) error {
	return this.conn.CloseHandler()
}

// 设置连接关闭的处理方法
func (this *Websocket) SetCloseHandler(h func(code int, text string) error) {
	if this.conn == nil {
		return
	}
	this.conn.SetCloseHandler(h)
}

func (this *Websocket) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}
