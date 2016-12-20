package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/emou/go-presentation/gorps/rps"
	"net"
	"strings"
)

type Message struct {
	Params map[string]string
}

func writeMsg(writer *bufio.Writer, params map[string]string) error {
	msg := Message{Params: params}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	b = append(b, '\n')

	_, err = writer.Write(b)
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
	if err != nil {
		return nil, err
	}

	msg := &Message{}
	err = json.Unmarshal([]byte(msgString), msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func login(reader *bufio.Reader, writer *bufio.Writer) (string, error) {
	for {
		err := writeMsg(writer, map[string]string{"name": "login"})
		if err != nil {
			return "", err
		}
		msg, err := readMsg(reader)
		if err != nil {
			return "", err
		}
		return msg.Params["username"], nil
	}
}

func serve(conn net.Conn, game *rps.Game) error {
	defer conn.Close()

	fmt.Printf("Incoming connection: %+v\n", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	name, err := login(reader, writer)
	if err != nil {
		fmt.Printf("Error on login: %s", name)
	}
	player := rps.NewPlayer(name)
	game.AddPlayer(player)

	// TODO: Pull function
	go func() {
		fmt.Println("Opening player message feed ", player.Name)
		defer fmt.Println("Closing player message feed ", player.Name)

		for {
			select {
			case msg := <-player.Messages:
				_, err := writer.Write([]byte(msg))
				if err != nil {
					fmt.Println(err)
					return
				}
				writer.Flush()
			case <-player.Finish:
				return
			}
		}
	}()

	game.StartMatch(player)

	// TODO: Pull function
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			game.RemovePlayer(player)
			return err
		}

		if player.State == rps.STATE_PLAYING && message != "" {
			player.Act(message)
		}

	}
}

func main() {
	game := rps.NewGame()
	listener, err := net.Listen("tcp", ":9000")

	fmt.Println("Listening on", listener.Addr())
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err == nil {
			go serve(conn, game)
		} else {
			fmt.Printf("%+v", err)
		}
	}

}
