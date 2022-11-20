package main

func main() {
	filename := "/tmp/help"

	cards := newDeck()

	hand, rest := deal(cards, 8)
	hand.print()
	rest.print()

	cards.saveToFile(filename)

	newD := newDeckFromFile(filename)
	newD.shuffle()
	newD.print()
}
