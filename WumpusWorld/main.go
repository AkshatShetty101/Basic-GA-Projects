package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

type coords struct {
	x int
	y int
}

type History struct {
	x      int
	y      int
	facing string
	status string
}

type gene struct {
	x       int
	y       int
	facing  string
	victory bool
	fire    bool
	score   int
	moves   int
}

func main() {
	var board [4][4]string
	var playerLocation gene
	var history []History
	var pitIssues = make(map[coords]string)
	var wumpusIssues = make(map[coords]string)
	var goldLocation = gene{x: 3, y: 2, facing: "default", victory: true, fire: false, score: 0, moves: 0}
	makeBoard(&board)
	fmt.Println("started")
	bestParent := setStartLocation(playerLocation)
	bestFitness := getFitness(bestParent, goldLocation, &board)
	display(bestParent, goldLocation, &board, bestFitness, &history)
	for i := 0; i < 10; i++ {
		child := move(bestParent, &history, &board, &pitIssues, &wumpusIssues)
		if child.victory == true {
			fmt.Println("Found Gold!!!")
			break
		} else {
			childFitness := getFitness(child, goldLocation, &board)
			display(child, goldLocation, &board, childFitness, &history)
			// if compareFitness(childFitness, bestFitness) {
			bestParent = child
			bestFitness = childFitness
			//display(bestParent, goldLocation, &board, bestFitness, &history)
			// }
		}
	}
}

func display(player gene, gold gene, board *[4][4]string, fitness float64, history *[]History) {
	fmt.Println("------------------------------------------------")
	fmt.Printf("Location: x-%v\t y-%v\t Direction:%v \tFitness:%v\nBoard Status: %s\nMoves:%v\nScore:%v\nGold Location: x-%v\t y-%v\n",
		player.x,
		player.y,
		player.facing,
		fitness,
		board[player.x][player.y],
		player.moves,
		player.score,
		gold.x,
		gold.y)
	//fmt.Println(history)
}
func makeBoard(board *[4][4]string) {

	// for i := 0; i < 4; i++ {
	// 	for j := 0; j < 4; j++ {
	// 		board[i][j] = ""
	// 	}
	// }
	board[1][0] = "breeze"
	board[2][0] = "pit"
	board[3][0] = "breeze"
	board[0][1] = "stench"
	board[2][1] = "breeze"
	board[0][2] = "wumpus"
	board[1][2] = "breeze,stench"
	board[2][2] = "pit"
	board[3][2] = "breeze"
	board[0][3] = "stench"
	board[2][3] = "breeze,gold"
	board[3][3] = "pit"
}

func setStartLocation(q gene) gene {
	q.x = 0
	q.y = 0
	q.facing = "E"
	q.victory = false
	q.fire = false
	q.score = 0
	q.moves = 0
	return q
}

func getFitness(player gene, gold gene, board *[4][4]string) float64 {
	boardStatus := strings.Split(board[player.x][player.y], ",")
	for _, status := range boardStatus {
		if status == "pit" || status == "wumpus" {
			return -1.0
		}
	}
	dist := math.Sqrt(math.Pow(float64(player.x-gold.x), 2) +
		math.Pow(float64(player.y-gold.y), 2))
	return float64(player.moves)/100 - dist/100
}

func compareFitness(f1 float64, f2 float64) bool {
	return f1 >= f2
}

func moveForward(player gene) gene {
	if player.facing == "W" {
		player.x = player.x - 1
	} else if player.facing == "E" {
		player.x = player.x + 1
	} else if player.facing == "N" {
		player.y = player.y + 1
	} else {
		player.y = player.y - 1
	}
	player.moves = player.moves + 1
	return player
}

func move(player gene, history *[]History, board *[4][4]string, pitIssues *map[coords]string, wumpusIssues *map[coords]string) gene {
	// s1 := rand.NewSource(time.Now().UnixNano())
	// r1 := rand.New(s1)
	prev := History{x: player.x, y: player.y, facing: player.facing, status: board[player.x][player.y]}
	*history = append(*history, prev)
	//fmt.Println(history)
	possibleLeft := checkTurnLeft(player, board)
	possibleRight := checkTurnRight(player, board)
	boardStatus := strings.Split(board[player.x][player.y], ",")
	//fmt.Println(boardStatus)
	//fmt.Println(len(boardStatus))
	if len(boardStatus[0]) > 0 {
		//fmt.Println("has status")
		for _, status := range boardStatus {
			if status == "gold" {
				player.victory = true
				return player
			}
		}
		for _, status := range boardStatus {
			if status == "breeze" {
				if next := moveForward(player); (*pitIssues)[coords{x: next.x, y: next.y}] == "-" {
					return next
				}
				addPoW(player, history, pitIssues, "pit")
				checkDiagonal(history, player, pitIssues, "breeze", "pit")
				fmt.Println("Backtracking!")
				stepBack := backtrack(player, history)
				return lOrR(checkTurnLeft(stepBack, board), checkTurnRight(stepBack, board), stepBack, board)
			}
			if status == "stench" {
				if next := moveForward(player); (*wumpusIssues)[coords{x: next.x, y: next.y}] == "-" {
					return next
				}
				addPoW(player, history, wumpusIssues, "wumpus")
				checkDiagonal(history, player, wumpusIssues, "stench", "wumpus")
				return checkAndFire(player, wumpusIssues, pitIssues, board)
			}
		}
	} else {
		move := moveForward(player)
		if checkMoveValidity(move) == false {
			return lOrR(possibleLeft, possibleRight, player, board)
		}
		return move
	}
	return lOrR(possibleLeft, possibleRight, player, board)
}

func checkAndFire(player gene, wumpusIssues *map[coords]string, pitIssues *map[coords]string, board *[4][4]string) gene {
	for key, wumpus := range *wumpusIssues {
		if wumpus == "wumpus" {
			if (*pitIssues)[key] != "pit" {
				if ok, direction := checkProximity(player, key); ok {
					if player.facing == direction {
						player.moves++
						if fire(player, wumpusIssues, board) {
							player.score++
						}

					} else {
						// player = align(player, direction)
					}
				}
			}

		}
	}
	return player
}

// func align(player gene, direction string gene) {

// 	return player
// }

func fire(player gene, issues *map[coords]string, board *[4][4]string) bool {
	var flag = false
	var index = -1
	for {
		if player = moveForward(player); (player.x > -1 && player.x < 4) && (player.y > -1 && player.y < 4) {
			status := strings.Split(board[player.x][player.y], ",")
			for i, val := range status {
				if val == "wumpus" {
					index = i
					flag = true
				}
				if flag == true {
					board[player.x][player.y] = strings.Join(append(status[:index], status[index+1:]...), ",")
				}
			}
			checkAndDelete(coords{x: player.x, y: player.y}, issues)
		} else {
			break
		}
	}
	return flag
}

func checkProximity(player gene, point coords) (bool, string) {

	if player.x-1 == point.x && player.y == point.y {
		return true, "E"
	}
	if player.x+1 == point.x && player.y == point.y {
		return true, "W"
	}
	if player.x == point.x && player.y-1 == point.y {
		return true, "S"
	}
	if player.x == point.x && player.y+1 == point.y {
		return true, "N"
	}
	return false, ""
}

func checkDiagonal(history *[]History, player gene, issues *map[coords]string, value1 string, value2 string) {
	from := gene{x: (*history)[len(*history)-2].x, y: (*history)[len(*history)-2].x}
	for _, value := range *history {
		if value.status == value1 {
			diags := getDiags(coords{x: value.x, y: value.y})
			for _, val := range diags {
				if checkHistoryForStatus(val, history) == value1 {
					switch val {
					case coords{x: value.x - 1, y: value.y - 1}:
						if value.x-1 != from.x && value.y != from.y {
							(*issues)[coords{x: value.x - 1, y: value.y}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						if value.x+1 != from.x && value.y-1 != from.y {
							(*issues)[coords{x: value.x + 1, y: value.y - 1}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						break
					case coords{x: value.x + 1, y: value.y - 1}:
						if value.x != from.x && value.y-1 != from.y {
							(*issues)[coords{x: value.x, y: value.y - 1}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						if value.x+1 != from.x && value.y != from.y {
							(*issues)[coords{x: value.x + 1, y: value.y}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						break
					case coords{x: value.x + 1, y: value.y + 1}:
						if value.x != from.x && value.y+1 != from.y {
							(*issues)[coords{x: value.x, y: value.y + 1}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						if value.x+1 != from.x && value.y != from.y {
							(*issues)[coords{x: value.x + 1, y: value.y}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						break
					case coords{x: value.x - 1, y: value.y + 1}:
						if value.x-1 != from.x && value.y != from.y {
							(*issues)[coords{x: value.x - 1, y: value.y}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						if value.x != from.x && value.y+1 != from.y {
							(*issues)[coords{x: value.x, y: value.y + 1}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						break
					}
				}
			}
		}
	}
}

func checkAndDelete(point coords, issues *map[coords]string) {
	_, ok := (*issues)[point]
	if ok {
		(*issues)[point] = "-"
	}
}

func checkHistoryForStatus(point coords, history *[]History) string {
	for _, value := range *history {
		if value.x == point.x && value.y == point.y {
			return value.status
		}
	}
	return "-"
}

func getDiags(point coords) []coords {
	var diag []coords
	if (point.x-1 > -1 && point.x-1 < 4) && (point.y-1 > -1 && point.y-1 < 4) {
		diag = append(diag, coords{x: point.x - 1, y: point.y - 1})
	}
	if (point.x-1 > -1 && point.x-1 < 4) && (point.y+1 > -1 && point.y+1 < 4) {
		diag = append(diag, coords{x: point.x - 1, y: point.y + 1})
	}
	if (point.x+1 > -1 && point.x+1 < 4) && (point.y-1 > -1 && point.y-1 < 4) {
		diag = append(diag, coords{x: point.x + 1, y: point.y - 1})
	}
	if (point.x+1 > -1 && point.x+1 < 4) && (point.y+1 > -1 && point.y+1 < 4) {
		diag = append(diag, coords{x: point.x + 1, y: point.y + 1})
	}
	return diag
}

func addPoW(player gene, history *[]History, issues *map[coords]string, value string) {
	from := gene{x: (*history)[len(*history)-2].x, y: (*history)[len(*history)-2].x}
	fmt.Println("From::")
	fmt.Println(from)
	if (player.x-1 > -1 && player.x-1 < 4) && (player.x-1 == from.x && player.y == from.y) == false {
		fmt.Println(coords{x: player.x - 1, y: player.y})
		(*issues)[coords{x: player.x - 1, y: player.y}] = value
	}
	if (player.x+1 > -1 && player.x+1 < 4) && (player.x+1 == from.x && player.y == from.y) == false {
		(*issues)[coords{x: player.x + 1, y: player.y}] = value
	}
	if (player.y-1 > -1 && player.y-1 < 4) && (player.x == from.x && player.y-1 == from.y) == false {
		(*issues)[coords{x: player.x, y: player.y - 1}] = value
	}
	if (player.y+1 > -1 && player.y+1 < 4) && (player.x == from.x && player.y+1 == from.y) == false {
		(*issues)[coords{x: player.x, y: player.y + 1}] = value
	}
}

func checkStagnant(history *[]History) bool {
	max := len(*history)
	if (*history)[max-1] == (*history)[max-2] {
		return true
	}
	return false
}

func backtrack(player gene, history *[]History) gene {
	max := len(*history)
	player.x = (*history)[max-2].x
	player.y = (*history)[max-2].y
	player.facing = (*history)[max-2].facing
	player.moves = player.moves + 1
	return player
}

func lOrR(possibleLeft bool, possibleRight bool, player gene, board *[4][4]string) gene {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	if possibleLeft == true {
		if possibleRight == true {
			if r1.Intn(1) == 0 {
				return turnLeft(player, board)
			}
			return turnRight(player, board)
		}
		return turnLeft(player, board)
	}
	return turnRight(player, board)
}

func checkMoveValidity(player gene) bool {
	if player.x > -1 && player.x < 4 && player.y > -1 && player.y < 4 {
		return true
	}
	return false
}

func turnLeft(player gene, board *[4][4]string) gene {
	if player.facing == "N" {
		player.facing = "W"
		// player.x = player.x - 1
	} else if player.facing == "S" {
		player.facing = "E"
		// player.x = player.x + 1
	} else if player.facing == "E" {
		player.facing = "N"
		// player.y = player.y + 1
	} else {
		// player.y = player.y - 1
		player.facing = "S"
	}
	player.fire = false
	player.moves = player.moves + 1
	return player
}

func turnRight(player gene, board *[4][4]string) gene {
	if player.facing == "S" {
		player.facing = "W"
		// player.x = player.x - 1
	} else if player.facing == "N" {
		player.facing = "E"
		// player.x = player.x + 1
	} else if player.facing == "W" {
		player.facing = "N"
		// player.y = player.y + 1
	} else {
		// player.y = player.y - 1
		player.facing = "S"
	}
	player.fire = false
	player.moves = player.moves + 1
	return player
}

func checkTurnLeft(player gene, board *[4][4]string) bool {
	if player.facing == "N" {
		if player.x-1 > -1 {
			return true
		}
	} else if player.facing == "S" {
		if player.x+1 < 4 {
			return true
		}
	} else if player.facing == "E" {
		if player.y+1 < 4 {
			return true
		}
	} else {
		if player.y-1 > -1 {
			return true
		}
	}
	return false
}

func checkTurnRight(player gene, board *[4][4]string) bool {
	if player.facing == "N" {
		if player.x+1 < 4 {
			return true
		}
	} else if player.facing == "S" {
		if player.x-1 > -1 {
			return true
		}
	} else if player.facing == "W" {
		if player.y+1 < 4 {
			return true
		}
	} else {
		if player.y-1 > -1 {
			return true
		}
	}
	return false
}
