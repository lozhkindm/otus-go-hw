package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	Close() error // this linter is hilarious - Error: client.Close undefined (type TelnetClient has no field or method Close) (typecheck)
	Send() error
	Receive() error
}

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *Client) Connect() (err error) {
	if c.conn, err = net.DialTimeout("tcp", c.address, c.timeout); err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stderr, "...connected to %s\n", c.address)
	return err
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Send() (err error) {
	if _, err = io.Copy(c.conn, c.in); err != nil {
		return err
	}
	_, err = fmt.Fprintln(os.Stderr, "...EOF")
	return err
}

func (c *Client) Receive() (err error) {
	if _, err = io.Copy(c.out, c.conn); err != nil {
		return err
	}
	_, err = fmt.Fprintln(os.Stderr, "...connection was closed by peer")
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
