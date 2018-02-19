package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jinzhu/copier"
)

// Chromosome represents each
type Chromosome struct {
	genes   [][]int
	fitness int
	age     int
}

func initialize(n int, expectedSum int) Chromosome {
	parent := Chromosome{}
	parent.genes = make([][]int, n, n)
	list := rand.Perm(n * n)
	for i := 0; i < (n * n); i++ {
		parent.genes[int(math.Floor(float64(i/n)))] = append(parent.genes[int(math.Floor(float64(i/n)))], list[i]+1)
	}
	parent.fitness = getFitness(parent, expectedSum)
	parent.age = 0
	return parent
}

func main() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	var historicalFitness []int
	maxAge := 500
	start := time.Now()
	n := 0
	fmt.Println(maxAge)
	fmt.Print("Enter the number whose magic square is required: ")
	fmt.Scanln(&n)
	expectedSum := n * (n*n + 1) / 2
	parent := initialize(n, expectedSum)
	bestParent := Chromosome{}
	copier.Copy(&bestParent, &parent)
	historicalFitness = append(historicalFitness, bestParent.fitness)
	display(bestParent, start)
	for {
		// fmt.Println(historicalFitness)
		child := mutate(parent)
		child.fitness = getFitness(child, expectedSum)
		if maxAge != -1 {
			parent.age++
			if maxAge < parent.age {
				index := locateIndex(historicalFitness, child.fitness)
				// fmt.Println("index: ", index)
				diff := len(historicalFitness) - index
				propotionalSimilar := math.Exp(float64(-1 * diff))
				// fmt.Println("proportional Similar: ", propotionalSimilar)
				x := r.Float64()
				// fmt.Println(x)
				if x < propotionalSimilar {
					// fmt.Println("Copy child to parent")
					copier.Copy(&parent, &child)
				} else {
					// fmt.Println("Copy bestparent to parent")
					copier.Copy(&parent, &bestParent)
					parent.age = 0
				}
			}
			if child.fitness < parent.fitness {
				// fmt.Println("Copy child to parent")
				copier.Copy(&parent, &child)
				parent.age = 0

			} else {
				// fmt.Println("Update child age")
				child.age = parent.age + 1
				copier.Copy(&parent, &child)
			}

		}
		if child.fitness < bestParent.fitness {
			// fmt.Println("Copy child to bestparent")
			copier.Copy(&bestParent, &child)
			historicalFitness = append(historicalFitness, child.fitness)
			display(child, start)
		}
		if bestParent.fitness == 0 {
			break
		}
	}
}

func locateIndex(h []int, f int) int {
	if len(h) == 0 {
		return -1
	}
	// fmt.Println(h)
	// fmt.Println("Value to match: ", f)
	for i := len(h) - 1; i >= 0; i-- {
		if f < h[i] {
			// fmt.Println("Found: ", h[i])
			return i
		}
	}
	return 0
}

func getFitness(p Chromosome, expectedSum int) int {
	seDiag, neDiag, sum := 0, 0, 0
	for i := 0; i < len(p.genes); i++ {
		rowSum := 0
		colSum := 0
		for j := 0; j < len(p.genes); j++ {
			rowSum += p.genes[i][j]
			colSum += p.genes[j][i]
			if i == j {
				seDiag += p.genes[i][j]
				neDiag += p.genes[i][len(p.genes)-i-1]
			}

		}
		// fmt.Println("RowSum: ", rowSum)
		// fmt.Println("ColSum: ", colSum)
		sum += int(math.Abs(float64(rowSum-expectedSum)) + math.Abs(float64(colSum-expectedSum)))
	}
	sum += int(math.Abs(float64(neDiag-expectedSum)) + math.Abs(float64(seDiag-expectedSum)))
	return sum
}

func display(p Chromosome, start time.Time) {
	t := time.Now()
	fmt.Println("\n-----------------------------------------------")
	fmt.Println("Fitness: ", p.fitness)
	fmt.Println("Time taken: ", t.Sub(start))

	for i := 0; i < len(p.genes); i++ {
		for j := 0; j < len(p.genes[0]); j++ {
			fmt.Printf("%v\t", p.genes[i][j])
		}
		fmt.Println()
	}
}

func mutate(p Chromosome) Chromosome {
	source := rand.NewSource(time.Now().UnixNano())
	source1 := rand.NewSource(time.Now().UnixNano() + time.Now().UnixNano())
	r := rand.New(source)
	r1 := rand.New(source1)
	for i := 0; i < 5; i++ {
		for {
			number := r.Intn(len(p.genes) * len(p.genes))
			number1 := r1.Intn(len(p.genes) * len(p.genes))
			if number != number1 {
				p.genes[int(math.Floor(float64(number/len(p.genes))))][number%len(p.genes)],
					p.genes[int(math.Floor(float64(number1/len(p.genes))))][number1%len(p.genes)] = p.genes[int(math.Floor(float64(number1/len(p.genes))))][number1%len(p.genes)],
					p.genes[int(math.Floor(float64(number/len(p.genes))))][number%len(p.genes)]
				// fmt.Println(parent, number, number1)
				break
			}
		}
	}
	return p
}
