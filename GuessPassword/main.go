package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func main() {
	geneset := " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!."
	sliceGeneset := strings.Split(geneset, "")
	// target := "Hello World!"
}

func generateParent(geneset []string, length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	output := ""
	fmt.Println(len(geneset))
	r := rand.New(source)
	for i := 0; i < length; i++ {
		newPosition := r.Intn(len(geneset) - 1)
		output += geneset[newPosition]
	}
	fmt.Println(output)
	return output
}

func getFitness() int {

}
