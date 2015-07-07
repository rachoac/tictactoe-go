package dao

import (
	"github.com/rachoac/tictactoe-go/game/api/model"
	"encoding/json"
)

type GameDAO struct {
}

func (dao *GameDAO) GetMatch(matchID string) *model.Match {
	return
}

func (dao *GameDAO) CreateMatch(matchID string) {
}

func (dao *GameDAO) UpdateMatch(match *model.Match) {
}

func (dao* GameDAO) SaveGame( matchID string )

// entry point
func NewGameDAO() *GameDAO {
	return
}

///////////////////////////////////////////////////////////////////////////////////////////
// private

// state
var daoObj *GameDAO