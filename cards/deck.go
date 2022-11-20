package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Create a new type of 'deck', which is a slice of strings
type deck []string

func newDeck() deck {
	cards := deck{}

	cardSuits := []string{"Spades", "Diamonds", "Hearts", "Clubs"}
	cardValues := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}

	for _, suit := range cardSuits {
		for _, value := range cardValues {
			cards = append(cards, value+" of "+suit)
		}
	}

	return cards
}

func (d deck) print() {
	for _, card := range d {
		fmt.Println(card)
	}
}

func (d deck) toString() string {
	return strings.Join([]string(d), ",")
}

func deal(d deck, handSize int) (deck, deck) {
	return d[:handSize], d[handSize:]
}

func (d deck) saveToFile(filename string) error {
	return os.WriteFile(filename, []byte(d.toString()), 0744)
}

func newDeckFromFile(filename string) deck {
	bsd, err := os.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	cards := strings.Split(string(bsd), ",")

	return deck(cards)
}

func (d deck) shuffle() {
	timeNow := time.Now().UnixNano()
	r := rand.New(rand.NewSource(timeNow))

	size := len(d)

	for i := range d {
		newPosition := r.Intn(size - 1)

		d[i], d[newPosition] = d[newPosition], d[i]
	}
}
