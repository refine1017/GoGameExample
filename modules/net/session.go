package net

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type Session struct {
	id       uint32
	conn     net.Conn
	incoming chan *Packet
	outgoing chan *Packet
}

func NewSession(conn net.Conn, incoming chan *Packet) *Session {
	s := &Session{}
	s.conn = conn
	s.incoming = incoming
	s.outgoing = make(chan *Packet, 100)

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
		return fmt.Errorf("Session#d overflow", s.id)
	}

	return nil
}

func (s *Session) recv() {
	for {
		if err := s.conn.SetReadDeadline(time.Now().Add(time.Minute)); err != nil {
			logrus.Warning("Conn setReadDeadline with err: %v", err)
			continue
		}

		packet := &Packet{}
		if err := packet.readHeader(s.conn); err != nil {
			return
		}

		if err := packet.readBody(s.conn); err != nil {
			return
		}

		select {
		case s.incoming <- packet:
		default:
			logrus.Warning("Incoming overflow")
		}
	}
}

func (s *Session) send() {
	for {
		select {
		case packet := <-s.outgoing:
			if err := packet.writeHeader(s.conn); err != nil {
				continue
			}

			if err := packet.writeBody(s.conn); err != nil {
				continue
			}
		}
	}
}
