package api

import (
	"github.com/rachoac/tictactoe-go/game/api/model"
	"encoding/json"
)

type GameController struct {
	board *model.TicTacToeBoard
	turnOwner string
	playerOne string
	playerTwo string
}

type Data struct {
	playerOne string
	playerTwo string
	board []byte
	turnOwner string
}

func (self *GameController) GetBoardData() []byte {
	if ( self.board	== nil ) {
		self.board = model.NewBoard()
	}
	data := Data{
		self.playerOne, self.playerTwo, self.board.Serialize(), self.turnOwner }
	return json.Marshal(data)
}

func (self *GameController) SetBoardData(jsonBytes []byte) {
	if ( self.board	== nil ) {
		self.board = model.NewBoard()
	}
	if ( jsonBytes == nil || len(jsonBytes) < 1 ) {
		return
	}

	var data Data
	json.Unmarshal(jsonBytes, &data)

	if ( data.board != nil ) {
		self.board.Deserialize(data.board)
	}

	self.playerOne = data.playerOne
	self.playerTwo = data.playerTwo
	self.turnOwner = data.turnOwner
}

func (self *GameController) Create(playerOne, playerTwo string) {
	self.playerOne = playerOne
	self.playerTwo = playerTwo
	if ( self.board == nil ) {
		self.board = model.NewBoard()
	} else {
		self.board.Reset()
	}
}

func (self *GameController) GetTurnOwner() {
	if ( self.turnOwner == nil ) {
		self.turnOwner = self.playerOne
	}

	return self.turnOwner
}

func (self *GameController) DoMove( x, y int ) error {
	return self.board.PerformMove( self.GetTurnOwner(), x, y )
}

func (self *GameController) GetState() string {
	return self.board.GetState()
}

func (self *GameController) GetWinner() string {
	return self.board.GetWinner()
}

func (self *GameController) GetWinningCells() []string {
	return self.board.GetWinningCells()
}

//
func NewGameController() *GameController {
	gameController := GameController{}
	return &gameController
}

///////////////////////////////////////////////////////////////////////////////////////////
// private


func (self *GameController) nextTurn() {
	if ( self.turnOwner == nil ) {
		self.turnOwner = self.playerOne
		return
	}

	if ( self.turnOwner == self.playerOne ) {
		self.turnOwner = self.playerTwo
	} else {
		self.turnOwner = self.playerOne
	}
}

