package main

import (

	"fmt"
	"testing"
)

// NOTE: only works for 4x4 board.
func mkDraw(g *game, playerA, playerB string) {
	for i := 0; i < 2; i++ {
		for j := 0; j < 4; j++ {
			if j < 2 {
				g.Move(playerA, 0 + j)
				g.Move(playerB, 2 + j)
			} else {
				g.Move(playerA, j)
				g.Move(playerB, j - 2)
			}
		}
	}
}


func Test_CreateGame(t *testing.T) {
	g := CreateGame(4, 4, 4, "a", "b")

	if g.sequentialWin != 4 {
		t.Error("expected sequentialWin to be 4")
	}

	for _, row := range g.board {
		for _, space := range row {
			if space != "" {
				t.Error("expected space to be empty string")
			}
		}
	}

	if g.id == "" {
		t.Error("expected game id to be a uuid found ", g.id)
	}

	if !g.players["a"] {
		t.Error("expected player a to be playing")
	}
	if !g.players["b"] {
		t.Error("expected player b to be playing")
	}

	if g.playerGraphs["a"] == nil {
		t.Error("expected initialized map")
	}
	if g.playerGraphs["b"] == nil {
		t.Error("expected initialized map")
	}

	if g.playerList[0] != "a" {
		t.Error("expected first player to be a")
	}
	if g.playerList[1] != "b" {
		t.Error("expected first player to be b")
	}

	if len(g.moves) != 0 {
		t.Error("expected moves to be empty")
	}

	if g.over {
		t.Error("expected game to be in play")
	}

	if g.winner != "" {
		t.Error("expected no winner")
	}
}
func Test_Quit(t *testing.T) {
	g := CreateGame(4, 4, 4, "a", "b", "c")

	if g.Quit("YYYYYY") != STATUS_INVALID_GAME {
		t.Error("expected STATUS_INVALID_GAME")
	}

	if g.Quit("b") != STATUS_LEFT_GAME {
		t.Error("expected STATUS_LEFT_GAME")
	}

	if g.Quit("b") != STATUS_QUIT_LEFT_GAME {
		t.Error("expected STATUS_QUIT_LEFT_GAME")
	}

	if g.Quit("c") != STATUS_LEFT_GAME {
		t.Error("expected STATUS_LEFT_GAME")
	}

	if g.Quit("a") != STATUS_GAME_OVER {
		t.Error("expected STATUS_GAME_OVER")
	}
}
func Test_boardIsFull(t *testing.T) {
	g := CreateGame(4, 4, 4, "a", "b")
	if g.boardIsFull() {
		t.Error("expected board not to be full")
	}
	for i := 0; i < 4; i++ {
		g.board[0][i] = "YYYYYY"
	}
	if !g.boardIsFull() {
		t.Error("expected board to be full")
	}

	_, status := g.Move("a", 1)
	if status != MoveBadRequest {
		t.Error("expected MoveBadRequest got", status)
	}
}


func Test_Move(t *testing.T) {
	g := CreateGame(4, 4, 4, "a", "b")
	_, status := g.Move("a", -1)
	if status != MoveBadRequest {
		t.Error("expected MoveBadRequest")
	}
	_, status = g.Move("YYYYY", 0)
	if status != MoveWrongGame {
		t.Error("expected MoveWrongGame")
	}
	_, status = g.Move("b", 0)
	if status != MoveWrongTurn {
		t.Error("expected MoveWrongTurn")
	}

	g.Move("a", 1)

	move := g.moves[0]
	if move.player != "a" {
		t.Error("expected player a move")
	}

	if move.row != 3 {
		t.Error("expected move row to be at index 3 got ", move.row)
	}
	if move.col != 1 {
		t.Error("expected move col to be at index 1 got ", move.row)
	}

	var confirmation *MoveConfirmation
	confirmation, status = g.Move("b", 2)
	if status != MoveOK {
		t.Error("expected MoveOK")
	}
	expected := fmt.Sprintf("%s/moves/1", g.id)
	if confirmation.Move != expected {
		t.Error("expected %s", expected)

	}
	_, status = g.Move("a", 1)
	if status != MoveOK {
		t.Error("expected MoveOK")
	}

	_, status = g.Move("b", 3)
	if status != MoveOK {
		t.Error("expected MoveOK")
	}

	_, status = g.Move("a", 1)
	if status != MoveOK {
		t.Error("expected MoveOK")
	}

	_, status = g.Move("b", 3)
	if status != MoveOK {
		t.Error("expected MoveOK")
	}

	_, status = g.Move("a", 1)
	if status != MoveOK {
		t.Error("expected MoveOK")
	}

	_, status = g.Move("a", 1)
	if status != MoveBadRequest {
		t.Error("expected MoveBadRequest")
	}

	if g.Winner() != "a" {
		t.Error("expected a as winner")
	}

	if !g.isDone() {
		t.Error("expected game over")
	}
}

func Test_Draw(t *testing.T) {
	g := CreateGame(4, 4, 4, "a", "b")
	mkDraw(g, "a", "b")

	if !g.isDone() {
		t.Error("expected game to be over")
	}
	if g.Winner() != "" {
		t.Error("expected no winner got", g.Winner())
	}
}

func Test_GetMove(t *testing.T) {
	g := CreateGame(4, 4, 4, "a", "b")
	g.Move("a", 1)
	_, err := g.GetMove(1)
	if err != nil && err.Error() != "invalid index" {
		t.Error("expected invalid index got", err)
	}
	if err == nil {
		t.Error("expected invalid index got nil")
	}
}

func Test_GetMoves(t *testing.T) {
	g := CreateGame(4, 4, 4, "a", "b")
	g.Move("a", 1)
	g.Move("b", 2)
	g.Move("a", 3)

	moves := g.GetMoves(0, 1)
	if len(moves) != 2 {
		t.Error("expected 2 moves got", len(moves))
	}

	if moves[0].col != 1 {
		t.Error("expected column 1 got", moves[0].col)
	}
	if moves[1].col != 2 {
		t.Error("expected column 2 got", moves[1].col)
	}

	moves = g.GetMoves(0, 0)
	if len(moves) != 1 {
		t.Error("expected 1 move got", len(moves))
	}
	moves = g.GetMoves(0, 2)
	if len(moves) != 3 {
		t.Error("expected 3 moves got", len(moves))
	}

	moves = g.GetMoves(0, 3)
	if len(moves) != 3 {
		t.Error("expected 3 moves got", len(moves))
	}

	if moves[2].col != 3 {
		t.Error("expected column 3 got", moves[2].col)
	}


	moves = g.GetMoves(1, 2)
	if len(moves) != 2 {
		t.Error("expected 2 moves got", len(moves))
	}

	if moves[1].col != 3 {
		t.Error("expected column 3 got", moves[1].col)
	}

	moves = g.GetMoves(0, -1)
	if len(moves) != 3 {
		t.Error("expected 3 moves got", len(moves))
	}

	moves = g.GetMoves(1, 1)
	if len(moves) != 1 {
		t.Error("expected 1 move got", len(moves))
	}

	if moves[0].col != 2 {
		t.Error("expected column 2 got", moves[0].col)
	}

}
