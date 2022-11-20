package main

import (
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) {
	expected := 52
	d := newDeck()

	if len(d) != 52 {
		t.Errorf("Expected deck length of %d, but got %v", expected, len(d))
	}

	firstCard := d[0]
	if firstCard != "Ace of Spades" {
		t.Errorf("Expected 'Ace of Spades', but got %s", firstCard)
	}

	lastCard := d[len(d)-1]
	if lastCard != "King of Clubs" {
		t.Errorf("Expected 'King of Clubs', but got %s", lastCard)
	}
}

func TestSaveToDeckAndNNewDeckFromFile(t *testing.T) {
	expected := 52
	filename := "_decktesting.testing"

	// Clean the 'decktesting' file before start of the test
	os.Remove(filename)

	d := newDeck()
	d.saveToFile(filename)

	dLoad := newDeckFromFile(filename)

	if len(dLoad) != expected {
		t.Errorf("Expected deck length of %d, but got %v", expected, len(dLoad))
	}

	os.Remove(filename)
}
