package main

import (
	"fmt"
	"math/rand"
	"time"
)

type gene struct {
	row    [8]int
	column [8]int
}

func main() {
	start := time.Now()
	var board [8][8]string
	var qLocation gene
	bestParent := setStartLocation(qLocation)
	bestFitness := getFitness(bestParent)
	display(&board, bestParent, bestFitness, start)
	// count := 0
	for {
		// fmt.Print("!")
		child := mutate(bestParent)
		childFitness := getFitness(child)
		if compareFitness(child, bestParent) == false {
			continue
		}
		display(&board, child, childFitness, start)
		if childFitness == 0 {
			break
		}
		bestFitness = childFitness
		bestParent = child

	}
	fmt.Println("Finished!")

}

func setStartLocation(q gene) gene {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	for i := 0; i < 8; i++ {
		for {
			flag := -1
			newRow := r.Intn(len(q.row))
			newCol := r.Intn(len(q.column))
			for j := 0; j < len(q.row); j++ {
				if newRow == q.row[j] && newCol == q.column[j] {
					flag = 0
					break
				}
			}
			if flag == -1 {
				q.row[i] = newRow
				q.column[i] = newCol
				break
			}

		}
	}
	return q
}

func clearBoard(b *[8][8]string) {
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b[i]); j++ {
			b[i][j] = "-"
		}
	}
}
func updateBoard(b *[8][8]string, q gene) {
	for i := 0; i < len(q.row); i++ {
		b[q.row[i]][q.column[i]] = "Q"
	}
}

func display(b *[8][8]string, q gene, fitness int, start time.Time) {
	t := time.Now()
	clearBoard(b)
	updateBoard(b, q)
	fmt.Println("\nBoard:-")
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b[i]); j++ {
			fmt.Print("\t", b[i][j])
		}
		fmt.Println()
	}
	fmt.Printf("Fitness :%v\tTime :%v", fitness, t.Sub(start))
	fmt.Println("\n-------------------------------------------------")
}

func compareFitness(test1 gene, test2 gene) bool {
	f1 := getFitness(test1)
	f2 := getFitness(test2)
	// fmt.Println(f1, f2)
	// fmt.Println(f1 < f2)

	return f1 <= f2
}

func getFitness(q gene) int {
	var row = make(map[int]bool)
	var col = make(map[int]bool)
	var seDiag = make(map[int]bool)
	var neDiag = make(map[int]bool)
	for i := 0; i < len(q.row); i++ {
		row[q.row[i]] = true
		col[q.column[i]] = true
		seDiag[len(q.row)-1-q.row[i]+q.column[i]] = true
		neDiag[q.row[i]+q.column[i]] = true
	}
	total := 4*len(q.row) - len(row) - len(col) - len(neDiag) - len(seDiag)
	return total
}

func mutate(q gene) gene {
	source := rand.NewSource(time.Now().UnixNano())
	source1 := rand.NewSource(time.Now().UnixNano() + time.Now().UnixNano())
	r := rand.New(source)
	r1 := rand.New(source1)
	qRowNumber := r.Intn(len(q.row))
	qColNumber := r.Intn(len(q.column))
	for {
		flag := -1
		newRow := r1.Intn(len(q.row))
		newCol := r1.Intn(len(q.column))
		if q.row[qRowNumber] == newRow && q.column[qColNumber] == newCol {
			continue
		}
		for i := 0; i < len(q.row); i++ {
			if newRow == q.row[i] && newCol == q.column[i] {
				flag = 0
				break
			}
		}
		if flag == -1 {
			q.row[qRowNumber] = newRow
			q.column[qColNumber] = newCol
			break
		}
	}
	return q
}
