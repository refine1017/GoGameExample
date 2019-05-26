package net

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

var serverAddr = "127.0.0.1:8000"

func TestServer(t *testing.T) {
	s, err := runServer()
	if err != nil {
		t.Errorf("runServer with err: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	conn1, err := runClientAndWrite(1000, "client1")
	if err != nil {
		t.Errorf("runClient with err: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	conn2, err := runClientAndWrite(2000, "client2")
	if err != nil {
		t.Errorf("runClient with err: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	if len(s.Incoming()) != 2 {
		t.Errorf("Server incoming length want 2, got %v", len(s.Incoming()))
	}

	packet1 := <-s.Incoming()
	if err := assertPacket(packet1, 1000, "client1", conn1.LocalAddr().String()); err != nil {
		t.Error(err)
	}

	packet2 := <-s.Incoming()
	if err := assertPacket(packet2, 2000, "client2", conn2.LocalAddr().String()); err != nil {
		t.Error(err)
	}
}

func assertPacket(packet *Packet, command uint32, msg string, connAddr string) error {
	sessionAddr := packet.Session().Conn().RemoteAddr().String()

	if packet.Command != command {
		return fmt.Errorf("Incoming packet.Command: want 1000, got %v", packet.Command)
	}

	body := string(packet.Body)
	if !strings.EqualFold(body, msg) {
		return fmt.Errorf("Incoming packet.Body: want hello, got %s", packet.Body)
	}

	if sessionAddr != connAddr {
		return fmt.Errorf("Incoming packet.session: want %v, got %s", connAddr, sessionAddr)
	}

	logrus.Infof("Recv: command=%v, body=%s from %v", packet.Command, packet.Body, sessionAddr)

	return nil
}

func runServer() (*Server, error) {
	s := NewServer(serverAddr)

	if err := s.Run(&sync.WaitGroup{}); err != nil {
		return nil, err
	}

	return s, nil
}

func runClientAndWrite(command uint32, msg string) (net.Conn, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}

	packet := &Packet{}
	packet.Command = command
	packet.Body = []byte(msg)

	if err := packet.write(conn); err != nil {
		return nil, err
	}

	logrus.Infof("%v Send: command=%v, body=%s", conn.LocalAddr(), packet.Command, packet.Body)

	return conn, nil
}
