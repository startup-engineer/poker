package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Suit string

const (
	Spade Suit = "♠️"
	Club       = "♣️"
	Heart      = "♥️"
	Diamond    = "♦️"
)

// Card

var faceCardStrings = [4]string{"J", "Q", "K", "A"}

type Card struct {
	Value int
	Suit Suit
}

func (c Card) String() string {
	// TODO: Update string methods to use string buffers instead of inefficient
	// string concatenation

	valueString := ""
	if c.Value > 1 && c.Value < 11 {
		valueString = fmt.Sprintf("%d", c.Value)
	} else {
		valueString = faceCardStrings[c.Value%11]
	}
	return fmt.Sprintf("%s%s", valueString, c.Suit)
}

// Deck
type Deck struct {
	Cards []Card
}

func (d Deck) String() string {
	// TODO: Update string methods to use string buffers instead of inefficient
	// string concatenation

	// Empty deck
	if len(d.Cards) == 0 {
		return "{}"
	}

	// Print deck in rows of four for visual clarity
	s := fmt.Sprintf("%s  ", d.Cards[0])
	i := 2
	for _, c := range d.Cards[1:] {
		s += fmt.Sprintf("%s", c)
		if i % 4 == 0 {
			s += "\n"
		} else {
			s += "  "
		}
		i++
	}
	return s
}

func (d *Deck) Init() {
	d.Cards = make([]Card, 52)
	suits := []Suit{ Spade, Club, Heart, Diamond }
	i := 0
	for v := 2; v <= 14; v++ {
		for _, s := range suits {
			d.Cards[i].Value = v
			d.Cards[i].Suit = s
			i++
		}
	}
}

func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

// Player
type Player struct {
	Name string
	Cards []Card
}

func (p Player) String() string {
	if len(p.Cards) == 2 {
		return fmt.Sprintf("%s, Cards: %s  %s", p.Name, p.Cards[0], p.Cards[1])
	} else {
		return p.Name
	}
}

// Poker Table
type PokerTable struct {
	Players []Player
	CommunityCards []Card
	BurnedCards []Card
	Deck Deck
}

func (pt *PokerTable) Init(playerCount int) {
	pt.Players = make([]Player, playerCount)
	for i := range pt.Players {
		pt.Players[i].Name = fmt.Sprintf("Player #%d", i)
	}

	pt.Deck.Init()
	pt.Deck.Shuffle()

	for i := 0; i < 2*len(pt.Players); i++ {
		pt.Players[i%len(pt.Players)].Cards = append(
			pt.Players[i%len(pt.Players)].Cards,
			pt.Deck.Cards[0],
		)
		pt.Deck.Cards = pt.Deck.Cards[1:]
	}

	pt.BurnedCards = append(pt.BurnedCards, pt.Deck.Cards[0]) // Burn one
	pt.Deck.Cards = pt.Deck.Cards[1:]

	pt.CommunityCards = pt.Deck.Cards[:3] // Flop
	pt.Deck.Cards = pt.Deck.Cards[3:]

	pt.BurnedCards = append(pt.BurnedCards, pt.Deck.Cards[0]) // Burn one
	pt.Deck.Cards = pt.Deck.Cards[1:]

	pt.CommunityCards = append(pt.CommunityCards, pt.Deck.Cards[0]) // Turn
	pt.Deck.Cards = pt.Deck.Cards[1:]

	pt.BurnedCards = append(pt.BurnedCards, pt.Deck.Cards[0]) // Burn one
	pt.Deck.Cards = pt.Deck.Cards[1:]

	pt.CommunityCards = append(pt.CommunityCards, pt.Deck.Cards[0]) // Turn
	pt.Deck.Cards = pt.Deck.Cards[1:]
}

func (pt PokerTable) String() string {
	// TODO: Update string methods to use string buffers instead of inefficient
	// string concatenation

	s := ""
	for _, p := range pt.Players {
		s += fmt.Sprintf("%s\n", p)
	}
	s += "\n"

	s += fmt.Sprintf("Community Cards: ")
	for i, c := range pt.CommunityCards {
		if i != len(pt.CommunityCards) - 1 {
			s += fmt.Sprintf("%s  ", c)
		} else {
			s += fmt.Sprintf("%s", c)
		}
	}
	s += fmt.Sprintf("\n\n")

	s += fmt.Sprintf("Burned Cards: ")
	for i, c := range pt.BurnedCards {
		if i != len(pt.BurnedCards) - 1 {
			s += fmt.Sprintf("%s  ", c)
		} else {
			s += fmt.Sprintf("%s", c)
		}
	}
	s += fmt.Sprintf("\n\n")

	s += fmt.Sprintf("Deck:\n")
	s += fmt.Sprintf("%s\n", pt.Deck)

	return s
}

func main() {
	pokerTable := PokerTable{}
	pokerTable.Init(9)

	fmt.Println(pokerTable)
}