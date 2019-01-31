package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
)

type GameStatus string
var STATUS_DONE = GameStatus("DONE")
var STATUS_IN_PROGRESS = GameStatus("IN_PROGRESS")
var STATUS_INVALID_GAME = GameStatus("INVALID_GAME")
var STATUS_GAME_OVER = GameStatus("GAME_OVER")
var STATUS_QUIT_LEFT_GAME = GameStatus("QUIT_LEFT_GAME")
var STATUS_LEFT_GAME = GameStatus("LEFT_GAME")

type MoveStatus string
var MoveOK = MoveStatus("OK")
var MoveBadRequest = MoveStatus("BAD_REQUEST")
var MoveWrongGame = MoveStatus("WRONG_GAME")
var MoveWrongTurn = MoveStatus("WRONG_TURN")

type MoveType string
var MoveMove = MoveType("MOVE")
var MoveQuit = MoveType("QUIT")

type Move struct  {
	player string
	row int
	col int

	Type MoveType
}

type MoveConfirmation struct {
	Move string `json:"move"`
}
func MkConfirmation(id string, moveNum int) *MoveConfirmation {
	return &MoveConfirmation{
		Move: fmt.Sprintf("%s/moves/%d", id, moveNum),
	}
}

var mkGameId = func () string {
	u, _ := uuid.NewV4()
	return fmt.Sprintf("%v", u)
}

type game struct {
	sync.RWMutex

	id string

	board [][]string

	// Players and status regarding if they are still playing this game.
	players map[string]bool

	playerList []string

	// Sequential list of moves
	moves []*Move

	// Location of player coins on the board.
	// playerId to PlayerGraph
	playerGraphs map[string]*PlayerGraph

	// If this game is over.
	over bool

	// Player id of the winner.
	winner string

	sequentialWin int
}


func (g *game) GameStatus() *GameStatusResponse {
	g.RLock()
	defer g.RUnlock()

	var status GameStatus

	if g.over {
		status = STATUS_DONE
	} else {
		status = STATUS_IN_PROGRESS
	}
	gameStatus := &GameStatusResponse{
		Players: g.currentlyPlaying(),
		Status: status,
	}
	if status == STATUS_DONE {
		gameStatus.Winner = g.winner
	}
	return gameStatus
}


func (g *game) GetMove(index int) (*Move, error) {
	g.RLock()
	defer g.RUnlock()
	var move *Move
	if index < 0 || index >= len(g.moves) {
		return nil, errors.New("invalid index")
	}
	move = g.moves[index]
	return move, nil
}

func (g *game) GetMoves(start, until int) []*Move {
	g.RLock()
	defer g.RUnlock()
	if until + 1 > len(g.moves) || until < 0 {
		until = len(g.moves)
	} else {
		until = until + 1
	}
	moves := []*Move{}
	for _, m := range g.moves[start:until] {
		moves = append(moves, m)
	}
	return moves
}

func (g *game) boardIsFull() bool {
	for _, item := range g.board[0] {
		if item == "" {
			return false
		}
	}
	return true
}

// makeMove performs the move on the board and sets related status.
func (g *game) makeMove(playerId string, col int) MoveStatus {
	spot := g.board[0][col]
	if spot != "" {
		// Column is full
		return MoveBadRequest
	}

	lastEmptyRow := len(g.board) - 1

	// Attempt to move on board
	for rowIdx := 1; rowIdx < len(g.board); rowIdx++ {
		spot = g.board[rowIdx][col]
		if spot != "" {
			lastEmptyRow = rowIdx - 1
			break
		}
	}

	g.board[lastEmptyRow][col] = playerId
	g.moves = append(g.moves, &Move{playerId, lastEmptyRow, col, MoveMove})

	playerGraph := g.playerGraphs[playerId]
	playerGraph.Add(lastEmptyRow, col)

	won := playerGraph.FindConsecutive(lastEmptyRow, col, g.sequentialWin)
	if won {
		g.winner = playerId
		g.over = true
	}
	if g.boardIsFull() {
		g.over = true
	}
	return MoveOK
}

// Move returns an error if there was a problem with the move.
func (g *game) Move(playerId string, col int) (*MoveConfirmation, MoveStatus) {
	g.Lock()
	defer g.Unlock()

	// Validate this column
	if col < 0 || col > len(g.board[0]) - 1 {
		return nil, MoveBadRequest
	}

	ok := g.isPlaying(playerId)
	if !ok {
		// Player is NOT playing this game.
		return nil, MoveWrongGame
	}

	if g.over {
		return nil, MoveBadRequest
	}

	// Check if it's my turn
	if g.nextMove() != playerId {
		return nil, MoveWrongTurn
	}
	status := g.makeMove(playerId, col)

	// Add response
	// {
	// 	"move": "{gameId}/moves/{move_number}"
	// }
	confirmation := MkConfirmation(g.id, len(g.moves) - 1)
	return confirmation, status
}

func (g *game) isDone() bool {
	g.RLock()
	defer g.RUnlock()
	return g.over
}

func (g *game) Winner() string {
	g.RLock()
	defer g.RUnlock()
	return g.winner
}

func (g *game) isPlaying(playerId string) bool {
	return g.players[playerId]
}

// currentlyPlaying returns an ORDERED list of players that have not quit.
func (g *game) currentlyPlaying() []string {
	players := []string{}
	for _, player := range g.playerList {
		if !g.players[player] {
			continue
		}
		players = append(players, player)
	}
	return players
}

func (g *game) Quit(playerId string) GameStatus {

	g.Lock()
	defer g.Unlock()

	playing, ok := g.players[playerId]
	if !ok {
		return STATUS_INVALID_GAME
	}

	if g.over {
		return STATUS_GAME_OVER
	}

	if !playing {
		return STATUS_QUIT_LEFT_GAME
	}

	// Can quit now.
	g.players[playerId] = false
	playersLeft := g.currentlyPlaying()
	if len(playersLeft) == 1 {
		g.over = true
		g.winner = playersLeft[0]
	}
	g.moves = append(g.moves, &Move{
		player: playerId,
		Type: MoveQuit,
	})
	return STATUS_LEFT_GAME
}

// NextMove returns the playerId of the user who has the next move.
func (g *game) nextMove() string {
	players := g.currentlyPlaying()
	if len(g.moves) == 0 {
		return players[0]
	}

	lastMove := g.moves[len(g.moves) - 1]
	var player string
	for i := 0; i < len(players); i++ {
		if players[i] == lastMove.player {
			player = players[(i + 1) % len(players)]
			break
		}
	}
	return player
}

func CreateGame(winningSequence, rows, cols int, players ...string) *game {
	g := &game{}
	g.sequentialWin = winningSequence

	board := [][]string{}
	for i := 0; i < rows; i++ {
		row := []string{}
		for j := 0; j < cols; j++ {
			row = append(row, "")
		}
		board = append(board, row)
	}
	g.board = board
	g.id = mkGameId()
	playerMap := map[string]bool{}
	graphs := map[string]*PlayerGraph{}
	for _, player := range players {
		playerMap[player] = true
		graphs[player] = &PlayerGraph{
			coins: map[CoinKey]bool{},
		}
	}
	g.players = playerMap

	g.playerList = players
	g.moves = []*Move{}
	g.playerGraphs = graphs
	return g
}
