package rps

import "fmt"

const addPlayer, removePlayer, startMatch, endMatch = "removePlayer", "addPlayer", "startMatch", "removeMatch"

// Action struct holds a player action
type Action struct {
	player *Player
	action string
	param  string
}

func (a Action) String() string {
	return fmt.Sprintf("<%s:%s:%s>", a.player.Name, a.action, a.param)
}

// Game struct holds the action channel and players
type Game struct {
	actions chan Action
	players []*Player
	waiting []*Player
	scores  *ScoreBoard
}

// NewGame starts a new game
func NewGame() *Game {
	scoreboard := NewScoreBoard()
	var game = Game{
		players: []*Player{},
		waiting: []*Player{},
		actions: make(chan Action),
		scores:  scoreboard}
	go game.cmdloop()
	return &game
}

// cmdloop starts a command processing loop, reading from the game command channel
func (game *Game) cmdloop() {

	for {
		var action = <-game.actions
		fmt.Printf("Handle: %s\n", action)

		switch action.action {
		case addPlayer:
			game.addPlayer(action)
		case removePlayer:
			game.removePlayer(action)
		case startMatch:
			game.startMatch(action)
		case endMatch:
			game.endMatch(action)
		}
	}
}

// addPlayer handles action, that adds a player to the game
func (game *Game) addPlayer(action Action) {
	fmt.Printf("Joining: %s\n", action.player.Name)
	action.player.State = STATE_JOINED
	game.players = append(game.players, action.player)
}

// AddPlayer dispatches action to adds a player to the game
func (game Game) AddPlayer(player *Player) {
	game.actions <- Action{player: player, action: addPlayer, param: ""}
}

// removePlayer removes a player from the game
func (game *Game) removePlayer(action Action) {
	game.removeWaiting(action.player)
	for idx, player := range game.players {
		if player == action.player {
			fmt.Printf("Leaving: %s\n", action.player.Name)
			game.players = append(game.players[:idx], game.players[idx+1:]...)
			break
		}
	}
	action.player.Leave()
}

// RemovePlayer dispatches action to remove a player from the game
func (game Game) RemovePlayer(player *Player) {
	game.actions <- Action{player: player, action: removePlayer, param: ""}
}

// startMatch handles request to start a new match
func (game *Game) startMatch(action Action) {
	waiting := game.getOpponent(action.player)
	if waiting == nil {
		game.addWaiting(action.player)
		fmt.Printf("Player %s waiting ...", action.player.Name)
		action.player.State = STATE_WAITING
	} else {

		var msg = fmt.Sprintf("Duel between %s and %s. FIGHT!\n", waiting.Name, action.player.Name)
		waiting.WriteMsg(msg)
		action.player.WriteMsg(msg)

		waiting.State = STATE_PLAYING
		action.player.State = STATE_PLAYING

		match := NewMatch(game, action.player, waiting)
		go match.start()
	}
}

// StartMatch dispatches action to start a new match
func (game Game) StartMatch(player *Player) {
	game.actions <- Action{player: player, action: startMatch, param: ""}
}

// endMatch handles request to end a new match
func (game *Game) endMatch(action Action) {
	action.player.WriteMsg(fmt.Sprintf("You %s \n", action.param))
	win := action.param == "win"
	game.scores.Add(action.player.Name, win)
	game.startMatch(action)
}

// EndMatch ends a match
func (game Game) EndMatch(player *Player, param string) {
	game.actions <- Action{player: player, action: endMatch, param: param}
}

// getOpponent returns an opponent for this player or nil
func (game *Game) getOpponent(player *Player) *Player {
	for idx, waiting := range game.waiting {
		if waiting == player || waiting.Name == player.Name {
			continue
		}
		game.waiting = append(game.waiting[:idx], game.waiting[idx+1:]...)
		return waiting
	}
	return nil
}

// addWaiting adds a player to the waitlist
func (game *Game) addWaiting(player *Player) {
	game.waiting = append(game.waiting, player)
}

// removeWaiting removes a player from the waitlist
func (game *Game) removeWaiting(player *Player) {
	for idx, waiting := range game.waiting {
		if waiting == player {
			game.waiting = append(game.waiting[:idx], game.waiting[idx+1:]...)
			return
		}
	}
}
