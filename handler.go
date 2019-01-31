package main

import (
	"fmt"
	"net/http"
)

func writeJSON(w http.ResponseWriter, content []byte) {
	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(content)))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	var content []byte
	var APIerr *APIError

	if r.Method == "GET" {
		content, APIerr = API_getGameList()
	} else if r.Method == "POST" {
		content, APIerr = API_createGame(r)
	} else {
		APIerr = &APIError{"method not allowed", 400}
	}

	if APIerr != nil {
		LOGGER.Println(fmt.Sprintf("error in game handler %s", APIerr))
		http.Error(w, APIerr.Msg, APIerr.Status)
		return
	}
	writeJSON(w, content)
}

func gameStatusHandler(w http.ResponseWriter, r *http.Request) {
	var content []byte
	var APIerr *APIError
	content, APIerr = API_gameStatus(r)
	if APIerr != nil {
		LOGGER.Println(fmt.Sprintf("error getting game status %s", APIerr.Msg))
		http.Error(w, APIerr.Msg, APIerr.Status)
		return
	}
	writeJSON(w, content)
}

func moveListHandler(w http.ResponseWriter, r *http.Request) {
	var content []byte
	var APIerr *APIError
	content, APIerr = API_moveList(r)
	if APIerr != nil {
		LOGGER.Println(fmt.Sprintf("error getting move list %s", APIerr.Msg))
		http.Error(w, APIerr.Msg, APIerr.Status)
		return
	}
	writeJSON(w, content)
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	var content []byte
	var APIerr *APIError
	content, APIerr = API_getMove(r)
	if APIerr != nil {
		LOGGER.Println(fmt.Sprintf("error getting move %s", APIerr.Msg))
		http.Error(w, APIerr.Msg, APIerr.Status)
		return
	}
	writeJSON(w, content)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusBadRequest
	if r.Method == "DELETE" {
		// Quit game.
		status = API_quitGame(r)
	} else if r.Method == "POST" {
		// Make move.
		content, APIerr := API_makeMove(r)
		if APIerr != nil {
			LOGGER.Println(fmt.Sprintf("error quiting game %s", APIerr.Msg))
			http.Error(w, APIerr.Msg, APIerr.Status)
			return
		}
		writeJSON(w, content)
		return
	}
	w.WriteHeader(status)
}
