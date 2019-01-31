package main

import (
	"fmt"
	"github.com/gorilla/mux"
)

// configureRouter allows leading URI customization for the API routes.
func configureRouter(custom string) *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	// GET in-progress games.
	// POST new game.
	r.HandleFunc(fmt.Sprintf("/%s", custom), gameHandler).Methods("GET", "POST")

	// Status of a game.
	r.HandleFunc(fmt.Sprintf("/%s/{gameId}", custom), gameStatusHandler).Methods("GET")

	// Get all or some moves in a game.
	r.HandleFunc(fmt.Sprintf("/%s/{gameId}/moves", custom), moveListHandler).Methods("GET")

	// Query a move number
	r.HandleFunc(
		fmt.Sprintf("/%s/{gameId}/moves/{move_number}", custom), moveHandler).Methods("GET")

	// POST move
	// DELETE quit
	r.HandleFunc(
		fmt.Sprintf("/%s/{gameId}/{playerId}", custom), playHandler).Methods("POST", "DELETE")

	return r
}
