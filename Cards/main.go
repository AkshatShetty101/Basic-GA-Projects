package main

import (
	"fmt"
)

func main() {
	cards := newDeck()
	cards.print()
	hand, remainingDeck := deal(cards, 5)
	fmt.Println("-------------------------------------")
	hand.print()
	fmt.Println("-------------------------------------")
	remainingDeck.print()
	x := cards.toString()
	fmt.Println(x)
	cards.saveToFile("myCards.txt")
}
