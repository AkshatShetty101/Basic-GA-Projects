package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()
	source := rand.NewSource(time.Now().UnixNano() + time.Now().UnixNano())
	r := rand.New(source)
	bestParent := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := range bestParent {
		j := r.Intn(i + 1)
		bestParent[i], bestParent[j] = bestParent[j], bestParent[i]
	}
	fmt.Println(bestParent)
	bestFitnessSum, bestFitnessProd := getFitness(bestParent)
	fmt.Println(bestFitnessSum, bestFitnessProd)
	display(bestParent, bestFitnessSum, bestFitnessProd, start)
	for {
		tmp := make([]int, len(bestParent))
		copy(tmp, bestParent)
		child := mutate(tmp)
		childFitnessSum, childFitnessProd := getFitness(child)
		if compareFitness(child, bestParent) == false {
			continue
		}
		display(child, childFitnessSum, childFitnessProd, start)
		if childFitnessSum == 36 && childFitnessProd == 360 {
			break
		}
		bestFitnessProd = childFitnessProd
		bestFitnessSum = childFitnessSum
		bestParent = child

	}
	fmt.Println("\nFinished!")

}

func display(s []int, fitnessSum int, fitnessProd int, start time.Time) {
	t := time.Now()
	fmt.Printf("\nSet :%v\t-\t%v\tSum :%v\tProduct :%v\tTime :%v", s[:5], s[5:], fitnessSum, fitnessProd, t.Sub(start))
}

func getFitness(p []int) (int, int) {
	sum := 0
	prod := 1
	for i, num := range p[:5] {
		sum += num
		prod *= p[i+5]
	}

	return sum, prod
}

func compareFitness(test1 []int, test2 []int) bool {
	fs1, fp1 := getFitness(test1)
	fs2, fp2 := getFitness(test2)
	totalDiff1 := math.Abs(float64(36-fs1)) + math.Abs(float64(360-fp1))
	totalDiff2 := math.Abs(float64(36-fs2)) + math.Abs(float64(360-fp2))
	return totalDiff1 <= totalDiff2
}

func mutate(s []int) []int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	for i := 0; i < 5; i++ {
		low := r.Intn(5)
		high := 5 + r.Intn(5)
		s[low], s[high] = s[high], s[low]
	}
	return s
}
