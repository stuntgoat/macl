package main

type Direction string

var Up = Direction("Up")
var Down = Direction("Down")

var Left = Direction("Left")
var Right = Direction("Right")

var UpLeft = Direction("UpLeft")
var DownRight = Direction("DownRight")

var UpRight = Direction("UpRight")
var DownLeft = Direction("DownLeft")

type Line string

var UpDown = Line("UpDown")
var LeftRight = Line("LeftRight")
var DiagonalLR_UD = Line("Diagonal_Left_Right_Up_Down")
var DiagonalLR_DU = Line("Diagonal_Left_Right_Down_Up")

var LINE_FOR_DIRECTION = map[Direction]Line{
	Up:   UpDown,
	Down: UpDown,

	Left:  LeftRight,
	Right: LeftRight,

	UpLeft:    DiagonalLR_UD,
	DownRight: DiagonalLR_UD,

	UpRight:  DiagonalLR_DU,
	DownLeft: DiagonalLR_DU,
}

type LineKey struct {
	Row  int
	Col  int
	Line Line
}

type CoinKey struct {
	Row int
	Col int
}
type DirectionKey struct {
	Row       int
	Col       int
	Direction Direction
}

type PlayerGraph struct {
	coins map[CoinKey]bool
}

func (pg *PlayerGraph) Add(row, col int) {
	pg.coins[CoinKey{row, col}] = true
}

func (pg *PlayerGraph) Get(row, col int) bool {
	return pg.coins[CoinKey{row, col}]
}

func mkNextDirectionKey(row, col int, direction Direction) DirectionKey {
	var nextRow int
	var nextCol int
	switch direction {
	case Up:
		nextRow = row - 1
		nextCol = col
	case Down:
		nextRow = row + 1
		nextCol = col

	case Left:
		nextRow = row
		nextCol = col - 1
	case Right:
		nextRow = row
		nextCol = col + 1

	case UpLeft:
		nextRow = row - 1
		nextCol = col - 1
	case DownRight:
		nextRow = row + 1
		nextCol = col + 1

	case DownLeft:
		nextRow = row + 1
		nextCol = col - 1
	case UpRight:
		nextRow = row - 1
		nextCol = col + 1
	}
	return DirectionKey{nextRow, nextCol, direction}
}

func dfs(key DirectionKey, lkey LineKey, graph *PlayerGraph, count int, seen map[LineKey]bool) int {
	// Mark coordinate as seen for this line.
	seen[lkey] = true

	nextDirectionKey := mkNextDirectionKey(key.Row, key.Col, key.Direction)
	if ok := graph.Get(nextDirectionKey.Row, nextDirectionKey.Col); ok {
		nextLineKey := LineKey{nextDirectionKey.Row, nextDirectionKey.Col, lkey.Line}
		count = dfs(nextDirectionKey, nextLineKey, graph, count+1, seen)
	}
	return count
}

func searchDirections(dKeyA, dKeyB DirectionKey, lkey LineKey, graph *PlayerGraph, seen map[LineKey]bool) int {
	result := dfs(dKeyA, lkey, graph, 0, seen)
	result += dfs(dKeyB, lkey, graph, 0, seen)
	return result
}

// FindConsecutive checks for `num` consecutive DirectionKeys, in
// in the sum of opposite directions, in the graph by performing depth first search
// on the CoinKeys of this graph.
func (pg *PlayerGraph) FindConsecutive(row, col, num int) bool {
	var result int
	var dKeyA DirectionKey
	var dKeyB DirectionKey
	var lkey LineKey

	// Store the line slope direction per coordinate so we skip already visited nodes for each
	// direction search.
	seen := map[LineKey]bool{}

	lkey = LineKey{row, col, LINE_FOR_DIRECTION[Up]}
	if !seen[lkey] {
		dKeyA = DirectionKey{row, col, Up}
		dKeyB = DirectionKey{row, col, Down}
		result = searchDirections(dKeyA, dKeyB, lkey, pg, seen)
		if result+1 >= num {
			return true
		}
	}

	lkey = LineKey{row, col, LINE_FOR_DIRECTION[Left]}
	if !seen[lkey] {
		dKeyA = DirectionKey{row, col, Left}
		dKeyB = DirectionKey{row, col, Right}
		result = searchDirections(dKeyA, dKeyB, lkey, pg, seen)
		if result+1 >= num {
			return true
		}
	}

	lkey = LineKey{row, col, LINE_FOR_DIRECTION[UpLeft]}
	if !seen[lkey] {
		dKeyA = DirectionKey{row, col, UpLeft}
		dKeyB = DirectionKey{row, col, DownRight}
		result = searchDirections(dKeyA, dKeyB, lkey, pg, seen)
		if result+1 >= num {
			return true
		}
	}

	lkey = LineKey{row, col, LINE_FOR_DIRECTION[UpRight]}
	if !seen[lkey] {
		dKeyA = DirectionKey{row, col, UpRight}
		dKeyB = DirectionKey{row, col, DownLeft}
		result = searchDirections(dKeyA, dKeyB, lkey, pg, seen)
		if result+1 >= num {
			return true
		}

	}
	return false
}
