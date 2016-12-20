package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"net"
	"os"
	"strings"
)

type Message struct {
	Type   string
	Params map[string]string
}

func writeMsg(writer *bufio.Writer, msg *Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(string(b) + "\n")
	if err != nil {
		return err
	}
	return writer.Flush()
}

func readMsg(reader *bufio.Reader) (*Message, error) {
	msgString, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	msgString = strings.TrimSpace(msgString)

	msg := &Message{}
	err = json.Unmarshal([]byte(msgString), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func login(writer *bufio.Writer, user string, password string) {
	msg := &Message{
		Type:   "login",
		Params: map[string]string{"username": user, "password": password},
	}
	err := writeMsg(writer, msg)
	if err != nil {
		panic(fmt.Sprintf("Error logging in: %s", err))
	}
	fmt.Println("Logged in, waiting for a game...")
}

func clientLoop(user string, password string, reader *bufio.Reader, writer *bufio.Writer) {
	for {
		msg, err := readMsg(reader)
		if err != nil {
			fmt.Printf("Lost connection to server: %s\n", err)
			return
		}
		switch msg.Type {
		case "login":
			login(writer, user, password)
		default:
			fmt.Printf("Ignoring unknown message: %s\n", msg)
		}
	}
}

func getLogin() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	fmt.Print("Password: ")
	password, err := gopass.GetPasswd()
	if err != nil {
		return "", "", err
	}
	return username, string(password), nil
}

func main() {
	serverAddr := "localhost:9000"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(fmt.Sprintf("Error connectinng to %s: %s", serverAddr, err))
	}
	fmt.Printf("Connected to %s!\n", serverAddr)
	username, password, err := getLogin()

	if err != nil {
		panic(fmt.Sprintf("Error getting login information: %s", err))
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	clientLoop(username, password, reader, writer)
}
