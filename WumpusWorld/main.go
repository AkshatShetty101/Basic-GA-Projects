package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type coords struct {
	x int
	y int
}

type History struct {
	c      coords
	facing string
	status string
}

type gene struct {
	c       coords
	facing  string
	victory bool
	fire    bool
	score   int
	moves   int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.Create("./output")
	check(err)
	var board [4][4]string
	var playerLocation gene
	var history []History
	var pitIssues = make(map[coords]string)
	var wumpusIssues = make(map[coords]string)
	var goldLocation = gene{c: coords{x: 2, y: 1}, facing: "default", victory: true, fire: false, score: 0, moves: 0}
	makeBoard(&board)
	bestParent := setStartLocation(playerLocation)
	bestFitness := getFitness(bestParent, goldLocation, &board)
	display(bestParent, goldLocation, &board, bestFitness, &history, &pitIssues, &wumpusIssues)
	for i := 0; i < 100; i++ {
		f.WriteString(strconv.Itoa(bestParent.c.y) +
			"-" + strconv.Itoa(bestParent.c.y) +
			"-" + string(bestParent.facing) +
			"-" + strconv.Itoa(bestParent.moves) +
			"-" + board[bestParent.c.x][bestParent.c.y] + "\n")
		child := move(bestParent, &history, &board, &pitIssues, &wumpusIssues)
		if child.victory == true {
			fmt.Println("Found Gold!!!")
			break
		} else {
			childFitness := getFitness(child, goldLocation, &board)
			f.WriteString(strconv.Itoa(child.c.x) +
				"-" + strconv.Itoa(child.c.y) +
				"-" + string(child.facing) +
				"-" + strconv.Itoa(child.moves) +
				"-" + board[child.c.x][child.c.y] + "\n")
			display(child, goldLocation, &board, childFitness, &history, &pitIssues, &wumpusIssues)
			bestParent = child
			bestFitness = childFitness
		}
	}
}

func display(player gene, gold gene, board *[4][4]string, fitness float64, history *[]History, pitIssues *map[coords]string, wumpusIssues *map[coords]string) {
	fmt.Println("------------------------------------------------")
	fmt.Printf("Location: x-%v\t y-%v\t Direction:%v \tFitness:%v\nBoard Status: %s\nMoves:%v\nScore:%v\nGold Location: x-%v\t y-%v\n",
		player.c.x,
		player.c.y,
		player.facing,
		fitness,
		board[player.c.x][player.c.y],
		player.moves,
		player.score,
		gold.c.x,
		gold.c.y)
	fmt.Println(history)
	fmt.Println(pitIssues)
	fmt.Println(wumpusIssues)
}
func makeBoard(board *[4][4]string) {

	board[1][0] = "breeze"
	board[2][0] = "pit"
	board[3][0] = "breeze"
	board[0][1] = "stench"
	board[2][1] = "breeze,gold"
	board[0][2] = "wumpus"
	board[1][2] = "breeze,stench"
	board[2][2] = "pit"
	board[3][2] = "breeze"
	board[0][3] = "stench"
	board[2][3] = "breeze"
	board[3][3] = "pit"
}

func setStartLocation(q gene) gene {
	q.c.x = 0
	q.c.y = 0
	q.facing = "E"
	q.victory = false
	q.fire = false
	q.score = 0
	q.moves = 0
	return q
}

func getFitness(player gene, gold gene, board *[4][4]string) float64 {
	boardStatus := strings.Split(board[player.c.x][player.c.y], ",")
	for _, status := range boardStatus {
		if status == "pit" || status == "wumpus" {
			return -1.0
		}
	}
	dist := math.Sqrt(math.Pow(float64(player.c.x-gold.c.x), 2) +
		math.Pow(float64(player.c.y-gold.c.y), 2))
	return float64(player.moves)/100 - dist/100
}

func compareFitness(f1 float64, f2 float64) bool {
	return f1 >= f2
}

func moveForward(player gene) gene {
	if player.facing == "W" {
		player.c.x = player.c.x - 1
	} else if player.facing == "E" {
		player.c.x = player.c.x + 1
	} else if player.facing == "N" {
		player.c.y = player.c.y + 1
	} else {
		player.c.y = player.c.y - 1
	}
	player.moves = player.moves + 1
	return player
}

func move(player gene, history *[]History, board *[4][4]string, pitIssues *map[coords]string, wumpusIssues *map[coords]string) gene {
	prev := History{c: player.c, facing: player.facing, status: board[player.c.x][player.c.y]}
	*history = append(*history, prev)
	possibleLeft := checkTurnLeft(player, board)
	possibleRight := checkTurnRight(player, board)
	boardStatus := strings.Split(board[player.c.x][player.c.y], ",")
	if len(boardStatus[0]) > 0 {
		sort.Strings(boardStatus)
		for _, status := range boardStatus {
			if status == "gold" {
				player.victory = true
				return player
			}
		}
		if boardStatus[0] == "breeze" && len(boardStatus) == 1 {
			if next := moveForward(player); (*pitIssues)[next.c] == "-" {
				return next
			}
			addPoW(player, history, pitIssues, "pit")
			checkDiagonal(history, player, pitIssues, wumpusIssues, "breeze", "pit")
			fmt.Println("Backtracking!")
			stepBack := backtrack(player, history)
			return lOrR(checkTurnLeft(stepBack, board), checkTurnRight(stepBack, board), stepBack, board)
		}
		if boardStatus[0] == "stench" && len(boardStatus) == 1 {
			if next := moveForward(player); (*wumpusIssues)[next.c] == "-" {
				return next
			}
			addPoW(player, history, wumpusIssues, "wumpus")
			checkDiagonal(history, player, wumpusIssues, pitIssues, "stench", "wumpus")
			return checkAndFire(player, wumpusIssues, pitIssues, board)
		}
		if len(boardStatus) == 2 && boardStatus[0] == "breeze" && boardStatus[1] == "stench" {
			fmt.Println("SPECIAL CASE!!!!")
			if next := moveForward(player); (*pitIssues)[next.c] == "-" {
				return next
			}
			return lOrR(possibleLeft, possibleRight, player, board)
			// stepBack := backtrack(player, history)
			// return lOrR(checkTurnLeft(stepBack, board), checkTurnRight(stepBack, board), stepBack, board)
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
					if player.facing != direction {
						player = align(player, direction, board)
					}
					player.moves++
					if fire(player, wumpusIssues, board) {
						player.score++
					}

				}
			}

		}
	}
	return player
}

func align(player gene, direction string, board *[4][4]string) gene {
	if player.facing == "N" {
		if direction == "S" {
			return turnLeft(turnLeft(player, board), board)
		}
		if direction == "W" {
			return turnLeft(player, board)
		}
		if direction == "E" {
			return turnRight(player, board)
		}
	} else if player.facing == "E" {
		if direction == "W" {
			return turnLeft(turnLeft(player, board), board)
		}
		if direction == "N" {
			return turnLeft(player, board)
		}
		if direction == "S" {
			return turnRight(player, board)
		}
	} else if player.facing == "W" {
		if direction == "E" {
			return turnLeft(turnLeft(player, board), board)
		}
		if direction == "S" {
			return turnLeft(player, board)
		}
		if direction == "N" {
			return turnRight(player, board)
		}
	} else if player.facing == "S" {
		if direction == "N" {
			return turnLeft(turnLeft(player, board), board)
		}
		if direction == "E" {
			return turnLeft(player, board)
		}
		if direction == "W" {
			return turnRight(player, board)
		}
	}
	return player
}

func fire(player gene, issues *map[coords]string, board *[4][4]string) bool {
	var flag = false
	var index = -1
	for {
		player = moveForward(player)
		fmt.Println(player)
		if (player.c.x > -1 && player.c.x < 4) && (player.c.y > -1 && player.c.y < 4) {
			status := strings.Split(board[player.c.x][player.c.y], ",")
			for i, val := range status {
				if val == "wumpus" {
					index = i
					flag = true
				}
			}
			if flag == true {
				board[player.c.x][player.c.y] = strings.Join(append(status[:index], status[index+1:]...), ",")
			}
			checkAndDelete(coords{x: player.c.x, y: player.c.y}, issues)
		} else {
			break
		}
	}
	return flag
}

func checkProximity(player gene, point coords) (bool, string) {

	if player.c.x-1 == point.x && player.c.y == point.y {
		return true, "E"
	}
	if player.c.x+1 == point.x && player.c.y == point.y {
		return true, "W"
	}
	if player.c.x == point.x && player.c.y-1 == point.y {
		return true, "S"
	}
	if player.c.x == point.x && player.c.y+1 == point.y {
		return true, "N"
	}
	return false, ""
}

func checkDiagonal(history *[]History, player gene, issues *map[coords]string, secondaryIssues *map[coords]string, value1 string, value2 string) {

	var value3 string
	if value1 == "breeze" {
		value3 = "stench"
	} else {
		value3 = "breeze"
	}
	if len(*history) < 2 {
		return
	}
	var flag = false
	from := (*history)[len(*history)-2].c
	for _, value := range *history {
		var status = strings.Split(value.status, ",")
		for _, v := range status {
			if v == value1 {
				flag = true
			}
		}
		if flag == true {
			diags := getDiags(value.c)
			fmt.Println(diags)
			for _, val := range diags {
				fmt.Println(val)
				if checkHistoryForStatus(val, history) == value1 {
					if (val == coords{x: value.c.x - 1, y: value.c.y - 1}) {
						fmt.Println("Case 1")
						if value.c.x-1 != from.x &&
							value.c.y != from.y &&
							(*issues)[coords{x: value.c.x - 1, y: value.c.y}] != "-" {
							(*issues)[coords{x: value.c.x - 1, y: value.c.y}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						if value.c.x+1 != from.x &&
							value.c.y-1 != from.y &&
							(*issues)[coords{x: value.c.x + 1, y: value.c.y - 1}] != "-" {
							(*issues)[coords{x: value.c.x + 1, y: value.c.y - 1}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
					} else if val.x == value.c.x+1 && val.y == value.c.y-1 {
						if value.c.x != from.x &&
							value.c.y-1 != from.y &&
							(*issues)[coords{x: value.c.x, y: value.c.y - 1}] != "-" {
							(*issues)[coords{x: value.c.x, y: value.c.y - 1}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						if value.c.x+1 != from.x &&
							value.c.y != from.y &&
							(*issues)[coords{x: value.c.x + 1, y: value.c.y}] != "-" {
							(*issues)[coords{x: value.c.x + 1, y: value.c.y}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
					} else if (val == coords{x: value.c.x + 1, y: value.c.y + 1}) {
						if value.c.x != from.x &&
							value.c.y+1 != from.y &&
							(*issues)[coords{x: value.c.x, y: value.c.y + 1}] != "-" {
							(*issues)[coords{x: value.c.x, y: value.c.y + 1}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						if value.c.x+1 != from.x &&
							value.c.y != from.y &&
							(*issues)[coords{x: value.c.x + 1, y: value.c.y}] != "-" {
							(*issues)[coords{x: value.c.x + 1, y: value.c.y}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
					} else if (val == coords{x: value.c.x - 1, y: value.c.y + 1}) {
						if value.c.x-1 != from.x &&
							value.c.y != from.y &&
							(*issues)[coords{x: value.c.x - 1, y: value.c.y}] != "-" {
							(*issues)[coords{x: value.c.x - 1, y: value.c.y}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
						if value.c.x != from.x &&
							value.c.y+1 != from.y &&
							(*issues)[coords{x: value.c.x, y: value.c.y + 1}] != "-" {
							(*issues)[coords{x: value.c.x, y: value.c.y + 1}] = value2
						} else {
							checkAndDelete(coords{x: from.x, y: from.y}, issues)
						}
					}
				} else {
					var status = strings.Split(checkHistoryForStatus(val, history), ",")
					for _, v := range status {
						if v == value3 {

							if (val == coords{x: value.c.x - 1, y: value.c.y - 1}) {
								checkAndDelete(coords{x: value.c.x - 1, y: value.c.y}, issues)
								checkAndDelete(coords{x: value.c.x + 1, y: value.c.y - 1}, issues)
								checkAndDelete(coords{x: value.c.x - 1, y: value.c.y}, secondaryIssues)
								checkAndDelete(coords{x: value.c.x + 1, y: value.c.y - 1}, secondaryIssues)
							} else if val.x == value.c.x+1 && val.y == value.c.y-1 {
								checkAndDelete(coords{x: value.c.x, y: value.c.y - 1}, issues)
								checkAndDelete(coords{x: value.c.x + 1, y: value.c.y}, issues)
								checkAndDelete(coords{x: value.c.x, y: value.c.y - 1}, secondaryIssues)
								checkAndDelete(coords{x: value.c.x + 1, y: value.c.y}, secondaryIssues)
							} else if (val == coords{x: value.c.x + 1, y: value.c.y + 1}) {
								checkAndDelete(coords{x: value.c.x, y: value.c.y + 1}, issues)
								checkAndDelete(coords{x: value.c.x + 1, y: value.c.y}, issues)
								checkAndDelete(coords{x: value.c.x, y: value.c.y + 1}, secondaryIssues)
								checkAndDelete(coords{x: value.c.x + 1, y: value.c.y}, secondaryIssues)
							} else if (val == coords{x: value.c.x - 1, y: value.c.y + 1}) {
								checkAndDelete(coords{x: value.c.x - 1, y: value.c.y}, issues)
								checkAndDelete(coords{x: value.c.x, y: value.c.y + 1}, issues)
								checkAndDelete(coords{x: value.c.x - 1, y: value.c.y}, secondaryIssues)
								checkAndDelete(coords{x: value.c.x, y: value.c.y + 1}, secondaryIssues)
							}

						}
					}
				}
			}
		}
	}
}

func checkAndDelete(point coords, issues *map[coords]string) {
	(*issues)[point] = "-"
}

func checkHistoryForStatus(point coords, history *[]History) string {
	for _, value := range *history {
		if value.c.x == point.x && value.c.y == point.y {
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
	from := (*history)[len(*history)-2].c
	fmt.Println("From::")
	fmt.Println(from)
	if (player.c.x-1 > -1 && player.c.x-1 < 4) && (player.c.x-1 == from.x && player.c.y == from.y) == false {
		(*issues)[coords{x: player.c.x - 1, y: player.c.y}] = value
	}
	if (player.c.x+1 > -1 && player.c.x+1 < 4) && (player.c.x+1 == from.x && player.c.y == from.y) == false {
		(*issues)[coords{x: player.c.x + 1, y: player.c.y}] = value
	}
	if (player.c.y-1 > -1 && player.c.y-1 < 4) && (player.c.x == from.x && player.c.y-1 == from.y) == false {
		(*issues)[coords{x: player.c.x, y: player.c.y - 1}] = value
	}
	if (player.c.y+1 > -1 && player.c.y+1 < 4) && (player.c.x == from.x && player.c.y+1 == from.y) == false {
		(*issues)[coords{x: player.c.x, y: player.c.y + 1}] = value
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
	player.c = (*history)[max-2].c
	player.facing = (*history)[max-2].facing
	player.moves = player.moves + 1
	return player
}

func lOrR(possibleLeft bool, possibleRight bool, player gene, board *[4][4]string) gene {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	if possibleLeft == true {
		if possibleRight == true {
			if r1.Intn(2) == 0 {
				return turnLeft(player, board)
			}
			return turnRight(player, board)
		}
		return turnLeft(player, board)
	}
	return turnRight(player, board)
}

func checkMoveValidity(player gene) bool {
	if player.c.x > -1 && player.c.x < 4 && player.c.y > -1 && player.c.y < 4 {
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
		if player.c.x-1 > -1 {
			return true
		}
	} else if player.facing == "S" {
		if player.c.x+1 < 4 {
			return true
		}
	} else if player.facing == "E" {
		if player.c.y+1 < 4 {
			return true
		}
	} else {
		if player.c.y-1 > -1 {
			return true
		}
	}
	return false
}

func checkTurnRight(player gene, board *[4][4]string) bool {
	if player.facing == "N" {
		if player.c.x+1 < 4 {
			return true
		}
	} else if player.facing == "S" {
		if player.c.x-1 > -1 {
			return true
		}
	} else if player.facing == "W" {
		if player.c.y+1 < 4 {
			return true
		}
	} else {
		if player.c.y-1 > -1 {
			return true
		}
	}
	return false
}
