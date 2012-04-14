package state

import (
	"bytes"
	"encoding/gob"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

const stockTypes = 6
const startingValue = 100

type Player struct {
	Cash     uint64
	Shares   [stockTypes]uint64
	Password []byte
}

type Stock struct {
	Name  string
	Value uint64
}

type gameState struct {
	Key    []byte
	Stock  [stockTypes]Stock
	Player map[string]*Player
}

type Game struct {
	g gameState
	sync.Mutex
}

type PlayerInfo struct {
	Cash   uint64
	Shares [stockTypes]uint64
	p      *Player
	g      *Game
}

func (p *PlayerInfo) SetPassword(pw []byte) {
	p.g.Lock()
	p.p.Password = pw
	p.g.Unlock()
}

func (p *PlayerInfo) CheckPassword(pw []byte) bool {
	p.g.Lock()
	defer p.g.Unlock()
	return bytes.Equal(p.p.Password, pw)
}

func (g *Game) ListStocks() []Stock {
	g.Lock()
	defer g.Unlock()
	rv := make([]Stock, len(g.g.Stock))
	copy(rv, g.g.Stock[:])
	return rv
}

func (g *Game) HasPlayer(name string) bool {
	g.Lock()
	_, ok := g.g.Player[name]
	g.Unlock()
	return ok
}

func (g *Game) Player(name string) *PlayerInfo {
	g.Lock()
	defer g.Unlock()

	if _, ok := g.g.Player[name]; !ok {
		p := &Player{Cash: 500000}
		g.g.Player[name] = p
	}
	p := g.g.Player[name]

	return &PlayerInfo{Cash: p.Cash, Shares: p.Shares, p: p, g: g}
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

func New(data string) *Game {
	year, month, day := time.Now().Date()
	rand.Seed(int64(year)*1000 + int64(month)*100 + int64(day))

	var g Game
	f, err := os.Open(data)
	if err == nil {
		defer f.Close()
		err = gob.NewDecoder(f).Decode(&g.g)
		if err == nil {
			return &g
		}
	}

	// File not found or gob invalid
	g.g.Player = make(map[string]*Player)
	for i := 0; i < stockTypes; i++ {
		g.g.Stock[i].Value = startingValue
		g.g.Stock[i].Name = g.g.pickName()
	}
	g.g.newKey()

	f, err = os.Create(data)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(&g.g)

	return &g
}
