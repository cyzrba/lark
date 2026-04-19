package client

import (
	"sync"

	"lark/apps/ws/internal/types"

	"github.com/gorilla/websocket"
)


type Hub interface {
	RemoveClient(string)
	SubmitMsg(types.ReadMsg)
}

type Client struct {
	Name 		string
	Conn 		*websocket.Conn
	Server 		Hub
	Writechan 	chan types.WriteMsg
	ClosedSignal chan struct{}
	CloseOnce   sync.Once

}

func NewClient(name string, conn *websocket.Conn, server Hub) *Client {
	return &Client{
		Name: name,
		Conn: conn,
		Server: server,
		Writechan: make(chan types.WriteMsg, 10),
		ClosedSignal: make(chan struct{}),
	}
}

func (c *Client) Close() { 
	c.CloseOnce.Do(func() {
		close(c.ClosedSignal)
		close(c.Writechan)
		c.Conn.Close()
		c.Server.RemoveClient(c.Name)
	})
}

func (c *Client) Write() { 
	defer c.Close()
	for {
		select {
		case <-c.ClosedSignal:
			return
		case msg := <-c.Writechan:
		err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Name + ": " + msg.Msg))
		if err != nil {
			c.Close()
			return
		}
		}
	}
}

func (c *Client) Read() {
	defer c.Close()
	for {
		messageType, payload, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}
		c.Server.SubmitMsg(types.ReadMsg{
			Name: c.Name,
			Msg:  string(payload),
		})

	}
} 

func (c *Client) Send(msg types.WriteMsg) {
	select {
	case <- c.ClosedSignal:
		return
	case c.Writechan <- msg:
	}
}


func (c *Client) Start() { 
	go c.Write()
	go c.Read()
}