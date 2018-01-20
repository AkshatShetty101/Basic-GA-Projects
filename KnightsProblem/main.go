package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Position is the struct to give the row and column of the knight
type Position struct {
	x int
	y int
}

func main() {
	start := time.Now()
	xMax := 8
	yMax := 8
	var board [8][8]string
	nKnights := 14
	knightsLocation := make([]Position, nKnights)
	bestParent := setStartLocation(knightsLocation, xMax, yMax)
	att := calcKnightAttack(bestParent, xMax, yMax)
	bestFitness := getFitness(att)
	display(&board, bestParent, att, bestFitness, start)
	for {
		tmp := make([]Position, len(bestParent))
		copy(tmp, bestParent)
		child := mutate(tmp, xMax, yMax)
		attChild := calcKnightAttack(child, xMax, yMax)
		childFitness := getFitness(attChild)
		if childFitness < bestFitness {
			continue
		}
		display(&board, child, attChild, childFitness, start)
		if childFitness == xMax*yMax {
			break
		}
		bestParent = child
		bestFitness = childFitness
	}
	fmt.Println("Finished!")
}

func setStartLocation(q []Position, xMax int, yMax int) []Position {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	for i := 0; i < len(q); i++ {
		for {
			flag := -1
			newRow := r.Intn(xMax)
			newCol := r.Intn(yMax)
			for j := 0; j < len(q); j++ {
				if newRow == q[j].x && newCol == q[j].y {
					flag = 0
					break
				}
			}
			if flag == -1 {
				q[i].x = newRow
				q[i].y = newCol
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

func calcKnightAttack(k []Position, xMax int, yMax int) [][]Position {
	posX := []int{2, 1}
	posY := []int{1, 2}
	att := make([][]Position, len(k))
	for i, val := range k {
		for j, pos := range posX {
			if (val.x+pos >= 0 && val.x+pos < xMax) && (val.y+posY[j] >= 0 && val.y+posY[j] < xMax) {
				att[i] = append(att[i], Position{x: val.x + pos, y: val.y + posY[j]})
			}
			if (val.x-pos >= 0 && val.x-pos < xMax) && (val.y+posY[j] >= 0 && val.y+posY[j] < xMax) {
				att[i] = append(att[i], Position{x: val.x - pos, y: val.y + posY[j]})
			}
			if (val.x+pos >= 0 && val.x+pos < xMax) && (val.y-posY[j] >= 0 && val.y-posY[j] < xMax) {
				att[i] = append(att[i], Position{x: val.x + pos, y: val.y - posY[j]})
			}
			if (val.x-pos >= 0 && val.x-pos < xMax) && (val.y-posY[j] >= 0 && val.y-posY[j] < xMax) {
				att[i] = append(att[i], Position{x: val.x - pos, y: val.y - posY[j]})
			}
		}
	}
	return att
}

func getFitness(att [][]Position) int {
	attMap := make(map[string]bool)
	st := ""
	for _, val := range att {
		for _, pos := range val {
			st = strconv.Itoa(pos.x) + "_" + strconv.Itoa(pos.y)
			attMap[st] = true
		}
	}
	return len(attMap)
}

func updateBoard(b *[8][8]string, k []Position, att [][]Position) {
	fmt.Println(k)
	for i := 0; i < len(k); i++ {
		b[k[i].x][k[i].y] = "K(" + strconv.Itoa(i) + ")"
	}
	for i, val := range att {
		for _, pos := range val {
			if b[pos.x][pos.y] == "-" {
				b[pos.x][pos.y] = strconv.Itoa(i)
			} else {
				b[pos.x][pos.y] += "," + strconv.Itoa(i)
			}
		}
	}
}

func display(b *[8][8]string, k []Position, att [][]Position, fitness int, start time.Time) {
	t := time.Now()
	clearBoard(b)
	updateBoard(b, k, att)
	fmt.Println("\nBoard:-")
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b[i]); j++ {
			fmt.Print(b[i][j], "\t\t|")
		}
		fmt.Println()
	}
	fmt.Printf("Fitness :%v\tTime :%v", fitness, t.Sub(start))
	fmt.Println("\n-------------------------------------------------")
}

func mutate(k []Position, xMax int, yMax int) []Position {
	source := rand.NewSource(time.Now().UnixNano())
	source1 := rand.NewSource(time.Now().UnixNano() + time.Now().UnixNano())
	r := rand.New(source)
	r1 := rand.New(source1)
	for j := 0; j < 1; j++ {
		number := r.Intn(len(k))
		for {
			flag := -1
			newRow := r1.Intn(xMax)
			newCol := r1.Intn(yMax)
			if k[number].x == newRow && k[number].y == newCol {
				continue
			}
			for i := 0; i < len(k); i++ {
				if newRow == k[i].x && newCol == k[i].y {
					flag = 0
					break
				}
			}
			if flag == -1 {
				k[number].x = newRow
				k[number].y = newCol
				break
			}
		}
	}
	return k
}
