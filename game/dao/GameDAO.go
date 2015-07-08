package dao

import (
	"github.com/rachoac/tictactoe-go/game/api/model"
)

type GameDAO struct {
	matches map[string]*model.Match
	matchBoards map[string][]byte
}

func (self *GameDAO) GetMatch(matchID string) *model.Match {
	return self.matches[matchID]
}

func (self *GameDAO) CreateMatch(match *model.Match) {
	self.matches[match.MatchID] = match
}

func (self *GameDAO) UpdateMatch(match *model.Match) {
	self.CreateMatch(match)
}

func (self* GameDAO) SaveGame( matchID string, boardData []byte ) {
	self.matchBoards[matchID] = boardData;
}

func (self *GameDAO) GetGame(matchID string) []byte {
	return self.matchBoards[matchID];
}

// entry point
func NewGameDAO() *GameDAO {
	daoObj := GameDAO{}
	daoObj.matches = make(map[string]*model.Match )
	daoObj.matchBoards = make(map[string][]byte )
	return &daoObj;
}
