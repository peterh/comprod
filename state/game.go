package state

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
	Stock  [stockTypes]Stock
	Player map[string]*Player
	News   []string
	Key    []byte
}

type Game struct {
	g gameState
	sync.Mutex
	changed chan<- struct{}
}

type PlayerInfo struct {
	Cash   uint64
	Shares [stockTypes]uint64
	p      *Player
	g      *Game
}

var ping struct{}

func (g *gameState) findStock(stock string) int {
	for k, v := range g.Stock {
		if v.Name == stock {
			return k
		}
	}
	return -1
}

func (p *PlayerInfo) Buy(stock string, lots uint64) error {
	p.g.Lock()
	defer p.g.Unlock()

	idx := p.g.g.findStock(stock)
	if idx < 0 {
		return fmt.Errorf("%s is not on the market", stock)
	}

	shares := lots * 100
	afford := p.p.Cash / p.g.g.Stock[idx].Value
	if shares > afford {
		return fmt.Errorf("You don't have enough cash to buy %d shares of %s", shares, stock)
	}

	p.p.Cash -= shares * p.g.g.Stock[idx].Value
	p.p.Shares[idx] += shares

	p.g.changed <- ping

	// Update caller-visible copies
	p.Cash = p.p.Cash
	p.Shares[idx] = p.p.Shares[idx]
	return nil
}

func (p *PlayerInfo) Sell(stock string, lots uint64) error {
	p.g.Lock()
	defer p.g.Unlock()

	idx := p.g.g.findStock(stock)
	if idx < 0 {
		return fmt.Errorf("%s is not on the market", stock)
	}

	shares := lots * 100
	if shares > p.p.Shares[idx] {
		return fmt.Errorf("You don't have %d shares of %s to sell", shares, stock)
	}

	p.p.Cash += shares * p.g.g.Stock[idx].Value
	p.p.Shares[idx] -= shares

	p.g.changed <- ping

	// Update caller-visible copies
	p.Cash = p.p.Cash
	p.Shares[idx] = p.p.Shares[idx]
	return nil
}

func (p *PlayerInfo) SetPassword(pw []byte) {
	p.g.Lock()
	p.p.Password = pw
	p.g.Unlock()
	p.g.changed <- struct{}{}
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
		p := &Player{Cash: 100000}
		g.g.Player[name] = p
		g.changed <- ping
	}
	p := g.g.Player[name]

	return &PlayerInfo{Cash: p.Cash, Shares: p.Shares, p: p, g: g}
}

func (g *Game) News() []string {
	g.Lock()
	defer g.Unlock()
	return g.g.News
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
	changed := make(chan struct{})
	g.changed = changed

	f, err := os.Open(data)
	if err == nil {
		defer f.Close()
		err = gob.NewDecoder(f).Decode(&g.g)
		if err == nil {
			go watcher(&g, data, changed)
			return &g
		}
	}

	// File not found or gob invalid
	g.g.Player = make(map[string]*Player)
	for i := 0; i < stockTypes; i++ {
		g.g.Stock[i].Value = startingValue
		g.g.Stock[i].Name = g.g.pickName()
	}
	g.g.News = []string{"A new season started\n"}
	g.g.newKey()

	go watcher(&g, data, changed)
	g.changed <- ping

	return &g
}
