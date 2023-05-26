package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	addr    string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

var ErrNotConnectionOpen = errors.New("соединение закрыто")

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	diealer := &net.Dialer{
		Timeout: timeout,
	}
	l, err := diealer.Dial("tcp", t.addr)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "...Connected to %s\n", t.addr)
	if err != nil {
		return err
	}

	t.conn = l

	return nil
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return ErrNotConnectionOpen
	}

	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stderr, "...EOF\n")

	return err
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return ErrNotConnectionOpen
	}

	_, err := io.Copy(t.out, t.conn)

	return err
}

func (t *telnetClient) Close() error {
	if t.conn == nil {
		return ErrNotConnectionOpen
	}

	err := t.conn.Close()
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stderr, "...Connection was closed by peer\n")

	return err
}
