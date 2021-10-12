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

func CardValueToString(v int) string {
	valueString := ""
	if v > 1 && v < 11 {
		valueString = fmt.Sprintf("%d", v)
	} else {
		valueString = faceCardStrings[v%11]
	}
	return valueString
}

func (c Card) String() string {
	// TODO: Update string methods to use string buffers instead of inefficient
	// string concatenation

	vs := CardValueToString(c.Value)
	return fmt.Sprintf("%s%s", vs, c.Suit)
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
type PokerGame struct {
	Players []Player
	CommunityCards []Card
	BurnedCards []Card
	Deck Deck
}

func (pg *PokerGame) Init(playerCount int) {
	pg.Players = make([]Player, playerCount)
	for i := range pg.Players {
		pg.Players[i].Name = fmt.Sprintf("Player #%d", i)
	}

	pg.Deck.Init()
	pg.Deck.Shuffle()

	// Deal cards to players
	for i := 0; i < 2*len(pg.Players); i++ {
		pg.Players[i%len(pg.Players)].Cards = append(
			pg.Players[i%len(pg.Players)].Cards,
			pg.Deck.Cards[0],
		)
		pg.Deck.Cards = pg.Deck.Cards[1:]
	}

	pg.BurnedCards = append(pg.BurnedCards, pg.Deck.Cards[0]) // Burn one
	pg.Deck.Cards = pg.Deck.Cards[1:]

	pg.CommunityCards = pg.Deck.Cards[:3] // Flop
	pg.Deck.Cards = pg.Deck.Cards[3:]

	pg.BurnedCards = append(pg.BurnedCards, pg.Deck.Cards[0]) // Burn one
	pg.Deck.Cards = pg.Deck.Cards[1:]

	pg.CommunityCards = append(pg.CommunityCards, pg.Deck.Cards[0]) // Turn
	pg.Deck.Cards = pg.Deck.Cards[1:]

	pg.BurnedCards = append(pg.BurnedCards, pg.Deck.Cards[0]) // Burn one
	pg.Deck.Cards = pg.Deck.Cards[1:]

	pg.CommunityCards = append(pg.CommunityCards, pg.Deck.Cards[0]) // River
	pg.Deck.Cards = pg.Deck.Cards[1:]
}

func (pg PokerGame) String() string {
	// TODO: Update string methods to use string buffers instead of inefficient
	// string concatenation

	s := ""
	for _, p := range pg.Players {
		s += fmt.Sprintf("%s\n", p)
	}
	s += "\n"

	s += fmt.Sprintf("Community Cards: ")
	for i, c := range pg.CommunityCards {
		if i != len(pg.CommunityCards) - 1 {
			s += fmt.Sprintf("%s  ", c)
		} else {
			s += fmt.Sprintf("%s", c)
		}
	}
	s += fmt.Sprintf("\n\n")

	s += fmt.Sprintf("Burned Cards: ")
	for i, c := range pg.BurnedCards {
		if i != len(pg.BurnedCards) - 1 {
			s += fmt.Sprintf("%s  ", c)
		} else {
			s += fmt.Sprintf("%s", c)
		}
	}
	s += fmt.Sprintf("\n\n")

	s += fmt.Sprintf("Deck:\n")
	s += fmt.Sprintf("%s\n", pg.Deck)

	return s
}

// Analyze a poker game for player and community cards
func HasFlush(p Player, cc []Card) (string, bool) {
	cards := make([]Card, len(p.Cards) + len(cc))
	copy(cards, p.Cards)
	copy(cards[len(p.Cards):], cc)

	suitCount := make(map[Suit]int)
	for _, card := range cards {
		suitCount[card.Suit]++
	}

	for suit, count := range suitCount {
		if count >= 5 {
			return string(suit), true
		}
	}
	return "", false
}

func HasPair(p Player, cc []Card) (string, bool) {
	cards := make([]Card, len(p.Cards) + len(cc))
	copy(cards, p.Cards)
	copy(cards[len(p.Cards):], cc)

	valueCount := make(map[int]int)
	for _, card := range cards {
		valueCount[card.Value]++
	}

	for value, count := range valueCount {
		if count == 2 {
			return CardValueToString(value), true
		}
	}
	return "", false
}

func main() {
	var flushHitCount int
	var pairHitCount int
	n := 100000

	pokerGame := PokerGame{}
	for i := 0; i < n; i++ {
		pokerGame.Init(2)
		// fmt.Println(pokerGame)

		for _, player := range pokerGame.Players {
			_, hasFlush := HasFlush(player, pokerGame.CommunityCards)
			if hasFlush {
				flushHitCount++
				break
			}
		}

		for _, player := range pokerGame.Players {
			_, hasPair := HasPair(player, pokerGame.CommunityCards)
			if hasPair {
				pairHitCount++
				break
			}
		}
	}

	fmt.Println(float64(flushHitCount) / float64(n))
	fmt.Println(float64(pairHitCount) / float64(n))
}