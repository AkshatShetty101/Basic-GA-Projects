package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jinzhu/copier"
)

// To define x and y co-ordinates
type Position struct {
	x int
	y int
}

type Status struct {
	level int
	moves int
}

type Instance struct {
	board           [height][width]rune
	threats         map[Position]Status
	attacks         map[Position]Status
	chances         [5]bool
	playerVictory   bool
	computerVictory bool
}

const height = 6
const width = 7

var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

// func tcpServer() {
// 	fmt.Println("Launching server...")

// 	// listen on all interfaces
// 	ln, _ := net.Listen("tcp", ":8081")

// 	// accept conne4ion on port
// 	conn, _ := ln.Accept()

// 	// run loop forever (or until 4rl-c)
// 	// for {
// 	// 	// will listen for message to process ending in newline (\n)
// 	// 	message, _ := bufio.NewReader(conn).ReadString('\n')
// 	// 	// output message received
// 	// 	fmt.Print("Message Received:", string(message))
// 	// 	// nextMove := move()
// 	// 	newmessage := strings.ToUpper(nextMove)
// 	// 	// send new string back to client
// 	// 	conn.Write([]byte(newmessage + "\n"))
// 	// }

// }
func main() {

	original := setBoard()
	for {
		original = *moveComputer(&original)
		display(&original)
		mimicPlayerMove(&original)
		evaluateStatus(&original)
		display(&original)
	}

}

func setupThreatStatus(game *Instance) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			getPlayerChances(Position{x: i, y: j}, game)
		}
	}
}

func setupAttackStatus(game *Instance) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			getComputerChances(Position{x: i, y: j}, game)
		}
	}
}

func evaluateStatus(game *Instance) {
	setupThreatStatus(game)
	setupAttackStatus(game)
	if game.playerVictory == true {
		display(game)
		fmt.Println("Player wins!!")
		os.Exit(1)
	}
	if game.computerVictory == true {
		display(game)
		fmt.Println("Computer wins!!")
		os.Exit(1)
	}

}

func mimicPlayerMove(game *Instance) {
	for {
		var col int
		fmt.Println("Which column to place coin in?")
		fmt.Scan(&col)
		// for {
		// 	col := r1.Intn(width - 1)
		// fmt.Println(col)
		if dropCoin(rune('o'), col, &game.board) {
			break
		}
		fmt.Println("Column full! Please enter column again")
	}
}

func display(game *Instance) {
	for i := 0; i < height; i++ {
		fmt.Print("|")
		for j := 0; j < width; j++ {
			fmt.Print(" " + string(game.board[i][j]) + " | ")
		}
		fmt.Println()
	}
	fmt.Println("Threats: ", game.threats)
	fmt.Println("Attacks: ", game.attacks)
}

func getMovesLeft(board *[height][width]rune, pos *Position) int {
	i := pos.x
	for ; i < height; i++ {
		if board[i][pos.y] != ' ' {
			break
		}
	}
	return i - pos.x
}

func getFitness(game *Instance) float64 {
	level3T := 0.0
	level2T := 0.0
	level2A := 0.0
	level3A := 0.0
	if game.computerVictory == true {
		display(game)
		fmt.Println("Computer wins")
		os.Exit(1)
	}
	for _, val := range game.attacks {
		if val.level == 3 {
			level3A += 1.0 / float64(val.moves)
		}
		level2A += 1.0 / float64(val.moves)

	}
	for _, val := range game.threats {
		if val.level == 3 {
			level3T += 1.0 / float64(val.moves)
		}
		level2T += 1.0 / float64(val.moves)
	}
	return level3A/10 + level2A/20 - level3T - level2T/2
}

func moveComputer(game *Instance) *Instance {
	parentFitness := getFitness(game)
	var nextP = Instance{}
	fmt.Println("parent:", parentFitness)
	// if status, _ := checkLevelThree(&game.threats); status == true {
	// 	fmt.Println("Major Threat detected!")
	// 	fitness := parentFitness
	// 	for i := 0; i < width; i++ {
	// 		if flag, child := getMoveEffects(i, game); flag == true {
	// 			// display(child)
	// 			childFitness := getFitness(child)
	// 			// fmt.Println("Child:", childFitness)
	// 			// display(child)
	// 			if childFitness >= fitness {
	// 				fitness = getFitness(child)
	// 				nextP = *child
	// 			}

	// 		}
	// 	}
	// } else {
	fitness := -100.0
	for i := 0; i < width; i++ {
		if flag, child := getMoveEffects(i, game); flag == true {
			// display(child)
			childFitness := getFitness(child)
			// fmt.Println("Child:", childFitness)
			// display(child)
			if childFitness >= fitness {
				fitness = getFitness(child)
				nextP = *child
			}

		}
	}
	// }
	return &nextP
}

func getMoveEffects(column int, game *Instance) (bool, *Instance) {
	var cloneG = Instance{}
	copier.Copy(&cloneG, game)
	cloneG.threats = make(map[Position]Status)
	cloneG.attacks = make(map[Position]Status)
	if dropCoin(rune('#'), column, &cloneG.board) == true {
		evaluateStatus(&cloneG)
		return true, &cloneG
	}
	return false, &Instance{}
}

func dropCoin(ch rune, column int, board *[height][width]rune) bool {
	i := 0
	if board[i][column] != rune(' ') {
		fmt.Println("Column full")
		return false
	}
	for ; i < height; i++ {
		if board[i][column] != rune(' ') {
			break
		}
	}
	i = i - 1
	board[i][column] = ch
	return true
}

func setBoard() Instance {
	var x = Instance{}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			x.board[i][j] = rune(' ')
		}
	}
	x.computerVictory = false
	x.playerVictory = false
	x.threats = make(map[Position]Status)
	x.attacks = make(map[Position]Status)
	return x
}

func getComputerChances(p Position, game *Instance) bool {
	game.chances[0] = checkAbove(p, rune('o'), game)
	game.chances[1] = checkLeft(p, rune('o'), game)
	game.chances[2] = checkRight(p, rune('o'), game)
	game.chances[3] = checkRUDiag(p, rune('o'), game)
	game.chances[4] = checkLUDiag(p, rune('o'), game)
	if game.chances[0] == false && game.chances[1] == false && game.chances[2] == false && game.chances[3] == false && game.chances[4] == false {
		return false
	}
	return true
}

func getPlayerChances(p Position, game *Instance) bool {
	// fmt.Println("Pos: ", p)
	game.chances[0] = checkAbove(p, rune('#'), game)
	game.chances[1] = checkLeft(p, rune('#'), game)
	game.chances[2] = checkRight(p, rune('#'), game)
	game.chances[3] = checkRUDiag(p, rune('#'), game)
	game.chances[4] = checkLUDiag(p, rune('#'), game)
	// fmt.Println(game.threats)
	if game.chances[0] == false && game.chances[1] == false && game.chances[2] == false && game.chances[3] == false && game.chances[4] == false {
		return false
	}
	return true
}

func checkAbove(p Position, ch rune, game *Instance) bool {
	i := 0
	var pos []Position
	// fmt.Println(p)
	if game.board[p.x][p.y] != ch {
		// fmt.Println(p)
		for ; i < 4; i++ {
			if p.x < 0 || game.board[p.x][p.y] == ch {
				return false
			}
			if game.board[p.x][p.y] == rune(' ') {
				pos = append(pos, p)
			}
			p.x--
		}
		if i == 4 && len(pos) == 0 {
			if ch == rune('#') {
				game.playerVictory = true
			}
			game.computerVictory = true
			return false
		}
		checkLevelAndStore(&pos, ch, game)
		return true
	}
	return false
}

func checkLeft(p Position, ch rune, game *Instance) bool {
	i := 0
	var pos []Position
	if game.board[p.x][p.y] != ch {
		for ; i < 4; i++ {
			if p.y < 0 || game.board[p.x][p.y] == ch {
				return false
			}
			if game.board[p.x][p.y] == rune(' ') {
				pos = append(pos, p)
			}
			p.y--
		}
		if i == 4 && len(pos) == 0 {
			if ch == rune('#') {
				game.playerVictory = true
			}
			game.computerVictory = true
			return false
		}
		checkLevelAndStore(&pos, ch, game)
		return true
	}
	return false
}

func checkRight(p Position, ch rune, game *Instance) bool {
	i := 0
	var pos []Position
	if game.board[p.x][p.y] != ch {
		for ; i < 4; i++ {
			if p.y >= width || game.board[p.x][p.y] == ch {
				return false
			}
			if game.board[p.x][p.y] == rune(' ') {
				pos = append(pos, p)
			}
			p.y++
		}
		if i == 4 && len(pos) == 0 {
			if ch == rune('#') {
				game.playerVictory = true
			}
			game.computerVictory = true
			return false
		}
		checkLevelAndStore(&pos, ch, game)
		return true
	}
	return false
}

func checkRUDiag(p Position, ch rune, game *Instance) bool {
	i := 0
	var pos []Position
	if game.board[p.x][p.y] != ch {
		for ; i < 4; i++ {
			if p.y >= width || p.x < 0 || game.board[p.x][p.y] == ch {
				return false
			}
			if game.board[p.x][p.y] == rune(' ') {
				pos = append(pos, p)
			}
			p.x--
			p.y++
		}
		if i == 4 && len(pos) == 0 {
			if ch == rune('#') {
				game.playerVictory = true
			}
			game.computerVictory = true
			return false
		}
		checkLevelAndStore(&pos, ch, game)
		return true
	}
	return false
}

func checkLUDiag(p Position, ch rune, game *Instance) bool {
	i := 0
	var pos []Position
	if game.board[p.x][p.y] != ch {
		for ; i < 4; i++ {
			if p.x < 0 || p.y < 0 || game.board[p.x][p.y] == ch {
				return false
			}
			if game.board[p.x][p.y] == rune(' ') {
				pos = append(pos, p)
			}
			p.y--
			p.x--
		}
		if i == 4 && len(pos) == 0 {
			if ch == rune('#') {
				game.playerVictory = true
			}
			game.computerVictory = true
			return false
		}
		checkLevelAndStore(&pos, ch, game)
		return true
	}
	return false
}

func checkLevelThree(toA *map[Position]Status) (bool, []Position) {
	var x = make([]Position, 0)
	flag := false
	for pos, val := range *toA {
		if val.level == 3 {
			flag = true
			x = append(x, pos)
		}
	}
	return flag, x
}

func checkLevelAndStore(pos *[]Position, ch rune, game *Instance) {
	// fmt.Println("GOT: ", pos)
	var toA *map[Position]Status
	if ch == rune('o') {
		toA = &game.attacks
	} else {
		toA = &game.threats
	}
	if len(*pos) == 2 {
		for _, value := range *pos {
			if val, ok := (*toA)[value]; ok && val.level == 3 {
				continue
			}
			(*toA)[value] = Status{level: 2, moves: getMovesLeft(&game.board, &value)}
		}
	} else if len(*pos) == 1 {
		(*toA)[(*pos)[0]] = Status{level: 3, moves: getMovesLeft(&game.board, &(*pos)[0])}
	}
}
