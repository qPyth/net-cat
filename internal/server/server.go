package server

import (
	"log"
	"net"
	h "net-cat/internal/helpers"
	"sync"
)

type TcpServer struct {
	Socket        string
	onlineClients map[string]net.Conn
	mutex         *sync.Mutex
}

func NewTcpServer(socket string) *TcpServer {
	return &TcpServer{Socket: socket, onlineClients: make(map[string]net.Conn), mutex: new(sync.Mutex)}
}

func (t *TcpServer) Start() {
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

	err = h.Write(conn, "Enter your name: \n")
	t.errorHandler(err)

	name, err := h.Read(conn)
	if err != nil {
		log.Fatal(err)
	}
	client := NewClient(conn)
	client.SetName(name)
	t.mutex.Lock()
	t.onlineClients[client.name] = conn
	t.mutex.Unlock()
	t.broadcast()
}

func (t *TcpServer) broadcast() {
	for name, conn := range t.onlineClients {
		message := "User " + name + " connected to chat!"
		err := h.Write(conn, message)
		t.errorHandler(err)
	}
}

func (t *TcpServer) errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
