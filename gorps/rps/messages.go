package rps

import (
	"bufio"
	"encoding/json"
	"strings"
)

type Message struct {
	Params map[string]string
}

func WriteMsg(writer *bufio.Writer, params map[string]string) error {
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

func ReadMsg(reader *bufio.Reader) (*Message, error) {
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
