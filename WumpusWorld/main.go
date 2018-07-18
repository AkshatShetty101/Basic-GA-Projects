package main

import (
	"fmt"
	"strings"
	"math/rand"
	"time"
	"math"
)

type gene struct {
	x int;
	y int;
	facing string;
	victory bool;
	fire bool;
}

func main() {
	var board [4][4]string;
	var playerLocation gene;
	var history []*gene;
	var goldLocation = gene{x:2,y:1,facing:"default",victory:true,fire:false}
	makeBoard(&board);
	fmt.Println("started");
	bestParent := setStartLocation(playerLocation);
	bestFitness := getFitness(bestParent,goldLocation)
	for{
		child := move(bestParent,history,&board);
		if child.victory == true{
			break;
		} else {
			childFitness := getFitness(child,goldLocation)
			if compareFitness(childFitness,bestFitness){
				bestParent = child;
				bestFitness = childFitness;
			}
		}
	}
}

func makeBoard(board *[4][4]string){

	board[0][0] = "start";
	board[0][1] = "breeze";
	board[0][1] = "pit";
	board[1][0] = "stench";
	board[0][3] = "breeze";
	board[1][2] = "breeze";
	board[2][0] = "wumpus";
	board[2][1] = "breeze,stench,gold";
	board[2][2] = "pit";
	board[2][3] = "breeze";
	board[3][0] = "stench";
	board[3][2] = "breeze";
	board[3][3] = "pit";	
}


func setStartLocation(q gene) gene {
	q.x = 0;
	q.y = 0;
	q.facing = "N";
	q.victory = false;
	q.fire = false;
	return q;
}

func getFitness(player gene,gold gene)float64{
	return math.Sqrt(math.Pow(float64(player.x-gold.x),2)+
	math.Pow(float64(player.y-gold.y),2));
}

func compareFitness(f1 float64,f2 float64)bool{
	return f1>f2;
}

func move(player gene,history []*gene, board *[4][4]string)gene{
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	possibleMove := gene{x:player.x,
		y:player.y+1,
		facing:player.facing,
		victory:player.victory,
		fire:false};
	history = append(history,&player);
	if (possibleMove.x>-1 && possibleMove.x<4) &&
	 (possibleMove.y>-1 && possibleMove.y<4){
		boardStatus := strings.Split(board[possibleMove.x][possibleMove.y], ",");
		if len(boardStatus)>0{
			for _, status := range boardStatus {
				if status == "gold" {
					possibleMove.victory = true;
					return possibleMove;
				}
			}
			for _, status := range boardStatus {
				if status == "wumpus" {
					player.fire = true;
					return player;
				} else
				if status == "pit"{
					possibleLeft := checkTurnLeft(player,board);
					possibleRight := checkTurnRight(player,board);
					if possibleLeft == true{
						if possibleRight == true{
							if r1.Intn(1) == 0{
								return TurnLeft(player,board);								
							} else {
								return TurnRight(player,board);
							}
						} else {
							return TurnRight(player,board);
						}
					} else {
						return TurnLeft(player,board);						
					}
				} 
			}
		} else {
			return possibleMove;
		}
	}
	return possibleMove;
}
func TurnLeft(player gene, board *[4][4]string) gene{
	if player.facing == "N" {
			player.facing  = "W"
			player.x = player.x-1;
	}else
	if player.facing == "S"{
			player.facing  = "E"
			player.x = player.x+1;
	}else
	if player.facing == "E"{
			player.facing  = "N"
			player.y = player.y+1;
	}else
	{
			player.y = player.y-1;
			player.facing  = "S"
	}
	return player;					
}
func TurnRight(player gene, board *[4][4]string) gene{
	if player.facing == "S" {
			player.facing  = "W"
			player.x = player.x-1;
	}else
	if player.facing == "N"{
			player.facing  = "E"
			player.x = player.x+1;
	}else
	if player.facing == "W"{
			player.facing  = "N"
			player.y = player.y+1;
	}else
	{
			player.y = player.y-1;
			player.facing  = "S"
	}
	return player;					
}

func checkTurnLeft(player gene, board *[4][4]string)bool{
	if player.facing == "N" {
		if player.x-1>-1{
			return true;
		}
	}else
	if player.facing == "S"{
		if player.x+1<4{
			return true;
		}
	}else
	if player.facing == "E"{
		if player.y+1<4{
			return true;
		}
	}else
	{
		if player.y-1>-1{
			return true;
		}
	}
	return false;					
}

func checkTurnRight(player gene, board *[4][4]string)bool{
	if player.facing == "N" {
		if player.x+1<4{
			return true;
		}
	}else
	if player.facing == "S"{
		if player.x-1>-1{
			return true;
		}
	}else
	if player.facing == "W"{
		if player.y+1<4{
			return true;
		}
	}else
	{
		if player.y-1>-1{
			return true;
		}
	}
	return false;					
}