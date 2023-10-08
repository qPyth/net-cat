package app

import (
	"flag"
	"net-cat/internal/server"
)

type App struct {
	Server *server.TcpServer
}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() {
	port := flag.String("port", "3000", "[USAGE]: ./TCPChat $port")
	flag.Parse()
	app.Server = server.NewTcpServer("localhost:" + *port)
	app.Server.Start()
}
