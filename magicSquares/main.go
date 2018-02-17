package main

import (
	"fmt"
	"math"
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
}

func getFitness(parent [][]int, expectedSum int) int {
	seDiag, neDiag, count := 0, 0, 0
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
		if rowSum == expectedSum {
			count++
		}
		if colSum == expectedSum {
			count++
		}
	}
	if neDiag == expectedSum {
		count++
	}
	if seDiag == expectedSum {
		count++
	}
	// fmt.Println("South East Diagonal: ", seDiag)
	// fmt.Println("North West Diagonal: ", neDiag)
	// fmt.Println("Exceeding: ", count)
	return count
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
