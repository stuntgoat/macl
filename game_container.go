package main

import (
	"sync"
)

type GamesContainer struct {
	sync.RWMutex
	games map[string]*game
}

func (gc *GamesContainer) Get(gameId string) (*game, bool) {
	gc.RLock()
	defer gc.RUnlock()
	g, ok := gc.games[gameId]

	return g, ok
}

func (gc *GamesContainer) GetGames() []string {
	gc.RLock()
	defer gc.RUnlock()
	glist := []string{}
	for key, _ := range gc.games {
		glist = append(glist, key)
	}

	return glist
}

func (gc *GamesContainer) Add(g *game) {
	gc.Lock()
	defer gc.Unlock()
	gc.games[g.id] = g
}
