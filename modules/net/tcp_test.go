package net

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
	"time"
)

var serverAddr = "127.0.0.1:8000"

func TestServer(t *testing.T) {
	waiter := &sync.WaitGroup{}

	server := runServer(t, waiter)

	time.Sleep(10 * time.Millisecond)

	conn1 := runClientAndWrite(t, 1000, "client1")

	time.Sleep(10 * time.Millisecond)

	conn2 := runClientAndWrite(t, 2000, "client2")

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 2, len(server.Incoming()), "Server incoming packet num")

	packet1 := <-server.Incoming()
	assertServerPacket(t, packet1, 1000, "client1", conn1.LocalAddr().String())

	packet2 := <-server.Incoming()
	assertServerPacket(t, packet2, 2000, "client2", conn2.LocalAddr().String())

	assertClientPacket(t, conn1, 1001, "client1")

	assertClientPacket(t, conn2, 2001, "client2")

	assert.Nil(t, conn1.Close(), "Conn1.Close")

	assert.Nil(t, conn2.Close(), "Conn2.Close")

	time.Sleep(10 * time.Millisecond)

	assert.Nil(t, server.Close(), "Server.Close")

	waiter.Wait()
}

func assertServerPacket(t *testing.T, packet *Packet, command uint32, msg string, connAddr string) {
	sessionAddr := packet.Session().Conn().RemoteAddr().String()

	assert.Equal(t, command, packet.Command, "Incoming packet.Command")

	assert.Equal(t, msg, string(packet.Body), "Incoming packet.Body")

	assert.Equal(t, connAddr, sessionAddr, "Incoming packet.session")

	logrus.Infof("Recv: command=%v, body=%s from %v", packet.Command, packet.Body, sessionAddr)

	replyPacket := &Packet{}
	replyPacket.Command = command + 1
	replyPacket.Body = []byte(msg)

	assert.Nil(t, packet.Session().Send(replyPacket), "Session.Send")

	logrus.Infof("Reply: command=%v, body=%s to %v", replyPacket.Command, replyPacket.Body, sessionAddr)
}

func assertClientPacket(t *testing.T, conn net.Conn, command uint32, msg string) {
	packet := &Packet{}

	assert.Nil(t, packet.read(conn), "Client read")

	assert.Equal(t, command, packet.Command, "Packet.Command")

	assert.Equal(t, msg, string(packet.Body), "Packet.Body")

	logrus.Infof("Recv Reply: command=%v, body=%s from %v", packet.Command, packet.Body, conn.RemoteAddr())
}

func runServer(t *testing.T, waiter *sync.WaitGroup) *Server {
	s := NewServer(serverAddr)

	assert.Nil(t, s.Run(waiter), "Server.Run")

	return s
}

func runClientAndWrite(t *testing.T, command uint32, msg string) net.Conn {
	conn, err := net.Dial("tcp", serverAddr)
	assert.Nil(t, err, "Dial Server")

	packet := &Packet{}
	packet.Command = command
	packet.Body = []byte(msg)

	assert.Nil(t, packet.write(conn), "Packet.write")

	logrus.Infof("%v Send: command=%v, body=%s", conn.LocalAddr(), packet.Command, packet.Body)

	return conn
}
