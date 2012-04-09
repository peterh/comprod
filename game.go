package main

import (
	"math/rand"
	"time"
)

const stockTypes = 6
const startingValue = 100

type Player struct {
	Cash   uint64
	Shares [stockTypes]uint64
}

type Stock struct {
	Name  string
	Value uint64
}

type gameState struct {
	Stock  [stockTypes]Stock
	Player map[string]*Player
}

func (g *gameState) listStocks() []Stock {
	return g.Stock[:]
}

func (g *gameState) player(name string) *Player {
	if p, ok := g.Player[name]; ok {
		return p
	}
	return &Player{Cash: 500000}
}

func (g *gameState) pickName() string {
	names := [...]string{"Coffee", "Soybeans", "Corn", "Wheat", "Cocoa", "Gold", "Silver", "Platinum", "Oil", "Natural Gas", "Cotton", "Sugar"}
	used := make(map[string]bool)
	for _, v := range g.Stock {
		used[v.Name] = true
	}
	for {
		i := rand.Intn(len(names))
		if !used[names[i]] {
			return names[i]
		}
	}
	return ""
}

func newGame() *gameState {
	year, month, day := time.Now().Date()
	rand.Seed(int64(year)*1000 + int64(month)*100 + int64(day))

	var g gameState
	g.Player = make(map[string]*Player)
	for i := 0; i < stockTypes; i++ {
		g.Stock[i].Value = startingValue
		g.Stock[i].Name = g.pickName()
	}
	return &g
}

