package main

import (
	"fmt"
	"log"
	"math"
)

type move int

type piece string

func (p piece) opposite() piece {
	if p == "B" {
		return "R"
	} else if p == "R" {
		return "B"
	} else {
		return " "
	}
}

type column struct {
	container []piece
	numRows   int
}

func (c *column) full() bool {
	return len(c.container) == c.numRows
}

func (c *column) push(item piece) {
	if c.full() {
		log.Fatal("Trying to push piece to full column")
	}
	c.container = append(c.container, item)
}

func (c *column) item(index int) piece {
	if index > len(c.container)-1 {
		return piece(" ")
	}
	return c.container[index]
}

func (c *column) copy() column {
	temp := make([]piece, len(c.container))
	copy(temp, c.container)
	return column{temp, c.numRows}
}

type tuple [2]int

type board struct {
	numRows       int
	numColumns    int
	segmentLength int
	segments      [][]tuple
	position      []column
	turn          piece
}

func (b *board) init(position []column, turn piece) {
	b.numColumns = 7
	b.numRows = 6
	b.segmentLength = 4
	b.turn = turn
	b.segments = generateSegments(b.numColumns, b.numRows, b.segmentLength)
	if position == nil {
		for i := 0; i < b.numColumns; i++ {
			b.position = append(b.position, column{[]piece{}, b.numRows})
		}
	} else {
		b.position = position
	}
}

func generateSegments(numColumns int, numRows int, segmentLength int) [][]tuple {
	segments := [][]tuple{}
	for c := 0; c < numColumns; c++ {
		for r := 0; r < numRows-segmentLength+1; r++ {
			segment := []tuple{}
			for t := 0; t < segmentLength; t++ {
				segment = append(segment, tuple{c, r + t})
			}
			segments = append(segments, segment)
		}
	}

	for c := 0; c < numColumns-segmentLength+1; c++ {
		for r := 0; r < numRows; r++ {
			segment := []tuple{}
			for t := 0; t < segmentLength; t++ {
				segment = append(segment, tuple{c + t, r})
			}
			segments = append(segments, segment)
		}
	}

	for c := 0; c < numColumns-segmentLength+1; c++ {
		for r := 0; r < numRows-segmentLength+1; r++ {
			segment := []tuple{}
			for t := 0; t < segmentLength; t++ {
				segment = append(segment, tuple{c + t, r + t})
			}
			segments = append(segments, segment)
		}
	}

	for c := 0; c < numColumns-segmentLength+1; c++ {
		for r := segmentLength - 1; r < numRows; r++ {
			segment := []tuple{}
			for t := 0; t < segmentLength; t++ {
				segment = append(segment, tuple{c + t, r - t})
			}
			segments = append(segments, segment)
		}
	}
	return segments
}

func (b *board) getTurn() piece {
	return b.turn
}

func (b *board) copyPosition() []column {
	ret := []column{}
	for c := 0; c < len(b.position); c++ {
		pieces := make([]piece, len(b.position[c].container))
		col := column{pieces, b.numRows}
		copy(col.container, b.position[c].container)
		ret = append(ret, col)
	}
	return ret
}

func (b *board) move(location move) board {
	tempPosition := b.copyPosition()
	tempPosition[location].push(b.turn)
	newboard := board{}
	newboard.init(tempPosition, b.turn.opposite())
	return newboard
}

func (b *board) legalMoves() []move {
	ret := []move{}
	for c := 0; c < b.numColumns; c++ {
		if !b.position[c].full() {
			ret = append(ret, move(c))
		}
	}
	return ret
}

func (b *board) isWin() bool {
	for _, s := range b.segments {
		blackCount, redCount := b.countSegment(s)
		if blackCount == 4 || redCount == 4 {
			return true
		}
	}
	return false
}

func (b *board) isDraw() bool {
	return !b.isWin() && (len(b.legalMoves()) == 0)
}

func (b *board) countSegment(segment []tuple) (int, int) {
	blackCount := 0
	redCount := 0
	for _, s := range segment {
		if b.position[s[0]].item(s[1]) == piece("B") {
			blackCount++
		} else if b.position[s[0]].item(s[1]) == piece("R") {
			redCount++
		}
	}
	return blackCount, redCount
}

func (b *board) evaluateSegment(segment []tuple, player piece) float64 {
	blackCount, redCount := b.countSegment(segment)
	if blackCount > 0 && redCount > 0 {
		return 0.0
	}
	count := max(float64(redCount), float64(blackCount))
	score := 0.0

	if count == 2.0 {
		score = 1.0
	} else if count == 3.0 {
		score = 100.0
	} else if count == 4.0 {
		score = 1000000.0
	}
	color := piece("B")
	if redCount > blackCount {
		color = piece("R")
	}
	if color != player {
		score = -score
	}

	return score
}

func (b *board) evaluate(player piece) float64 {
	total := 0.0
	for _, segment := range b.segments {
		total += b.evaluateSegment(segment, player)
	}
	return total
}

func (b board) String() string {
	str := ""
	for r := b.numRows - 1; r >= 0; r-- {
		str += "|"
		for c := 0; c < b.numColumns; c++ {
			str += fmt.Sprintf(" %s |", b.position[c].item(r))
		}
		str += "\n"
	}
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

func alphaBeta(b board, maximizing bool, originalPlayer piece, maxDepth int, alpha float64, beta float64) float64 {
	if b.isWin() || b.isDraw() || maxDepth == 0 {
		return b.evaluate(originalPlayer)
	}

	if maximizing {
		for _, m := range b.legalMoves() {
			result := alphaBeta(b.move(m), false, originalPlayer, maxDepth-1, alpha, beta)
			alpha = max(result, alpha)
			if beta <= alpha {
				break
			}
		}
		return alpha
	} else {
		for _, m := range b.legalMoves() {
			result := alphaBeta(b.move(m), true, originalPlayer, maxDepth-1, alpha, beta)
			beta = min(result, beta)
			if beta <= alpha {
				break
			}
		}
		return beta
	}
}

func findBestMove(b board, maxDepth int) move {
	bestEval := -math.MaxFloat64
	bestMove := move(-1)
	for _, m := range b.legalMoves() {
		// result := minimax(b.move(m), false, b.turn, maxDepth)
		result := alphaBeta(b.move(m), false, b.turn, maxDepth, -math.MaxFloat64, math.MaxFloat64)
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
		fmt.Printf("Enter a legal column (0-6): ")
		fmt.Scanf("%d\n", &play)
		playermove = move(play)
	}
	return playermove
}

func main() {
	b := board{}
	b.init(nil, piece("B"))
	fmt.Println(b)
	for {
		humamMove := getPlayerMove(b)
		b = b.move(humamMove)
		fmt.Println(b)
		if b.isWin() {
			fmt.Println("Human wins!")
			break
		} else if b.isDraw() {
			fmt.Println("Draw!")
			break
		}
		computerMove := findBestMove(b, 5)
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
