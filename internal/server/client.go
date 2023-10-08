package server

import "net"

type Client struct {
	name string
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) GetName() string {
	return c.name
}

func (c *Client) SetName(name string) {
	c.name = name
}
