package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"encoding/json"
	"io/ioutil"
	"net/http"
)


type MoveResponse struct {
	Type MoveType `json:"type"`
	Player string `json:"player"`

	Column int `json:"column,omitempty"`

}
type MovesRangeResponse struct {
	Moves []MoveResponse `json:"moves"`
}

type MovesRangeRequest struct {
	Start int
	Until int
}

type GameStatusResponse struct {
	Players []string `json:"players"`
	Status GameStatus `json:"state"`
	Winner string `json:"winner,omitempty"`
}

type CreateGameRequest struct {
	Players []string `json:"players"`
	Columns int `json:"columns"`
	Rows int `json:"rows"`
}

type CreateGameResponse struct {
	GameId string `json:"gameId"`
}
type MoveRequest struct {
	Column int `json:"column"`
}

type GameList struct {
	Games []string `json:"games"`
}

// validateMoveList returns a range between 0 and -1, where -1 means to the end of the list.
func validateMoveList(r *http.Request) (*MovesRangeRequest, error) {
	var err error

	vals := r.URL.Query()
	if len(vals) == 0 {
		return &MovesRangeRequest{0, -1}, nil
	}
	start := 0
	until := -1

	startStrings, ok := vals["start"]
	if ok {
		startStr := strings.TrimSpace(startStrings[0])
		start, err = strconv.Atoi(startStr)
		if err != nil {
			// Invalid query parameter value.
			return nil, errors.New("invalid start conversion")
		}
	}
	untilStrings, ok := vals["until"]
	LOGGER.Println(fmt.Sprintf("%s %b", untilStrings, ok))
	if ok {
		untilStr := strings.TrimSpace(untilStrings[0])
		until, err = strconv.Atoi(untilStr)
		if err != nil {
			// Invalid query parameter value.
			return nil, errors.New("invalid until conversion")
		}
	}
	if start < 0 || start > until {
		LOGGER.Println(fmt.Sprintf("%+v", vals))
		return nil, errors.New("bad range request")
	}

	rangeReq := &MovesRangeRequest{
		Start: start,
		Until: until,
	}

	return rangeReq, nil
}

// validateCreateGame takes an
func validateMakeMove(r *http.Request) (*MoveRequest, *APIError)  {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		LOGGER.Println(fmt.Sprintf("failed to read body %s", err))
		return nil, &APIError{"server error", http.StatusInternalServerError}
	}
	mr := &MoveRequest{}
	err = json.Unmarshal(b, mr)
	if err != nil {
		return nil, &APIError{"malformed input", http.StatusBadRequest}
	}
	return mr, nil
}

// validateCreateGame takes an
func validateCreateGame(r *http.Request) (*CreateGameRequest, *APIError)  {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, &APIError{"server error", http.StatusInternalServerError}
	}
	cgr := &CreateGameRequest{}
	err = json.Unmarshal(b, cgr)
	if err != nil {
		return nil, &APIError{"malformed input", http.StatusBadRequest}
	}
	if len(cgr.Players) != *NUM_PLAYERS {
		return nil, &APIError{fmt.Sprintf("num players is not %d", *NUM_PLAYERS),
			http.StatusBadRequest}
	}
	if cgr.Rows != *BOARD_WIDTH || cgr.Columns != *BOARD_LENGTH {
		msg := fmt.Sprintf("expecting 4 rows and 4 columns, got rows %d cols %d",
			cgr.Rows, cgr.Columns)
		return nil, &APIError{msg, http.StatusBadRequest}
	}
	return cgr, nil
}
