package net

import (
	"encoding/binary"
	"errors"
	"github.com/sirupsen/logrus"
	"net"
)

type Packet struct {
	session *Session
	Length  uint32
	Command uint32
	Body    []byte
}

func (p *Packet) Session() *Session {
	return p.session
}

func (p *Packet) readHeader(conn net.Conn) error {
	header := make([]byte, 8)
	headerLen, err := conn.Read(header)
	if err != nil {
		logrus.Errorf("Conn Read header with err: %v", err)
		return err
	}

	if headerLen != 8 {
		logrus.Errorf("Conn Read header size != 8, len = %v", headerLen)
		return errors.New("Read header size error")
	}

	p.Length = binary.LittleEndian.Uint32(header[0:4])
	p.Command = binary.LittleEndian.Uint32(header[4:8])

	return nil
}

func (p *Packet) read(conn net.Conn) error {
	if err := p.readHeader(conn); err != nil {
		return err
	}

	if err := p.readBody(conn); err != nil {
		return err
	}

	return nil
}

func (p *Packet) readBody(conn net.Conn) error {
	body := make([]byte, p.Length)
	bodyLen, err := conn.Read(body)
	if err != nil {
		logrus.Errorf("Conn Read body with err: %v", err)
		return err
	}

	if bodyLen != int(p.Length) {
		logrus.Errorf("Conn Read body size != %v, len = %v", p.Length, bodyLen)
		return errors.New("Read body size error")
	}

	p.Body = body

	return nil
}

func (p *Packet) write(conn net.Conn) error {
	if err := p.writeHeader(conn); err != nil {
		return err
	}

	if err := p.writeBody(conn); err != nil {
		return err
	}

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
