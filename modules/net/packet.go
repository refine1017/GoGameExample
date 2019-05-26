package net

import (
	"encoding/binary"
	"errors"
	"github.com/sirupsen/logrus"
	"net"
)

type Packet struct {
	Length  uint32
	Command uint32
	Body    []byte
}

func (p *Packet) readHeader(conn net.Conn) error {
	header := make([]byte, 8)
	headerLen, err := conn.Read(header)
	if err != nil {
		logrus.Error("Conn Read header with err: %v", err)
		return err
	}

	if headerLen != 8 {
		logrus.Error("Conn Read header size != 8, len = %v", headerLen)
		return errors.New("Read header size error")
	}

	p.Length = binary.LittleEndian.Uint32(header[0:4])
	p.Command = binary.LittleEndian.Uint32(header[4:8])

	return nil
}

func (p *Packet) readBody(conn net.Conn) error {
	body := make([]byte, p.Length)
	bodyLen, err := conn.Read(body)
	if err != nil {
		logrus.Error("Conn Read body with err: %v", err)
		return err
	}

	if bodyLen != 8 {
		logrus.Error("Conn Read body size != %v, len = %v", p.Length, bodyLen)
		return errors.New("Read body size error")
	}

	p.Body = body

	return nil
}

func (p *Packet) writeHeader(conn net.Conn) error {
	p.Length = uint32(len(p.Body))

	bytes := make([]byte, 8)

	binary.LittleEndian.PutUint32(bytes[0:4], p.Length)
	binary.LittleEndian.PutUint32(bytes[4:8], p.Command)

	n, err := conn.Write(bytes)
	if err != nil {
		return err
	}

	if n != 8 {
		return errors.New("Write header size error")
	}

	return nil
}

func (p *Packet) writeBody(conn net.Conn) error {
	n, err := conn.Write(p.Body)
	if err != nil {
		return err
	}

	if n != int(p.Length) {
		return errors.New("Write body size error")
	}

	return nil
}
