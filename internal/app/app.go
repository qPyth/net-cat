package app

import (
	"fmt"
	"net-cat/internal/server"
	"os"
)

type App struct {
	Server *server.TcpServer
}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() {
	var port string
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	} else if len(os.Args) == 2 {
		port = os.Args[1]
	} else {
		port = "8989"
	}

	app.Server = server.NewTcpServer("localhost:" + port)
	app.Server.Start()
}
