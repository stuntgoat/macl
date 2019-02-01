package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"net/http"
)

var (
	GAMES              *GamesContainer
	LOGGER             *log.Logger
	API_PREFIX         = flag.String("api_prefix", "game", "api URL prefix")
	NUM_PLAYERS        = flag.Int("num_players", 2, "required number of players")
	BOARD_WIDTH        = flag.Int("board_width", 4, "board width")
	BOARD_LENGTH       = flag.Int("board_length", 4, "board length")
	CONSECUTIVE_LENGTH = flag.Int("consecutive_length", 4,
		"consecutive line length required for a win")
	LOG_PATH = flag.String("log_path", "macl.log", "logging path")
	PORT     = flag.Int("port", 8080, "server port")
)

func init() {
	flag.Parse()

	logfile, err := os.OpenFile(*LOG_PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("unable to open log file")
	}
	prefix := fmt.Sprintf("[%s] ", *API_PREFIX)
	LOGGER = log.New(logfile, prefix, log.LstdFlags|log.Lshortfile)

	GAMES = &GamesContainer{
		games: map[string]*game{},
	}

}

func main() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *PORT),
		Handler: configureRouter(*API_PREFIX),
	}

	LOGGER.Println(fmt.Sprintf("serving on port: %d", *PORT))
	err := server.ListenAndServe()
	if err != nil {
		panic("Error: " + err.Error())
	}
}
