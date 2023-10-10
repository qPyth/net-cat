package server

import (
	"fmt"
	"io"
	"log"
	"net"
	h "net-cat/internal/helpers"
	"strings"
	"sync"
)

const (
	maxClients = 10
)

type TcpServer struct {
	Socket        string
	onlineClients map[string]net.Conn
	mutex         *sync.Mutex
	messages      []string
}

func NewTcpServer(socket string) *TcpServer {
	return &TcpServer{Socket: socket, onlineClients: make(map[string]net.Conn), mutex: new(sync.Mutex), messages: []string{}}
}

func (t *TcpServer) Start() {
	log.Printf("Server starting on %s\n", t.Socket)
	listener, err := net.Listen("tcp", t.Socket)
	if err != nil {
		log.Fatal(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go t.handleClientConnection(conn)
	}
}

func (t *TcpServer) handleClientConnection(conn net.Conn) {
	err := h.Write(conn, h.Greeting())
	t.errorHandler(err)

	err = h.Write(conn, "[Enter your name: ]\n")
	t.errorHandler(err)

	name, err := h.Read(conn)
	if err != nil {
		log.Fatal(err)
	}
	name = strings.ReplaceAll(name, "\n", "")
	client := NewClient(conn)
	client.SetName(name)

	t.mutex.Lock()
	t.onlineClients[client.name] = conn
	t.mutex.Unlock()

	t.broadcast(h.MessageFromServerNewConnect(client.name), client.name)

	if t.messages != nil {
		t.mutex.Lock()
		for _, s := range t.messages {
			s = strings.Replace(s, "\n", "", 1)
			err := h.Write(conn, s)
			if err != nil {
				return
			}
		}
		t.mutex.Unlock()
	}

	for {
		h.Write(client.conn, h.NewInput(client.name))
		message, err := h.Read(conn)
		message = strings.ReplaceAll(message, "\n", "")
		if err == io.EOF {
			delete(t.onlineClients, client.name)
			t.broadcast(h.MessageFromServerLeftFromChat(client.name), client.name)
			h.Write(client.conn, h.MessageFromUser(client.name))
			err := conn.Close()
			if err != nil {
				return
			}
			return
		}
		message = fmt.Sprintf("%s%s\n", h.MessageFromUser(client.name), message)
		t.broadcast(message, client.name)
		t.messages = append(t.messages, message)
	}
}

func (t *TcpServer) broadcast(message string, writingUser string) {
	t.mutex.Lock()
	for name, conn := range t.onlineClients {
		if name == writingUser {
			continue
		}
		err := h.Write(conn, message)
		if err != nil {
			return
		}
		t.errorHandler(err)
		h.Write(conn, h.NewInput(name))
	}
	t.mutex.Unlock()
}

func (t *TcpServer) errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
