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

	message := fmt.Sprintf("User %s connected to chat ...\n", client.name)
	t.broadcast(message, client.name)
	t.mutex.Lock()
	if t.messages != nil {
		for _, s := range t.messages {
			err := h.Write(conn, s)
			if err != nil {
				return
			}
		}
	}
	t.mutex.Unlock()
	for {
		message, err := h.Read(conn)
		if err == io.EOF {
			message = fmt.Sprintf("User %s left from the chat\n", client.name)
			t.broadcast(message, client.name)
			err := conn.Close()
			if err != nil {
				return
			}
			delete(t.onlineClients, client.name)
			return
		}
		message = fmt.Sprintf("%s%s", h.MessageFromUser(client.name), message)
		t.broadcast(message, client.name)
		t.messages = append(t.messages, message)
	}
}

func (t *TcpServer) broadcast(message string, self string) {
	for name, conn := range t.onlineClients {
		if name == self {
			continue
		}
		err := h.Write(conn, message)
		t.errorHandler(err)
	}
}

func (t *TcpServer) errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
