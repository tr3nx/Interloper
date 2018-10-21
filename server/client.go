package main

import (
	"bufio"
	"log"
	"net"
)

type Client struct {
	conn     net.Conn
	incoming chan []byte
	outgoing chan []byte
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func NewClient(conn net.Conn) *Client {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	client := &Client{
		conn:     conn,
		incoming: make(chan []byte),
		outgoing: make(chan []byte),
		reader:   reader,
		writer:   writer,
	}

	client.Listen()

	return client
}

func (c *Client) Read() {
	var errors int
	for {
		if errors > 5 {
			c.Disconnect()
			break
		}

		message, err := unpackMessage(c.reader)
		if err != nil {
			errors++
			continue
		}

		c.incoming <- packMessage(message)
	}
}

func (c *Client) Write() {
	for data := range c.outgoing {
		c.writer.Write(data)
		c.writer.Flush()
	}
}

func (c *Client) Listen() {
	go c.Read()
	go c.Write()
}

func (c *Client) Disconnect() {
	err := c.conn.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Disconnected.")
}
