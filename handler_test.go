package main

import (
	"errors"
	"fmt"
	"flag"
	"strings"
	"testing"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)


func mockUUID() string {
	return "cats"
}

func expectWithWriter(w *httptest.ResponseRecorder, status int, expected string) error {
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != status {
		return errors.New(fmt.Sprintf("expected %d got %d", status, resp.StatusCode))
	}


	got := strings.TrimSpace(string(body))
	if got != expected {
		return errors.New(fmt.Sprintf("expected %s got %s", expected, string(body)))
	}
	return nil
}

func TestMain(m *testing.M) {
	// Mock the game id so game id in tests are consistent.
	flag.Parse()
	old := mkGameId
	mkGameId = mockUUID
	m.Run()
	mkGameId = old
}

func apiURL(resource string) string {
	if resource != "" {
		return fmt.Sprintf(fmt.Sprintf("http://localhost/%s/%s", *API_PREFIX, resource))
	}
	return fmt.Sprintf(fmt.Sprintf("http://localhost/%s", *API_PREFIX))
}

func Test_gameHandler(t *testing.T) {
	createGameBlob := strings.NewReader(`{"players": ["a", "b"],"rows": 4, "columns": 4}`)

	// create game
	r := httptest.NewRequest("POST", apiURL(""), createGameBlob)
	w := httptest.NewRecorder()

	gameHandler(w, r)
	err := expectWithWriter(w, http.StatusOK, `{"gameId":"cats"}`)
	if err != nil {
		t.Error(err)
	}

	// return in progress game.
	r = httptest.NewRequest("GET", apiURL(""), createGameBlob)
	w = httptest.NewRecorder()

	gameHandler(w, r)
	err = expectWithWriter(w, http.StatusOK, `{"games":["cats"]}`)
	if err != nil {
		t.Error(err)
	}

	// bad input to create game
	createGameBlob = strings.NewReader(`{"players": ["a", "b"],"rows": 4, "columns": z}`)
	r = httptest.NewRequest("POST", apiURL(""), createGameBlob)
	w = httptest.NewRecorder()

	gameHandler(w, r)
	err = expectWithWriter(w, http.StatusBadRequest, `malformed input`)
	if err != nil {
		t.Error(err)
	}

	// only one player
	createGameBlob = strings.NewReader(`{"players": ["a"], "rows": 4, "columns": 4}`)
	r = httptest.NewRequest("POST", apiURL(""), createGameBlob)
	w = httptest.NewRecorder()

	gameHandler(w, r)
	err = expectWithWriter(w, http.StatusBadRequest, `num players is not 2`)
	if err != nil {
		t.Error(err)
	}

	// invalid row count
	createGameBlob = strings.NewReader(`{"players": ["a", "b"], "rows": 30, "columns": 4}`)
	r = httptest.NewRequest("POST", apiURL(""), createGameBlob)
	w = httptest.NewRecorder()

	gameHandler(w, r)
	err = expectWithWriter(w, http.StatusBadRequest,
		`expecting 4 rows and 4 columns, got rows 30 cols 4`)
	if err != nil {
		t.Error(err)
	}

	// invalid col count
	createGameBlob = strings.NewReader(`{"players": ["a", "b"], "rows": 4, "columns": 30}`)
	r = httptest.NewRequest("POST", apiURL(""), createGameBlob)
	w = httptest.NewRecorder()

	gameHandler(w, r)
	err = expectWithWriter(w, http.StatusBadRequest,
		`expecting 4 rows and 4 columns, got rows 4 cols 30`)
	if err != nil {
		t.Error(err)
	}


}


func Test_gameStatusHandler(t *testing.T) {
	g := CreateGame(4, 4, 4, "a", "b")
	GAMES.Add(g)

	// check status
	r := httptest.NewRequest("GET", apiURL("cats"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats"})

	w := httptest.NewRecorder()
	gameStatusHandler(w, r)

	err := expectWithWriter(w, http.StatusOK, `{"players":["a","b"],"state":"IN_PROGRESS"}`)
	if err != nil {
		t.Error(err)
	}

	// Check draw response
	g = CreateGame(4, 4, 4, "a", "b")
	GAMES.Add(g)
	mkDraw(g, "a", "b")

	r = httptest.NewRequest("GET", apiURL("cats"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats"})

	w = httptest.NewRecorder()
	gameStatusHandler(w, r)

	err = expectWithWriter(w, http.StatusOK, `{"players":["a","b"],"state":"DONE"}`)
	if err != nil {
		t.Error(err)
	}

	// wrong game
	g = CreateGame(4, 4, 4, "a", "b")
	GAMES.Add(g)
	mkDraw(g, "a", "b")

	r = httptest.NewRequest("GET", apiURL("dogs"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "dogs"})

	w = httptest.NewRecorder()
	gameStatusHandler(w, r)

	err = expectWithWriter(w, http.StatusNotFound, `unknown game`)
	if err != nil {
		t.Error(err)
	}


	// Check winner response
	g = CreateGame(4, 4, 4, "a", "b")
	GAMES.Add(g)
	g.Move("a", 3)
	g.Move("b", 1)
	g.Move("a", 3)
	g.Move("b", 2)
	g.Move("a", 3)
	g.Move("b", 0)
	g.Move("a", 3)

	r = httptest.NewRequest("GET", apiURL("cats"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats"})

	w = httptest.NewRecorder()
	gameStatusHandler(w, r)

	err = expectWithWriter(w, http.StatusOK, `{"players":["a","b"],"state":"DONE","winner":"a"}`)
	if err != nil {
		t.Error(err)
	}
}


func Test_moveListHandler(t *testing.T) {
	g := CreateGame(4, 4, 4, "a", "b")
	GAMES.Add(g)
	// No moves
	r := httptest.NewRequest("GET", apiURL("/cats/moves"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats"})

	w := httptest.NewRecorder()
	moveListHandler(w, r)
	err := expectWithWriter(w, http.StatusOK, `{"moves":[]}`)
	if err != nil {
		t.Error(err)
	}

	g = CreateGame(4, 4, 4, "a", "b")
	GAMES.Add(g)
	g.Move("a", 3)
	g.Move("b", 1)

	// get all moves
	r = httptest.NewRequest("GET", apiURL("/cats/moves"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats"})

	w = httptest.NewRecorder()
	moveListHandler(w, r)
	err = expectWithWriter(w, http.StatusOK,
		`{"moves":[{"type":"MOVE","player":"a","column":3},{"type":"MOVE","player":"b","column":1}]}`)
	if err != nil {
		t.Error(err)
	}

	g = CreateGame(4, 4, 4, "a", "b")
	GAMES.Add(g)
	g.Move("a", 3)
	g.Move("b", 1)
	g.Move("a", 2)
	g.Move("b", 2)

	// move range
	r = httptest.NewRequest("GET", apiURL("/cats/moves?start=1&until=2"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats"})

	w = httptest.NewRecorder()
	moveListHandler(w, r)
	err = expectWithWriter(w, http.StatusOK,
		`{"moves":[{"type":"MOVE","player":"b","column":1},{"type":"MOVE","player":"a","column":2}]}`)
	if err != nil {
		t.Error(err)
	}

	// wrong game
	r = httptest.NewRequest("GET", apiURL("/dogs/moves?start=1&until=2"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "dogs"})

	w = httptest.NewRecorder()
	moveListHandler(w, r)
	err = expectWithWriter(w, http.StatusNotFound, `game not found`)
	if err != nil {
		t.Error(err)
	}

	// bad input start
	r = httptest.NewRequest("GET", apiURL("/cats/moves?start=Y&until=2"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats"})

	w = httptest.NewRecorder()
	moveListHandler(w, r)
	err = expectWithWriter(w, http.StatusBadRequest, `invalid start conversion`)
	if err != nil {
		t.Error(err)
	}

	// bad input until
	r = httptest.NewRequest("GET", apiURL("/cats/moves?start=1&until=Y"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats"})

	w = httptest.NewRecorder()
	moveListHandler(w, r)
	err = expectWithWriter(w, http.StatusBadRequest, `invalid until conversion`)
	if err != nil {
		t.Error(err)
	}

	// bad range
	r = httptest.NewRequest("GET", apiURL("/cats/moves?start=3&until=2"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats"})

	w = httptest.NewRecorder()
	moveListHandler(w, r)
	err = expectWithWriter(w, http.StatusBadRequest, `bad range request`)
	if err != nil {
		t.Error(err)
	}


}


func Test_moveHandler(t *testing.T) {
	g := CreateGame(4, 4, 4, "a", "b")
	GAMES.Add(g)
	g.Move("a", 3)
	g.Move("b", 1)

	r := httptest.NewRequest("GET", apiURL("cats/moves/1"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "move_number": "1"})

	w := httptest.NewRecorder()
	moveHandler(w, r)
	err := expectWithWriter(w, http.StatusOK, `{"type":"MOVE","player":"b","column":1}`)
	if err != nil {
		t.Error(err)
	}

	// Test missing game
	r = httptest.NewRequest("GET", apiURL("DOGS/moves/1"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "DOGS", "move_number": "1"})

	w = httptest.NewRecorder()
	moveHandler(w, r)
	err = expectWithWriter(w, http.StatusNotFound, `game not found`)
	if err != nil {
		t.Error(err)
	}

	// Test missing move
	r = httptest.NewRequest("GET", apiURL("cats/moves/999"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "move_number": "999"})

	w = httptest.NewRecorder()
	moveHandler(w, r)
	err = expectWithWriter(w, http.StatusNotFound, `invalid index`)
	if err != nil {
		t.Error(err)
	}

	// bad input
	r = httptest.NewRequest("GET", apiURL("cats/moves/dogs"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "move_number": "dogs"})

	w = httptest.NewRecorder()
	moveHandler(w, r)
	err = expectWithWriter(w, http.StatusBadRequest, `unable to parse move number`)
	if err != nil {
		t.Error(err)
	}
}


func Test_playHandler(t *testing.T) {
	g := CreateGame(4, 4, 4, "a", "b")
	GAMES.Add(g)

	playGameBlob := strings.NewReader(`{"column" : 2}`)
	r := httptest.NewRequest("POST", apiURL("cats/a"), playGameBlob)

	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "playerId": "a"})

	w := httptest.NewRecorder()
	playHandler(w, r)
	err := expectWithWriter(w, http.StatusOK, `{"move":"cats/moves/0"}`)
	if err != nil {
		t.Error(err)
	}


	playGameBlob = strings.NewReader(`{"column" : 1}`)
	r = httptest.NewRequest("POST", apiURL("cats/b"), playGameBlob)

	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "playerId": "b"})
	w = httptest.NewRecorder()

	playHandler(w, r)
	err = expectWithWriter(w, http.StatusOK, `{"move":"cats/moves/1"}`)
	if err != nil {
		t.Error(err)
	}

	// Wrong game
	playGameBlob = strings.NewReader(`{"column" : 1}`)
	r = httptest.NewRequest("POST", apiURL("dogs/b"), playGameBlob)
	r = mux.SetURLVars(r, map[string]string{"gameId": "dogs", "playerId": "b"})

	w = httptest.NewRecorder()

	playHandler(w, r)
	err = expectWithWriter(w, http.StatusNotFound, "unknown game")
	if err != nil {
		t.Error(err)
	}

	// play out of turn
	playGameBlob = strings.NewReader(`{"column" : 1}`)
	r = httptest.NewRequest("POST", apiURL("cats/b"), playGameBlob)

	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "playerId": "b"})

	w = httptest.NewRecorder()
	playHandler(w, r)
	err = expectWithWriter(w, http.StatusConflict, `WRONG_TURN`)
	if err != nil {
		t.Error(err)
	}

	// wrong player
	playGameBlob = strings.NewReader(`{"column" : 1}`)
	r = httptest.NewRequest("POST", apiURL("cats/dog"), playGameBlob)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "playerId": "dog"})

	w = httptest.NewRecorder()
	playHandler(w, r)
	err = expectWithWriter(w, http.StatusBadRequest, `WRONG_GAME`)
	if err != nil {
		t.Error(err)
	}


	// bad input
	playGameBlob = strings.NewReader(`{"column" : "forty-two"}`)
	r = httptest.NewRequest("POST", apiURL("cats/b"), playGameBlob)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "playerId": "b"})

	w = httptest.NewRecorder()
	playHandler(w, r)
	err = expectWithWriter(w, http.StatusBadRequest, `malformed input`)
	if err != nil {
		t.Error(err)
	}

	// bad input
	playGameBlob = strings.NewReader(`{"column" : -1}`)
	r = httptest.NewRequest("POST", apiURL("cats/b"), playGameBlob)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "playerId": "b"})

	w = httptest.NewRecorder()
	playHandler(w, r)
	err = expectWithWriter(w, http.StatusBadRequest, `BAD_REQUEST`)
	if err != nil {
		t.Error(err)
	}

	// quit
	r = httptest.NewRequest("DELETE", apiURL("cats/b"), nil)

	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "playerId": "b"})

	w = httptest.NewRecorder()
	playHandler(w, r)
	err = expectWithWriter(w, http.StatusAccepted, ``)
	if err != nil {
		t.Error(err)
	}

	// quit wrong game as player
	r = httptest.NewRequest("DELETE", apiURL("cats/dogs"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "playerId": "dogs"})

	w = httptest.NewRecorder()
	playHandler(w, r)
	err = expectWithWriter(w, http.StatusNotFound, ``)
	if err != nil {
		t.Error(err)
	}

	// quit wrong game id
	r = httptest.NewRequest("DELETE", apiURL("dogs/a"), nil)
	r = mux.SetURLVars(r, map[string]string{"gameId": "dogs", "playerId": "a"})

	w = httptest.NewRecorder()
	playHandler(w, r)
	err = expectWithWriter(w, http.StatusNotFound, ``)
	if err != nil {
		t.Error(err)
	}


	// quit twice
	r = httptest.NewRequest("DELETE", apiURL("cats/b"), nil)

	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "playerId": "b"})

	w = httptest.NewRecorder()
	playHandler(w, r)
	err = expectWithWriter(w, http.StatusGone, ``)
	if err != nil {
		t.Error(err)
	}

	// quit by other player
	r = httptest.NewRequest("DELETE", apiURL("cats/a"), nil)

	r = mux.SetURLVars(r, map[string]string{"gameId": "cats", "playerId": "a"})

	w = httptest.NewRecorder()
	playHandler(w, r)
	err = expectWithWriter(w, http.StatusGone, ``)
	if err != nil {
		t.Error(err)
	}
}
