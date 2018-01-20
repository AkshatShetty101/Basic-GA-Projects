package main

import (
	"fmt"
)

func main() {
	var colors map[string]string
	col := make(map[string]string)
	col["black"] = "#000000"
	colours := map[string]string{
		"red":   "#ff0000",
		"green": "#00ff00",
		"blue":  "#0000ff",
		"white": "#ffffff",
	}
	fmt.Println(colours, colors, col)
	delete(colours, "red")
	fmt.Println(colours, colors, col)
	printMap(colours)
}

func printMap(m map[string]string) {
	for col, hex := range m {
		fmt.Println(col, hex)
	}
}
