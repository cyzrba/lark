package server

import (
	"lark/apps/ws/internal/types"
	"lark/apps/ws/internal/client"
)
type Server struct {
	MsgChan chan types.ReadMsg
	ClientsMap map[string]*client.Client
}

func NewServer() *Server {
	return &Server{
		MsgChan: make(chan types.ReadMsg, 100),
		ClientsMap: make(map[string]*client.Client),
	}
}

func (s *Server) AddClient(c *client.Client) {
	s.ClientsMap[c.Name] = c
}

func (s *Server) RemoveClient(name string) {
	delete(s.ClientsMap, name)
}


func (s *Server) SubmitMsg(msg types.ReadMsg) {
	s.MsgChan <- msg
}

func (s *Server) PushMsg(msg string, name string) {
	c, ok := s.ClientsMap[name]
	if ok {
		c.Send(types.WriteMsg{
			Name: name,
			Msg:  msg,
		})
	}
	
}

func (s *Server) Start() { 
	for msg := range s.MsgChan { 
		s.PushMsg(msg.Msg, msg.Name)
	}
}