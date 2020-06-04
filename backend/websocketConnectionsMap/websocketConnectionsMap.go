package websocketConnectionsMap

import (
	"encoding/json"
	"github.com/fasthttp/websocket"
	"sync"
)

type WebsocketConnections struct {
	mu      sync.RWMutex
	connMap map[int][]*websocket.Conn
}

func (*WebsocketConnections) remove(list []*websocket.Conn, i int) {
	listLen := len(list)
	list[i] = list[listLen-1]
	list = list[:listLen-1]
}

func (c *WebsocketConnections) RemoveConn(userId int, connToRemove *websocket.Conn) {
	c.mu.Lock()
	if len(c.connMap[userId]) <= 1 {
		delete(c.connMap, userId)
	} else {
		for i, conn := range c.connMap[userId] {
			if conn == connToRemove {
				c.remove(c.connMap[userId], i)
				break
			}
		}
	}
	c.mu.Unlock()
}

func (c *WebsocketConnections) AddConn(userId int, connToAdd *websocket.Conn) {
	c.mu.Lock()
	c.connMap[userId] = append(c.connMap[userId], connToAdd)
	c.mu.Unlock()
}

func (c *WebsocketConnections) PushMessageToConnections(messageTo int, message json.RawMessage) {
	c.mu.RLock()
	for _, conn := range c.connMap[messageTo] {
		_ = conn.WriteJSON(message)
	}
	c.mu.RUnlock()
}

func CreateWebsocketConnections() *WebsocketConnections {
	return &WebsocketConnections{
		connMap: make(map[int][]*websocket.Conn),
	}
}
