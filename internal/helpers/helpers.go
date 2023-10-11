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

	return fmt.Sprintf("\n[%s][%s]: ", formattedTime, name)
}

func NewInput(name string) string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	return fmt.Sprintf("[%s][%s]: ", formattedTime, name)
}

func MessageFromServerNewConnect(name string) string {
	return fmt.Sprintf("\nUser %s connected to chat ...\n", name)
}

func MessageFromServerLeftFromChat(name string) string {
	return fmt.Sprintf("\nUser %s left from chat ...\n", name)
}
func MessageFromServerIncorrectInput() string {
	return fmt.Sprintf("Incorrect Input\n")
}

func MessageIsValid(message string) bool {
	messInRune := []rune(message)
	if len(messInRune) == 0 {
		return false
	}
	for _, w := range messInRune {
		if w < 32 || w > 127 {
			return false
		}
	}
	return true
}

func NameIsValid(message string) bool {
	messInRune := []rune(message)
	if len(messInRune) == 0 || len(messInRune) > 15 {
		return false
	}
	for _, w := range messInRune {
		if w < 32 || w > 127 {
			return false
		}
	}
	return true
}
