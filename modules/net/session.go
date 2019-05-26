package net

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"time"
)

type Session struct {
	conn     net.Conn
	incoming chan *Packet
	outgoing chan *Packet
	closing  chan *Session
}

func NewSession(conn net.Conn, incoming chan *Packet, closing chan *Session) *Session {
	s := &Session{}
	s.conn = conn
	s.incoming = incoming
	s.outgoing = make(chan *Packet, 100)
	s.closing = closing

	return s
}

func (s *Session) Go() {
	go s.recv()
	go s.send()
}

func (s *Session) Send(p *Packet) error {
	select {
	case s.outgoing <- p:
	default:
		return fmt.Errorf("Session#%d overflow", s.conn.RemoteAddr())
	}

	return nil
}

func (s *Session) Conn() net.Conn {
	return s.conn
}

func (s *Session) Close() {
	if s.closing != nil {
		s.closing <- s
		s.closing = nil
		_ = s.Conn().Close()
	}
}

func (s *Session) recv() {
	defer s.Close()

	for {
		if err := s.conn.SetReadDeadline(time.Now().Add(time.Minute)); err != nil {
			logrus.Warningf("Conn setReadDeadline with err: %v", err)
			continue
		}

		packet := &Packet{}
		if err := packet.read(s.conn); err != nil {
			if err == io.EOF {
				return
			}

			logrus.Warningf("Conn read with err: %v", err)
			return
		}

		packet.session = s

		select {
		case s.incoming <- packet:
		default:
			logrus.Warning("Incoming overflow")
		}
	}
}

func (s *Session) send() {
	defer s.Close()

	for {
		select {
		case packet := <-s.outgoing:
			if err := packet.write(s.conn); err != nil {
				logrus.Warningf("Conn write with err: %v", err)
				return
			}
		}
	}
}
