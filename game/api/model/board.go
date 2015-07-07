package model

import (
	"errors"
	"encoding/json"
)

func (self *TicTacToeBoard) Reset() {
	reset()
}

func (self *TicTacToeBoard) PerformMove(  playerName string, x, y int ) error {
	err := validateMove(x, y)
	board[x][y] = playerName;
	moves++;
	computeState();
	return err;
}

func (self *TicTacToeBoard) GetState() string {
	return state
}

func (self *TicTacToeBoard) Serialize() []byte {
	var object map[string]string
	object = make(map[string]string)

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			key := makeKey(x, y);
			value := board[x][y];
			if ( value != nil ) {
				object[key] = value;
			}
		}
	}
	object["moves"] = moves
	jsonBytes, err := json.Marshal(object)

	return jsonBytes, err
}

func (self *TicTacToeBoard) Deserialize(jsonBytes []byte) {
	// bytes -> map
	var object map[string]string = map[string]string(jsonBytes)

	reset();
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			key := makeKey(x, y);
			value := object[key];
			if ( value != nil ) {
				board[x][y] = value;
			}
		}
	}
	moves = object["moves"];
	if ( moves == nil ) {
	 	moves = 0;
	} else {
		moves = int(moves);
	}

	computeState();
}

func (self *TicTacToeBoard) GetWinner() string {
	return winner
}

func (self *TicTacToeBoard) ComputeWinner() string {
	return computeWinner()
}
func (self *TicTacToeBoard) GetWinningCells() []string {
	return winningCells
}

///////////////////////////////////////////////////////////////////////////////////////////
// private

type TicTacToeBoard struct {}

// state
var board [][]string
var moves int
var winner string
var winningCells []string
var state string

// entry point
func NewBoard() *TicTacToeBoard {
	reset()
	return &(TicTacToeBoard{})
}

func reset() {
	board = nil
	moves = 0
	winner = nil
	winningCells = nil
	state = "init"
}

func validateMove( x, y int) error {
	if ( x < 0 || y < 0 ) {
		errors.New("Out of bounds: " + x + ", " + y);
	}

	if ( x > 3 || y > 3 ) {
		errors.New("Out of bounds: " + x + ", " + y);
	}

	if ( board[x][y] != nil ) {
		errors.New("Out of bounds: " + x + ", " + y);
	}

	return nil;
}

func computeState() {
	winner := computeWinner()
	if ( winner == nil && moves > 8) {
		state = "stalemate"
		return
	}
	if ( winner == nil ) {
		state = "running"
	} else {
		state = "won"
	}
}

func computeWinner() string {
	var current string = nil
	var last string = nil

	// check horizontal
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x ++  {
			current = board[x][y]
			if ( current == nil ) {
				// empty slot means short circuit row check
				break
			}

			if ( last != nil && last != current ) {
				break
			}

			last = current
		}
		if ( last != nil && current != nil && last == current ) {
			append( winningCells, makeKey(0,y) )
			append( winningCells, makeKey(1,y) )
			append( winningCells, makeKey(2,y) )
			return last
		}

		last = nil
	}

	// check vertical
	last = nil
	current = nil
	winningCells = nil
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x ++ {
			current = board[y][x];
			if ( current == nil ) {
				// empty slot means short circuit row check
				break
			}

			if ( last != nil && last != current ) {
				break
			}

			last = current
		}
		if ( last != nil && current != nil && last == current ) {
			append( winningCells, makeKey(y,0) )
			append( winningCells, makeKey(y,1) )
			append( winningCells, makeKey(y,2) )
			return last
		}
		last = nil
	}

	// check diagonal
	last = nil
	current = nil
	winningCells = nil
	for i := 0; i < 3; i++ {
		current = board[i][i];
		if ( current == nil ) {
		// empty slot means short circuit row check
			break
		}

		if ( last != nil && last != current ) {
			break
		}

		last = current
		}
		if ( last != nil && current != nil && last == current ) {
			append( winningCells, makeKey(0,0) )
			append( winningCells, makeKey(1,1) )
			append( winningCells, makeKey(2,2) )
			return last
		}

	// check diagonal
	last = nil
	current = nil
	winningCells = nil
	for i := 2; i >= 0; i-- {
		x := i
		y := 2 - i
		current = board[x][y]
		if ( current == nil ) {
		// empty slot means short circuit row check
			break
		}

		if ( last != nil && last != current ) {
			break
		}

		last = current
	}

	if ( last != nil && current != nil && last == current ) {
		append( winningCells, makeKey(2,0) )
		append( winningCells, makeKey(1,1) )
		append( winningCells, makeKey(0,2) )
		return last
	}

	winningCells = nil
	return nil
}

func makeKey(x, y string) string {
	return x + "_" + y
}



