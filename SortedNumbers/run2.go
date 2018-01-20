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
	// Optimal defination
	optimal := sliceGeneset[:10]
	fmt.Println(optimal)
	fmt.Println(getFitness2(optimal))
	//Generating starting list of numbers
	bestParent := generateParent2(sliceGeneset, 10)
	bestFitness, bestGap := getFitness2(bestParent)
	display2(bestParent, bestFitness, bestGap, start)
	for {
		tmp := make([]int, len(bestParent))
		copy(tmp, bestParent)
		child := mutate2(tmp, sliceGeneset)
		childFitness, childGap := getFitness2(child)
		if compareFitness2(child, bestParent) == false {
			continue
		}
		display2(child, childFitness, childGap, start)
		if compareFitness2(optimal, child) == false {
			break
		}
		bestGap = childGap
		bestFitness = childFitness
		bestParent = child

	}
	fmt.Println("Finished!")
}

func generateParent2(geneset []int, length int) []int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	var output []int
	for i := 0; i < length; i++ {
		newPosition := r.Intn(len(geneset) - 1)
		output = append(output, geneset[newPosition])
	}
	return output
}

func getFitness2(test []int) (int, int) {
	match := 1
	gap := 0
	for i := 1; i < len(test); i++ {
		if test[i] > test[i-1] {
			match++
		} else {
			gap += test[i-1] - test[i]
		}
	}
	return match, gap
}

func compareFitness2(test1 []int, test2 []int) bool {
	f1, g1 := getFitness2(test1)
	f2, g2 := getFitness2(test2)
	if f1 != f2 {
		return f1 > f2
	}
	return g1 < g2
}

func mutate2(parent []int, geneset []int) []int {

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

func display2(guess []int, fitness int, gap int, start time.Time) {
	t := time.Now()
	fmt.Printf("Guess :%v\t Fitness :%v\tGap :%v\tTime :%v\n", guess, fitness, gap, t.Sub(start))
}
