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
	sessions *sync.Map
	incoming chan *Packet
	outgoing chan *Packet
	closing  chan *Session
}

func NewServer(addr string) *Server {
	s := &Server{}
	s.addr = addr
	s.sessions = new(sync.Map)
	s.incoming = make(chan *Packet, 1000)
	s.outgoing = make(chan *Packet, 1000)
	s.closing = make(chan *Session, 100)

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

	go s.accept()
	go s.monitor()

	return nil
}

func (s *Server) Incoming() chan *Packet {
	return s.incoming
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) accept() {
	defer s.waiter.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				logrus.Error("Server.accept:", ne.Error(), "[temporary]")
				time.Sleep(time.Millisecond)
				continue
			}
			return
		}

		session := NewSession(conn, s.incoming, s.closing)
		s.sessions.Store(conn.RemoteAddr(), session)
		session.Go()

		logrus.Infof("Accept new session: %v", conn.RemoteAddr())
	}
}

func (s *Server) monitor() {
	for {
		session := <-s.closing

		addr := session.Conn().RemoteAddr()
		s.sessions.Delete(addr)

		logrus.Infof("Session %v closed", addr)
	}
}
