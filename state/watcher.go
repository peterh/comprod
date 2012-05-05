package state

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"
)

func (g *Game) write(fn string) {
	g.Lock()
	defer g.Unlock()

	new := fn + ".new"
	old := fn + ".old"
	f, err := os.OpenFile(fn+".new", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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

func (g *Game) nextTurn() <-chan time.Time {
	g.Lock()
	defer g.Unlock()

	now := time.Now().UTC()
	prev := g.g.Previous

	if prev.Day() != now.Day() {
		return time.After(1)
	}

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

	g.Lock()
	defer g.Unlock()

	now := time.Now().UTC()
	defer func() { g.g.Previous = now }()

	before := g.g.Stock
	var divpaid [stockTypes]uint64
	news := make([]string, 0, stockTypes)

	for i := 0; i < rounds; i++ {
		adjust := uint64(math.Pow(rand.Float64()*.8+1.2, 5.0))
		stock := rand.Intn(stockTypes)
		switch rand.Intn(3) {
		case up:
			g.g.Stock[stock].Value += adjust
			if g.g.Stock[stock].Value >= splitValue {
				news = append(news, g.g.Stock[stock].Name+" split 2 for 1")
				g.g.Stock[stock].Value = (g.g.Stock[stock].Value + 1) / 2
				before[stock].Value = (before[stock].Value + 1) / 2
				for _, p := range g.g.Player {
					p.Shares[stock] *= 2
				}
			}
		case down:
			if g.g.Stock[stock].Value <= adjust {
				news = append(news, g.g.Stock[stock].Name+" went bankrupt, and was removed from the market")
				for _, p := range g.g.Player {
					p.Shares[stock] = 0
				}
				g.g.Stock[stock].Value = startingValue
				before[stock].Value = startingValue
				newname := g.g.pickName()
				news = append(news, newname+" was added to the market")
				g.g.Stock[stock].Name = newname
			} else {
				g.g.Stock[stock].Value -= adjust
			}
		case dividend:
			if g.g.Stock[stock].Value >= startingValue {
				divpaid[stock] += adjust
				for _, p := range g.g.Player {
					p.Cash += adjust * p.Shares[stock]
				}
			}
		}
	}

	for k, v := range g.g.Stock {
		var item string
		switch {
		case v.Value == before[k].Value:
			item = v.Name + " did not change price"
		case v.Value < before[k].Value:
			item = fmt.Sprintf("%s fell %.1f%%", v.Name, float64(before[k].Value-v.Value)/float64(before[k].Value)*100)
		default: // case v.Value > before[k].Value:
			item = fmt.Sprintf("%s rose %.1f%%", v.Name, float64(v.Value-before[k].Value)/float64(before[k].Value)*100)
		}
		if divpaid[k] > 0 {
			item = fmt.Sprintf("%s, and paid $%d in dividends", item, divpaid[k])
		}
		news = append(news, item)
	}
	g.g.News = news

	if now.Month() != g.g.Previous.Month() {
		leader := g.g.Leaders()
		sort.Sort(LeaderSort(leader))
		announce := fmt.Sprintf("The winner of the %s %d season was %s, with a net worth of $%d",
			g.g.Previous.Month().String(), g.g.Previous.Year(), leader[0].Name, leader[0].Worth)
		g.g.News = append(g.g.News, announce)
		g.g.History = append([]string{announce}, g.g.History...)
		if len(leader) > 1 {
			rup := []string{}
			for i := 1; i < len(leader); i++ {
				rup = append(rup, fmt.Sprintf("%s had $%d", leader[i].Name, leader[i].Worth))
			}
			g.g.News = append(g.g.News,
				fmt.Sprintf("(%s)", strings.Join(rup, ", ")))
		}
		g.g.reset()
	}
}

func watcher(g *Game, filename string, changed chan struct{}) {
	var tick <-chan time.Time
	tock := g.nextTurn()

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
			tock = g.nextTurn()
		case <-sigint:
			if tick != nil {
				g.write(filename)
			}
			log.Println("Exiting")
			os.Exit(0)
		}
	}
}
