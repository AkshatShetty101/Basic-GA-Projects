package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()
	n := 0
	fmt.Print("Enter the number whose magic square is required: ")
	fmt.Scanln(&n)
	geneSet := make([]int, n*n)
	bestParent := make([][]int, n, n)
	for i := 0; i < len(geneSet); i++ {
		geneSet[i] = i + 1
		bestParent[int(math.Floor(float64(i/n)))] = append(bestParent[int(math.Floor(float64(i/n)))], (i + 1))
	}
	expectedSum := n * (n*n + 1) / 2
	bestFitness := getFitness(bestParent, expectedSum)
	display(bestParent, bestFitness, start)
	for {
		child := mutate(bestParent)
		childFitness := getFitness(child, expectedSum)
		if childFitness < bestFitness {
			bestParent, bestFitness = child, childFitness
			display(child, childFitness, start)
		}
		if bestFitness == 0 {
			break
		}
	}
}

func getFitness(parent [][]int, expectedSum int) int {
	seDiag, neDiag, sum := 0, 0, 0
	for i := 0; i < len(parent); i++ {
		rowSum := 0
		colSum := 0
		for j := 0; j < len(parent); j++ {
			rowSum += parent[i][j]
			colSum += parent[j][i]
			if i == j {
				seDiag += parent[i][j]
				neDiag += parent[i][len(parent)-i-1]
			}

		}
		// fmt.Println("RowSum: ", rowSum)
		// fmt.Println("ColSum: ", colSum)
		sum += int(math.Abs(float64(rowSum-expectedSum)) + math.Abs(float64(colSum-expectedSum)))
	}
	sum += int(math.Abs(float64(neDiag-expectedSum)) + math.Abs(float64(seDiag-expectedSum)))
	return sum
}

func display(parent [][]int, fitness int, start time.Time) {
	t := time.Now()
	fmt.Println("\n-----------------------------------------------")
	fmt.Println("Fitness: ", fitness)
	fmt.Println("Time taken: ", t.Sub(start))

	for i := 0; i < len(parent); i++ {
		for j := 0; j < len(parent); j++ {
			fmt.Printf("%v\t", parent[i][j])
		}
		fmt.Println()
	}
}

func mutate(parent [][]int) [][]int {
	source := rand.NewSource(time.Now().UnixNano())
	source1 := rand.NewSource(time.Now().UnixNano() + time.Now().UnixNano())
	r := rand.New(source)
	r1 := rand.New(source1)
	for i := 0; i < 5; i++ {
		for {
			number := r.Intn(len(parent) * len(parent))
			number1 := r1.Intn(len(parent) * len(parent))
			if number != number1 {
				parent[int(math.Floor(float64(number/len(parent))))][number%len(parent)],
					parent[int(math.Floor(float64(number1/len(parent))))][number1%len(parent)] = parent[int(math.Floor(float64(number1/len(parent))))][number1%len(parent)],
					parent[int(math.Floor(float64(number/len(parent))))][number%len(parent)]
				// fmt.Println(parent, number, number1)
				break
			}
		}
	}
	return parent
}
