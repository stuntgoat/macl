package main

import (
	"bytes"
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"strconv"
	"github.com/gorilla/mux"
)

type APIError struct {
	Msg string
	Status int
}
func (e *APIError) Error() string {
	return e.Msg
}


func API_getMove(r *http.Request) ([]byte, *APIError) {
	vars := mux.Vars(r)
	gid := vars["gameId"]
	g, ok := GAMES.Get(gid)
	if !ok {
		// 404 - Game not found or player is not a part of it.
		return nil, &APIError{"game not found", http.StatusNotFound}
	}
	moveNumStr := vars["move_number"]
	moveNumStr = strings.TrimSpace(moveNumStr)
	moveNum, err := strconv.Atoi(moveNumStr)
	if err != nil {
		return nil, &APIError{"unable to parse move number", http.StatusBadRequest}
	}

	move, err := g.GetMove(moveNum)
	if err != nil {
		return nil, &APIError{err.Error(), http.StatusNotFound}
	}

	mr := mkMoveResponse(move)

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(mr)
	if err != nil {
		LOGGER.Println(fmt.Sprintf("Encode Error: %s", err))
		return nil, &APIError{fmt.Sprintf("JSON Encode error"), http.StatusInternalServerError}
	}
	return buf.Bytes(), nil
}

func mkMoveResponse(move *Move) MoveResponse {
	mRes := MoveResponse{
			Type: move.Type,
			Player: move.player,
		}
		if move.Type == MoveMove {
			mRes.Column = move.col
		}
	return mRes
}

func moveResponses(moves []*Move) *MovesRangeResponse {
	mrr := &MovesRangeResponse{
		Moves: []MoveResponse{},
	}

	for _, m := range moves {
		mRes := mkMoveResponse(m)
		mrr.Moves = append(mrr.Moves, mRes)
	}
	return mrr
}


func API_moveList(r *http.Request) ([]byte, *APIError) {
	vars := mux.Vars(r)
	gid := vars["gameId"]
	g, ok := GAMES.Get(gid)
	if !ok {
		return nil, &APIError{"game not found", http.StatusNotFound}
	}

	rangeReq, err := validateMoveList(r)
	if err != nil {
		return nil, &APIError{err.Error(), http.StatusBadRequest}
	}
	moves := g.GetMoves(rangeReq.Start, rangeReq.Until)

	mRangeRes := moveResponses(moves)

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(mRangeRes)
	if err != nil {
		LOGGER.Println(fmt.Sprintf("JSON Encode error: %s", err))
		return nil, &APIError{"server error", http.StatusInternalServerError}
	}
	return buf.Bytes(), nil
}


func API_makeMove(r *http.Request) ([]byte, *APIError) {
	vars := mux.Vars(r)
	gid := vars["gameId"]
	g, ok := GAMES.Get(gid)
	if !ok {
		return nil, &APIError{"unknown game", http.StatusNotFound}
	}

	mr, APIerr := validateMakeMove(r)
	if APIerr != nil {
		return nil, APIerr
	}

	confirmation, status := g.Move(vars["playerId"], mr.Column)

	var err error
	switch status {
	case MoveOK:
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		err = enc.Encode(confirmation)
		if err != nil {
			LOGGER.Println(fmt.Sprintf("JSON Encode error: %s", err))
			return nil, &APIError{"server error", http.StatusInternalServerError}
		}
		return buf.Bytes(), nil

	case MoveWrongGame:
		return nil, &APIError{string(status), http.StatusBadRequest}
	case MoveBadRequest:
		return nil, &APIError{string(status), http.StatusBadRequest}
	case MoveWrongTurn:
		return nil, &APIError{string(status), http.StatusConflict}
	default:
		return nil, &APIError{string(status), http.StatusNotFound}
	}
}


func API_quitGame(r *http.Request) int {
	vars := mux.Vars(r)
	g, ok := GAMES.Get(vars["gameId"])
	if !ok {
		return http.StatusNotFound
	}

	gameStatus := g.Quit(vars["playerId"])
	switch gameStatus {
	case STATUS_INVALID_GAME:
		return http.StatusNotFound
	case STATUS_QUIT_LEFT_GAME:
		// NOTE: only valid when more than 2 players.
		return http.StatusNotFound
	case STATUS_GAME_OVER:
		return http.StatusGone
	case STATUS_LEFT_GAME:
		return http.StatusAccepted
	default:
		return http.StatusNotFound
	}
}



func API_gameStatus(r *http.Request) ([]byte, *APIError) {
	vars := mux.Vars(r)

	g, ok := GAMES.Get(vars["gameId"])
	if !ok {
		return nil, &APIError{"unknown game", http.StatusNotFound}
	}
	status := g.GameStatus()

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(status)
	if err != nil {
		LOGGER.Println(fmt.Sprintf("error encoding JSON %s", err))
		return nil, &APIError{"server error", http.StatusInternalServerError}
	}
	return buf.Bytes(), nil
}

func API_getGameList() ([]byte, *APIError) {
	glist := GAMES.GetGames()

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(&GameList{
		Games: glist,
	})
	if err != nil {
		LOGGER.Println(fmt.Sprintf("error encoding JSON %s", err))
		return nil, &APIError{"server error", http.StatusInternalServerError}
	}
	return buf.Bytes(), nil

}

// API_createGame validates a request to create a game and returns the JSON response or
// an error on failure.
func API_createGame(r *http.Request) ([]byte, *APIError) {
	cgr, APIerr := validateCreateGame(r)
	if APIerr != nil {
		return nil, APIerr
	}

	game := CreateGame(*CONSECUTIVE_LENGTH, cgr.Rows, cgr.Columns, cgr.Players...)

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(&CreateGameResponse{game.id})
	if err != nil {
		LOGGER.Println(fmt.Sprintf("JSON Encode error: %s", err))
		return nil, &APIError{"server error", http.StatusInternalServerError}
	}
	GAMES.Add(game)
	return buf.Bytes(), nil
}
