package net

import (
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

type Server struct {
	listener net.Listener
	waiter   *sync.WaitGroup
	addr     string
	sessions map[net.Addr]*Session
	incoming chan *Packet
	outgoing chan *Packet
}

func NewServer(addr string) *Server {
	s := &Server{}
	s.addr = addr
	s.sessions = make(map[net.Addr]*Session)
	s.incoming = make(chan *Packet)
	s.outgoing = make(chan *Packet)

	return s
}

func (s *Server) Run(waiter *sync.WaitGroup) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.listener = listener

	s.waiter = waiter
	waiter.Add(1)

	go s.Accept()

	return nil
}

func (s *Server) Accept() {
	defer s.waiter.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				logrus.Error("Server.Accept:", ne.Error(), "[temporary]")
				time.Sleep(time.Millisecond)
				continue
			}
			return
		}

		session := NewSession(conn, s.incoming)
		s.sessions[conn.RemoteAddr()] = session
		session.Go()

		logrus.Infof("Accept new session: %v", conn.RemoteAddr())
	}
}
