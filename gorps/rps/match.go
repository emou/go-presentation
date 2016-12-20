package rps

import (
	"fmt"
	"github.com/satori/uuid"
	"math"
)

var ROCK = 0
var PAPER = 1
var SCISSORS = 2

const MSG_CHOOSE = "CHOOSE"

var MSG_MAP = map[int]string{
	0: "rock",
	1: "paper",
	2: "scissors",
}

var ACT_MAP = map[string]int{
	"1": ROCK,
	"2": PAPER,
	"3": SCISSORS,
}

// msgPlayers messages all players with a message
func msgPlayers(msg Message, players ...*Player) {
	for _, player := range players {
		player.WriteMsg(msg)
	}
}

// Match struct holds the match players and decisions
type Match struct {
	id        uuid.UUID
	game      *Game
	player1   *Player
	player2   *Player
	decision1 *int
	decision2 *int
}

// NewMatch creates a new match
func NewMatch(game *Game, player1 *Player, player2 *Player) *Match {
	return &Match{id: uuid.NewV4(), game: game, player1: player1, player2: player2}
}

// handleUserAction handles a player
func (m Match) handlePlayerAction(player *Player, action string) *int {
	act, found := ACT_MAP[action]
	if !found {
		return nil
	}
	return &act
}

// gameLogic
func gameLogic(decision1 int, decision2 int) bool {
	if decision1 == decision2 {
		return false
	}
	if math.Mod(float64(decision1+3)-float64(decision2), 3) == 1 {
		return true
	}
	return false
}

// checkWinner
func (m *Match) checkWinner() (*Player, *Player) {
	switch {
	case gameLogic(*m.decision1, *m.decision2):
		return m.player1, m.player2
	case gameLogic(*m.decision2, *m.decision1):
		return m.player2, m.player1
	default:
		return nil, nil
	}
}

// start starts a new match
func (m *Match) start() {
	fmt.Println("Begin match ", m.id)
	defer fmt.Println("End match ", m.id)
	msgPlayers(Message{Type: MSG_CHOOSE}, m.player1, m.player2)

	for {
		select {
		case action := <-m.player1.action:
			m.decision1 = m.handlePlayerAction(m.player1, action)
			if m.decision1 == nil {
				msgPlayers(Message{Type: MSG_CHOOSE}, m.player1)
			}

		case action := <-m.player2.action:
			m.decision2 = m.handlePlayerAction(m.player2, action)
			if m.decision2 == nil {
				msgPlayers(Message{Type: MSG_CHOOSE}, m.player2)
			}
		case <-m.player1.Finish:
			m.game.EndMatch(m.player2, "win")
			return
		case <-m.player2.Finish:
			m.game.EndMatch(m.player1, "win")
			return
		default:
			if m.decision1 != nil && m.decision2 != nil {
				msgPlayers(Message{
					Type: "choices",
					Params: map[string]string{
						m.player1.Name: MSG_MAP[*m.decision1],
						m.player2.Name: MSG_MAP[*m.decision2],
					},
				}, m.player1, m.player2)

				winner, loser := m.checkWinner()
				if winner == nil && loser == nil {
					msgPlayers(Message{
						Type: "tie",
					})
					msgPlayers(Message{
						Type: MSG_CHOOSE,
					}, m.player1, m.player2)

					m.decision1 = nil
					m.decision2 = nil
				} else {
					m.game.EndMatch(winner, "win")
					m.game.EndMatch(loser, "lose")
					return
				}
			}
		}
	}
}
