package client

import (
	"sync"
	"time"

	pb "lark/pkg/proto/pb/ws"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const(
	READTIMEOUT = 300000 * time.Second
)

type Hub interface {
	RemoveClient(string)
	SubmitMsg(*pb.Packet)
}

type Client struct {
	Uuid 		string
	UserId 		string
	Name 		string
	LoginTime 	time.Time
	Conn 		*websocket.Conn
	Server 		Hub
	Writechan 	chan *pb.Packet
	ClosedSignal chan struct{}
	CloseOnce   sync.Once

}

func NewClient(name string, conn *websocket.Conn, server Hub) *Client {
	return &Client{
		Name: name,
		Conn: conn,
		Server: server,
		Writechan: make(chan *pb.Packet, 10),
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
			data,err := proto.Marshal(msg)
			if err != nil {
				return
			}
			err = c.Conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				c.Close()
				return
			}
		}
	}
}

func (c *Client) Read() {
	defer c.Close()
	c.Conn.SetReadDeadline(time.Now().Add(READTIMEOUT))
	for {
		messageType, payload, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}
		c.Conn.SetReadDeadline(time.Now().Add(READTIMEOUT))
		
		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}
		var msg pb.Packet
		if err := proto.Unmarshal(payload, &msg); err != nil {
			return
		}
		
		switch msg.Command {
			case pb.
		}

	}
} 

func (c *Client) Send(msg *pb.Packet) {
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