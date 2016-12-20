package main

import (
	"bufio"
	"fmt"
	"github.com/emou/go-presentation/gorps/rps"
	"net"
)

func login(reader *bufio.Reader, writer *bufio.Writer) (string, error) {
	for {
		err := rps.WriteMsg(writer, &rps.Message{Type: "login"})
		if err != nil {
			return "", err
		}
		msg, err := rps.ReadMsg(reader)
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
				err := rps.WriteMsg(writer, &msg)
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
		msg, err := rps.ReadMsg(reader)
		if err != nil {
			game.RemovePlayer(player)
			return err
		}

		if msg.Type != "action" {
			game.RemovePlayer(player)
			return rps.WriteMsg(writer, &rps.Message{
				Type:   "error",
				Params: map[string]string{"message": fmt.Sprintf("Unexpected message: %+v", msg)},
			})
		}

		if player.State == rps.STATE_PLAYING && msg != nil {
			player.Act(msg.Params["action"])
		}
	}
}

func main() {
	game := rps.NewGame()
	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err == nil {
			go serve(conn, game)
		} else {
			fmt.Printf("%+v", err)
		}
	}

}
