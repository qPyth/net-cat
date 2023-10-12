package server

import (
	"fmt"
	"io"
	"log"
	"net"
	h "net-cat/internal/helpers"
	"runtime"
	"strings"
	"sync"
)

const (
	maxClients = 3
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
			err := conn.Close()
			if err != nil {
				log.Printf("Problem with close connection: %s", err)
			}
		}
		numGoroutines := runtime.NumGoroutine()
		if numGoroutines > maxClients {
			err := h.Write(conn, "Server is full, please try again later\n")
			if err != nil {
				return
			}
			_ = conn.Close()
			continue
		}
		go t.handleClientConnection(conn)
	}
}

func (t *TcpServer) handleClientConnection(conn net.Conn) {
	err := h.Write(conn, h.Greeting())
	t.errorHandler(err)
	var name string
	for {
		err = h.Write(conn, "[Enter your name]: ")
		t.errorHandler(err)

		name, err = h.Read(conn)
		if err != nil {
			if err == io.EOF {
				conn.Close()
				return
			}
		}
		name = strings.ReplaceAll(name, "\n", "")
		if _, ok := t.onlineClients[name]; ok {
			err = h.Write(conn, "this name is already taken, please choose other\n")
			continue
		}
		if !h.NameIsValid(name) {
			err = h.Write(conn, "The name contains illegal characters or is longer than 15 characters\n")
			continue
		}
		break
	}

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
		err := h.Write(client.conn, h.NewInput(client.name))
		if err != nil {
			return
		}
		message, err := h.Read(conn)
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
		message = strings.ReplaceAll(message, "\n", "")
		message = strings.TrimSpace(message)
		if h.MessageIsValid(message) {
			message = fmt.Sprintf("%s%s\n", h.MessageFromUser(client.name), message)
			t.broadcast(message, client.name)
			t.messages = append(t.messages, message)
		} else {
			err := h.Write(conn, h.MessageFromServerIncorrectInput())
			if err != nil {
				return
			}
		}
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
