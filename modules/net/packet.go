package net

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

const HeaderSize = 8

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
	header := make([]byte, HeaderSize)
	headerLen, err := conn.Read(header)
	if err != nil {
		return err
	}

	if headerLen != HeaderSize {
		return fmt.Errorf("header size != %v, len = %v", HeaderSize, headerLen)
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
		return err
	}

	if bodyLen != int(p.Length) {
		return fmt.Errorf("body size != %v, len = %v", p.Length, bodyLen)
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
		return errors.New("write header size error")
	}

	return nil
}

func (p *Packet) writeBody(conn net.Conn) error {
	n, err := conn.Write(p.Body)
	if err != nil {
		return err
	}

	if n != int(p.Length) {
		return errors.New("write body size error")
	}

	return nil
}
