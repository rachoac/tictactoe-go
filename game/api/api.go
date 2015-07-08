package api

import (
	"github.com/rachoac/tictactoe-go/game/dao"
	"github.com/rachoac/tictactoe-go/game/util"
	"github.com/rachoac/tictactoe-go/game/api/model"
	log "github.com/Sirupsen/logrus"
	"errors"
	"fmt"
)

type GameAPI struct {
	gameDAO *dao.GameDAO
}

func (self *GameAPI) CreateMatch(challengerPlayer, challengedPlayer string ) *model.Match {
	matchID := challengerPlayer + "_" + challengedPlayer
	match := model.NewMatch( matchID, challengerPlayer, challengedPlayer )
	match.MatchStartTs = util.Now()

	log.Info(challengerPlayer + " challenged " + challengedPlayer)

	self.gameDAO.CreateMatch(match);

	controller := model.NewGameController()
	controller.Create(challengerPlayer, challengedPlayer)

	data, failure := controller.GetBoardData()

	if ( failure == nil ) {
		self.gameDAO.SaveGame(match.MatchID, data)
	} else {
		// xxx TODO!
		fmt.Println("BAAAAD!", failure)
	}

	return match
}

func (self *GameAPI) GetMatchTurnOwner(matchID string) string {
	controller := getGameController(self, matchID);
	return controller.GetTurnOwner()
}

func (self *GameAPI) PerformMove(matchID, player string, x, y int) error {
	controller := getGameController(self, matchID)
	if ( controller.GetTurnOwner() == player ) {
		controller.DoMove(x, y);
		data, failure := controller.GetBoardData()
		if ( failure == nil ) {
			self.gameDAO.SaveGame(matchID, data)
		} else {
			return failure
		}
	} else {
		return errors.New(matchID + ": Not player " + player + "'s turn, its " + controller.GetTurnOwner() )
	}

	if ( controller.GetState() == "won" || controller.GetState() == "stalemate" ) {
		// end of the game
		// todo
		// inform lobby that the game is complete
		// ...
	}

	return nil
}

func (self *GameAPI) GetMatchStatus( matchID string ) *model.GameStatus {
	match := self.gameDAO.GetMatch(matchID);
	controller := getGameController(self, matchID)
	state := controller.GetState()
	winner := controller.GetWinner()
	turnOwner := controller.GetTurnOwner()
	winingCells := controller.GetWinningCells()

	data, failure := controller.GetBoardData()
	if ( failure == nil ) {
		status := model.NewGameStatus(state, winner, turnOwner, match, data)
		status.WinningCells = winingCells
		return status
	} else {
		return nil
	}
}

func (self *GameAPI) StopMatch(matchID string) *model.GameStatus {
	match := self.gameDAO.GetMatch(matchID)
	match.State = "stopped"
	self.gameDAO.UpdateMatch(match)

	return self.GetMatchStatus(matchID)
}

func (self *GameAPI) GetGameController(matchID string) *model.GameController {
	gameData := self.gameDAO.GetGame(matchID)
	controller := model.NewGameController()
	controller.SetBoardData(gameData)
	return controller
}

// entrypoint
func NewGameAPI(_dao ... *dao.GameDAO) *GameAPI {
	apiObj := GameAPI{}
	if ( len(_dao) > 0 ) {
		apiObj.gameDAO = _dao[0]
	} else {
		apiObj.gameDAO = dao.NewGameDAO()
	}
	return &apiObj
}

// private
func getGameController(api *GameAPI, matchID string) *model.GameController {
	gameData := api.gameDAO.GetGame( matchID );
	controller := model.NewGameController();
	controller.SetBoardData(gameData);
	return controller;
}
