package server

import (
	"lark/apps/ws/internal/client"
	"lark/apps/chat/rpc/chatrpc"
	pb "lark/pkg/proto/pb/ws"

	"time"
	"sync"
)

const sendMessageTimeout = 3 * time.Second

type Server struct {
	MsgChan    chan *pb.Packet
	ClientsMap sync.Map // map[string]*client.Client
	ChatRpc    chatrpc.ChatRpc
}

func NewServer(chatRpc chatrpc.ChatRpc) *Server {
	return &Server{
		MsgChan:    make(chan *pb.Packet, 100),
		ClientsMap: sync.Map{},
		ChatRpc:    chatRpc,
	}
}

func (s *Server) AddClient(c *client.Client) {
	s.ClientsMap.Store(c.Name, c)
}

func (s *Server) RemoveClient(name string) {
	s.ClientsMap.Delete(name)
}


func (s *Server) SubmitMsg(msg *pb.Packet) {
	s.MsgChan <- msg
}

func (s *Server) PushMsg(msg *pb.Packet) {
	// c, ok := s.ClientsMap[msg.GetChat().To]
	// if ok {
	// 	c.Send(msg)
	// }
}

func (s *Server) Start() {
	for msg := range s.MsgChan {
		switch msg.Command {
		case pb.Command_COMMAND_HEARTBEAT_PING:
			
			
		}
	}
}
