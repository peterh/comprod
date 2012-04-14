package state

import (
	"encoding/gob"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func (g *Game) write(fn string) {
	g.Lock()
	defer g.Unlock()

	new := fn + ".new"
	old := fn + ".old"
	f, err := os.Create(fn + ".new")
	if err != nil {
		log.Fatal(err)
	}
	err = gob.NewEncoder(f).Encode(&g.g)
	if err != nil {
		log.Fatal(err)
	}
	f.Sync()
	f.Close()

	os.Rename(fn, old)
	err = os.Rename(new, fn)
	if err != nil {
		log.Fatal(err)
	}
	os.Remove(old)
}

func nextTurn() <-chan time.Time {
	now := time.Now().UTC()
	tomorrow := now.Add(time.Hour * 23)
	for tomorrow.Day() == now.Day() {
		tomorrow = tomorrow.Add(time.Hour)
	}
	next := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.UTC)
	return time.After(next.Sub(now))
}

func (g *Game) newDay() {
	const rounds = 15
	const (
		up = iota
		down
		dividend
	)

	for i := 0; i < rounds; i++ {
		adjust := uint64(math.Pow(rand.Float64()*.8+1.2, 5.0))
		stock := rand.Intn(stockTypes)
		switch rand.Intn(3) {
		case up:
			g.g.Stock[stock].Value += adjust
		case down:
			g.g.Stock[stock].Value -= adjust
		case dividend:
			if g.g.Stock[stock].Value >= startingValue {
				for _, p := range g.g.Player {
					p.Cash += adjust * p.Shares[stock]
				}
			}
		}
	}
}

func watcher(g *Game, filename string, changed chan struct{}) {
	var tick <-chan time.Time
	tock := nextTurn()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, os.Kill)

	for {
		select {
		case <-changed:
			if tick == nil {
				tick = time.After(5 * time.Minute)
			}
		case <-tick:
			tick = nil
			g.write(filename)
		case <-tock:
			g.newDay()
			g.write(filename)
			tock = nextTurn()
		case <-sigint:
			if tick != nil {
				g.write(filename)
			}
			log.Println("Exiting")
			os.Exit(0)
		}
	}
}
