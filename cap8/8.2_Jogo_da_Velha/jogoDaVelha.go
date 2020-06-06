package main

import (
	"fmt"
	"math"
)

type move int

type piece string

func (p piece) opposite() piece {
	if p == "X" {
		return "O"
	} else if p == "O" {
		return "X"
	} else {
		return " "
	}
}

type board struct {
	position []piece
	turn     piece
}

func (b *board) init() {
	turn := piece("X")
	for i := 0; i < 9; i++ {
		b.position = append(b.position, piece(" "))
	}
	b.turn = turn
}

func (b *board) getTurn() piece {
	return b.turn
}

func (b *board) copyPosition() []piece {
	ret := []piece{}
	ret = append(ret, b.position...)
	copy(ret, b.position)
	return ret
}

func (b *board) move(location move) board {
	tempPosition := b.copyPosition()
	tempPosition[location] = b.turn
	return board{tempPosition, b.turn.opposite()}
}

func (b *board) legalMoves() []move {
	ret := []move{}
	for l := 0; l < len(b.position); l++ {
		if b.position[l] == piece(" ") {
			ret = append(ret, move(l))
		}
	}
	return ret
}

func (b *board) isDraw() bool {
	return !b.isWin() && (len(b.legalMoves()) == 0)
}

func (b *board) isWin() bool {
	return b.position[0] == b.position[1] &&
		b.position[0] == b.position[2] &&
		b.position[0] != piece(" ") ||
		b.position[3] == b.position[4] &&
			b.position[3] == b.position[5] &&
			b.position[3] != piece(" ") ||
		b.position[6] == b.position[7] &&
			b.position[6] == b.position[8] &&
			b.position[6] != piece(" ") ||
		b.position[0] == b.position[3] &&
			b.position[0] == b.position[6] &&
			b.position[0] != piece(" ") ||
		b.position[1] == b.position[4] &&
			b.position[1] == b.position[7] &&
			b.position[1] != piece(" ") ||
		b.position[2] == b.position[5] &&
			b.position[2] == b.position[8] &&
			b.position[2] != piece(" ") ||
		b.position[0] == b.position[4] &&
			b.position[0] == b.position[8] &&
			b.position[0] != piece(" ") ||
		b.position[2] == b.position[4] &&
			b.position[2] == b.position[6] &&
			b.position[2] != piece(" ")
}

func (b *board) evaluate(player piece) float64 {
	if b.isWin() && b.turn == player {
		return -1.0
	} else if b.isWin() && b.turn != player {
		return 1.0
	} else {
		return 0.0
	}
}

func (b board) String() string {
	str := ""
	str += fmt.Sprintf(" %s | %s | %s\n", b.position[0], b.position[1], b.position[2])
	str += fmt.Sprintf("---------\n")
	str += fmt.Sprintf(" %s | %s | %s\n", b.position[3], b.position[4], b.position[5])
	str += fmt.Sprintf("---------\n")
	str += fmt.Sprintf(" %s | %s | %s\n", b.position[6], b.position[7], b.position[8])
	return str
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func minimax(b board, maximizing bool, originalPlayer piece, maxDepth int) float64 {
	if b.isWin() || b.isDraw() || maxDepth == 0 {
		return b.evaluate(originalPlayer)
	}

	if maximizing {
		bestEval := -math.MaxFloat64
		for _, m := range b.legalMoves() {
			result := minimax(b.move(m), false, originalPlayer, maxDepth-1)
			bestEval = max(result, bestEval)
		}
		return bestEval
	} else {
		worstEval := math.MaxFloat64
		for _, m := range b.legalMoves() {
			result := minimax(b.move(m), true, originalPlayer, maxDepth-1)
			worstEval = min(result, worstEval)
		}
		return worstEval
	}
}

func findBestMove(b board, maxDepth int) move {
	bestEval := -math.MaxFloat64
	bestMove := move(-1)

	for _, m := range b.legalMoves() {
		result := minimax(b.move(m), false, b.turn, maxDepth)
		if result > bestEval {
			bestEval = result
			bestMove = m
		}
	}
	return bestMove
}

func in(a move, list []move) bool {
	for _, l := range list {
		if a == l {
			return true
		}
	}
	return false
}

func getPlayerMove(b board) move {
	var play int
	playermove := move(-1)
	for !in(playermove, b.legalMoves()) {
		fmt.Printf("Enter a legal square (0-8): ")
		fmt.Scanf("%d\n", &play)
		playermove = move(play)
	}
	return playermove
}

func main() {
	b := board{}
	b.init()
	fmt.Println(b)
	for {
		humamMove := getPlayerMove(b)
		b = b.move(humamMove)
		if b.isWin() {
			fmt.Println("Human wins!")
			break
		} else if b.isDraw() {
			fmt.Println("Draw!")
			break
		}
		computerMove := findBestMove(b, 8)
		fmt.Printf("\nComputer move is %d\n\n", int(computerMove))
		b = b.move(computerMove)
		fmt.Println(b)
		if b.isWin() {
			fmt.Println("Computer wins!")
			break
		} else if b.isDraw() {
			fmt.Println("Draw!")
			break
		}
	}
}
