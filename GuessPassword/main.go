package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	geneset := " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!."
	sliceGeneset := strings.Split(geneset, "")
	target := "Hello World!"
	fmt.Println("Password : ", target)
	bestParent := generateParent(sliceGeneset, 12)
	bestFitness := getFitness(bestParent, target)
	display(bestParent, bestFitness, start)
	for {
		child := mutate(bestParent, sliceGeneset)
		childFitness := getFitness(child, target)
		if bestFitness >= childFitness {
			continue
		}
		display(child, childFitness, start)
		bestFitness = childFitness
		bestParent = child
		if childFitness == len(bestParent) {
			break
		}
	}
	fmt.Println("Finished!")
}

func generateParent(geneset []string, length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	output := ""
	for i := 0; i < length; i++ {
		newPosition := r.Intn(len(geneset) - 1)
		output += geneset[newPosition]
	}
	return output
}

func getFitness(test string, target string) int {
	match := 0
	if len(test) == len(target) {
		for i := 0; i < len(test); i++ {
			if test[i] == target[i] {
				match++
			}
		}
	} else {
		fmt.Println("The two string lengths do not match")
	}
	return match
}

func mutate(parent string, geneset []string) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	parentSlice := strings.Split(string(parent), "")
	location := r.Intn(len(parent))
	element := parentSlice[location]
	x := element
	for {
		x = geneset[r.Intn(len(geneset)-1)]
		if x != element {
			break
		}
	}
	parentSlice[location] = x
	return strings.Join(parentSlice, "")
}

func display(guess string, fitness int, start time.Time) {
	t := time.Now()
	fmt.Printf("%v\t%v\t%v\n", guess, fitness, t.Sub(start))
}
