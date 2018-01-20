package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()
	var sliceGeneset []int
	for i := 0; i < 100; i++ {
		sliceGeneset = append(sliceGeneset, i)
	}
	bestParent := generateParent(sliceGeneset, 12)
	bestFitness := getFitness(bestParent)
	display(bestParent, bestFitness, start)
	for {
		child := mutate(bestParent, sliceGeneset)
		childFitness := getFitness(child)
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

func generateParent(geneset []int, length int) []int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	var output []int
	for i := 0; i < length; i++ {
		newPosition := r.Intn(len(geneset) - 1)
		output = append(output, geneset[newPosition])
	}
	return output
}

func getFitness(test []int) int {
	match := 1
	for i := 1; i < len(test); i++ {
		if test[i] > test[i-1] {
			match++
		}
	}
	return match
}

func mutate(parent []int, geneset []int) []int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	location := r.Intn(len(parent))
	element := parent[location]
	x := element
	for {
		x = geneset[r.Intn(len(geneset)-1)]
		if x != element {
			break
		}
	}
	parent[location] = x
	return parent
}

func display(guess []int, fitness int, start time.Time) {
	t := time.Now()
	fmt.Printf("%v\t%v\t%v\n", guess, fitness, t.Sub(start))
}
