package rps

import (
	"strings"
)

const STATE_NEW = "new"
const STATE_JOINED = "joined"
const STATE_WAITING = "waiting"
const STATE_PLAYING = "playing"

type Player struct {
	Name     string
	State    string
	action   chan string
	Messages chan Message
	Finish   chan string
}

func NewPlayer(name string) *Player {
	return &Player{Name: name, action: make(chan string), Messages: make(chan Message), Finish: make(chan string), State: STATE_NEW}
}

// Act dispatches a user action
func (p *Player) Act(message string) {
	p.action <- strings.TrimSpace(message)
}

// WriteMsg sends a message to the user
func (p *Player) WriteMsg(message Message) {
	p.Messages <- message
}

func (p *Player) Leave() {
	close(p.Finish)
}
