package helpers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
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
	reader := bufio.NewReader(conn)
	return reader.ReadString('\n')
}

func MessageFromUser(name string) string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	return fmt.Sprintf("[%s][%s]: ", formattedTime, name)
}

func MessageFromServer(message string) string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	return fmt.Sprintf("[%s][%s]: ", formattedTime, message)
}
