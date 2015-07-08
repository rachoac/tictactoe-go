package model

import (
	"encoding/json"
)

type GameController struct {
	board *TicTacToeBoard
	turnOwner string
	playerOne string
	playerTwo string
}

type Data struct {
	PlayerOne string `json:"playerOne"`
	PlayerTwo string `json:"playerTwo"`
	Board []byte `json:"board"`
	TurnOwner string `json:"turnOwner"`
}

func (self *GameController) GetBoardData() ([]byte, error) {
	if ( self.board	== nil ) {
		self.board = NewBoard()
	}
	serialized, failure := self.board.Serialize()
	if ( failure != nil ) {
		return nil, failure
	}

	data := &Data{
		PlayerOne : self.playerOne,
		PlayerTwo: self.playerTwo,
		Board : serialized,
		TurnOwner : self.turnOwner }

	var toReturn, err = json.Marshal(data)

	if ( err != nil  ) {
		return nil, err
	}
	return toReturn, err
}

func (self *GameController) SetBoardData(jsonBytes []byte) {
	if ( self.board	== nil ) {
		self.board = NewBoard()
	}
	if ( jsonBytes == nil || len(jsonBytes) < 1 ) {
		return
	}

	var data Data
	json.Unmarshal(jsonBytes, &data)

	if ( data.Board != nil ) {
		self.board.Deserialize(data.Board)
	}

	self.playerOne = data.PlayerOne
	self.playerTwo = data.PlayerTwo
	self.turnOwner = data.TurnOwner
}

func (self *GameController) Create(playerOne, playerTwo string) {
	self.playerOne = playerOne
	self.playerTwo = playerTwo
	if ( self.board == nil ) {
		self.board = NewBoard()
	} else {
		self.board.Reset()
	}
}

func (self *GameController) GetTurnOwner() string {
	if ( self.turnOwner == "" ) {
		self.turnOwner = self.playerOne
	}

	return self.turnOwner
}

func (self *GameController) DoMove( x, y int ) error {
	failure := self.board.PerformMove( self.GetTurnOwner(), x, y )
	if ( failure != nil ) {
		return failure
	}
	nextTurn(self)
	return nil
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


func nextTurn(self *GameController) {
	if ( self.turnOwner == "" ) {
		self.turnOwner = self.playerOne
		return
	}

	if ( self.turnOwner == self.playerOne ) {
		self.turnOwner = self.playerTwo
	} else {
		self.turnOwner = self.playerOne
	}
}

