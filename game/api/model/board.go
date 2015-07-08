package model

import (
	"errors"
	"encoding/json"
	"strconv"
)

func (self *TicTacToeBoard) Reset() {
	reset(self)
}

func (self *TicTacToeBoard) PerformMove(  playerName string, x, y int ) error {
	err := validateMove(self, x, y)
	self.board[x][y] = playerName;
	self.moves++;
	computeState(self);
	return err;
}

func (self *TicTacToeBoard) GetState() string {
	return self.state
}

func (self *TicTacToeBoard) Serialize() ([]byte, error) {
	var object map[string]string = make(map[string]string)

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			key := makeKey(x, y);
			value := self.board[x][y];
			if ( value != "" ) {
				object[key] = value;
			}
		}
	}
	object["moves"] = strconv.FormatInt(int64(self.moves), 10)
	jsonBytes, err := json.Marshal(object)


	return jsonBytes, err
}

func (self *TicTacToeBoard) Deserialize(jsonBytes []byte) error {
	// bytes -> map
	var object map[string]string = make( map[string]string )
	err := json.Unmarshal(jsonBytes, &object)
	if ( err !=  nil ) {
		return err
	}

	reset(self);

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			key := makeKey(x, y);
			value := object[key];
			if ( value != "" ) {
				self.board[x][y] = value;
			}
		}
	}
	if ( object["moves"] == "" ) {
		self.moves = 0;
	} else {
		self.moves, err = strconv.ParseInt(object["moves"], 10, 32);
	}

	computeState(self);

	return nil
}

func (self *TicTacToeBoard) GetWinner() string {
	return self.winner
}

func (self *TicTacToeBoard) ComputeWinner() string {
	return computeWinner(self)
}
func (self *TicTacToeBoard) GetWinningCells() []string {
	return self.winningCells
}

///////////////////////////////////////////////////////////////////////////////////////////
// private

type TicTacToeBoard struct {
	 board [][]string
	 moves int64
	 winner string
	 winningCells []string
	 state string
}

// entry point
func NewBoard() *TicTacToeBoard {
	theBoard := &(TicTacToeBoard{})
	reset(theBoard)
	return theBoard
}

func reset(self* TicTacToeBoard) {
	self.board = nil
	self.moves = 0
	self.winner = ""
	self.winningCells = nil
	self.state = "init"

	self.board = make([][]string, 3)
	for i := 0; i < 3; i++ {
		self.board[i] = make([]string, 3)
	}
}

func validateMove( self* TicTacToeBoard, x, y int) error {
	if ( x < 0 || y < 0 ) {
		errors.New("Out of bounds: " + strconv.Itoa(x) + ", " + strconv.Itoa(y));
	}

	if ( x > 3 || y > 3 ) {
		errors.New("Out of bounds: " + strconv.Itoa(x) + ", " + strconv.Itoa(y));
	}

	if ( self.board[x][y] != "" ) {
		errors.New("Out of bounds: " + strconv.Itoa(x) + ", " + strconv.Itoa(y));
	}

	return nil;
}

func computeState(self* TicTacToeBoard) {
	self.winner = computeWinner(self)
	if ( self.winner == "" && self.moves > 8) {
		self.state = "stalemate"
		return
	}
	if ( self.winner == "" ) {
		self.state = "running"
	} else {
		self.state = "won"
	}
}

func computeWinner(self* TicTacToeBoard) string {
	var current string = ""
	var last string = ""

	// check horizontal
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x ++  {
			current = self.board[x][y]
			if ( current == "" ) {
				// empty slot means short circuit row check
				break
			}

			if ( last != "" && last != current ) {
				break
			}

			last = current
		}
		if ( last != "" && current != "" && last == current ) {
			self.winningCells = append( self.winningCells, makeKey(0,y) )
			self.winningCells = append( self.winningCells, makeKey(1,y) )
			self.winningCells = append( self.winningCells, makeKey(2,y) )
			return last
		}

		last = ""
	}

	// check vertical
	last = ""
	current = ""
	self.winningCells = nil
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x ++ {
			current = self.board[y][x];
			if ( current == "" ) {
				// empty slot means short circuit row check
				break
			}

			if ( last != "" && last != current ) {
				break
			}

			last = current
		}
		if ( last != "" && current != "" && last == current ) {
			self.winningCells = append( self.winningCells, makeKey(y,0) )
			self.winningCells = append( self.winningCells, makeKey(y,1) )
			self.winningCells = append( self.winningCells, makeKey(y,2) )
			return last
		}
		last = ""
	}

	// check diagonal
	last = ""
	current = ""
	self.winningCells = nil
	for i := 0; i < 3; i++ {
		current = self.board[i][i];
		if ( current == "" ) {
		// empty slot means short circuit row check
			break
		}

		if ( last != "" && last != current ) {
			break
		}

		last = current
		}
		if ( last != "" && current != "" && last == current ) {
			self.winningCells = append( self.winningCells, makeKey(0,0) )
			self.winningCells = append( self.winningCells, makeKey(1,1) )
			self.winningCells = append( self.winningCells, makeKey(2,2) )
			return last
		}

	// check diagonal
	last = ""
	current = ""
	self.winningCells = nil
	for i := 2; i >= 0; i-- {
		x := i
		y := 2 - i
		current = self.board[x][y]
		if ( current == "" ) {
		// empty slot means short circuit row check
			break
		}

		if ( last != "" && last != current ) {
			break
		}

		last = current
	}

	if ( last != "" && current != "" && last == current ) {
		self.winningCells = append( self.winningCells, makeKey(2,0) )
		self.winningCells = append( self.winningCells, makeKey(1,1) )
		self.winningCells = append( self.winningCells, makeKey(0,2) )
		return last
	}

	self.winningCells = nil
	return ""
}

func makeKey(x, y int) string {
	return strconv.Itoa(x) + "_" + strconv.Itoa(y)
}



