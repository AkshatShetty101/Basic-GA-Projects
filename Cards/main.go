package main

import (
	"fmt"
)

func main() {
	// Creating a new Deck
	cards := newDeck()
	fmt.Println("Newly created deck:-\n")
	cards.print()
	fmt.Println("-------------------------------------")
	// Creating Hands
	fmt.Println("Creating hands from the deck:-\n")
	hand, remainingDeck := deal(cards, 5)
	hand.print()
	fmt.Println("-------------------------------------")
	remainingDeck.print()
	fmt.Println("-------------------------------------")
	// Saving deck to File
	fmt.Println("Saving current deck to file:-\n")
	cards.saveToFile("myCards")
	// Reading deck from file
	fmt.Println("-------------------------------------")
	fmt.Println("Loading deck from file:-\n")
	d := newDeckFromFile("myCards")
	fmt.Println(d)
	fmt.Println("-------------------------------------")
	// Shuffling the cards
	fmt.Println("Shuffling the deck:-\n")
	cards.shuffle()
	cards.print()
}
