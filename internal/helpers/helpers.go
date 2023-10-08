package helpers

import (
	"log"
	"net"
	"os"
)

var (
	pathToGreetingImage = "assets/greetings.txt"
)

func Greeting() string {
	image, err := os.ReadFile(pathToGreetingImage)
	if err != nil {
		log.Fatal(err)
	}
	return string(image)
}

// TODO need a writer and reader

func Write(conn net.Conn, message string) error {
	bytes := []byte(message)
	_, err := conn.Write(bytes)
	return err
}

func Read(conn net.Conn) (string, error) {
	buffer := make([]byte, 4)
	_, err := conn.Read(buffer)
	return string(buffer), err
}
