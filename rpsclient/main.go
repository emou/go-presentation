package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type Message struct {
	Params map[string]string
}

func writeMsg(writer *bufio.Writer, msg map[string]string) error {
	b, err := json.Marshal(Message{Params: msg})
	if err != nil {
		return err
	}
	_, err = writer.WriteString(string(b) + "\n")
	if err != nil {
		return err
	}
	return writer.Flush()
}

func readMsg(reader *bufio.Reader) (map[string]string, error) {
	msgString, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	msgString = strings.TrimSpace(msgString)

	msg := Message{}
	err = json.Unmarshal([]byte(msgString), &msg)
	if err != nil {
		return nil, err
	}

	return msg.Params, nil
}

func login(writer *bufio.Writer, user string, password string) {
	msg := map[string]string{"name": "login", "username": user, "password": password}
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
		switch msg["name"] {
		case "login":
			login(writer, user, password)
		default:
			fmt.Printf("Ignoring unknown message: %s\n", msg)
		}
	}
}

func main() {
	serverAddr := "localhost:9000"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(fmt.Sprintf("Error connectinng to %s: %s", serverAddr, err))
	}
	fmt.Printf("Connected to %s!\n", serverAddr)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	clientLoop("emo", "123", reader, writer)
}
