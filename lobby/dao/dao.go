package dao

import (
	"github.com/rachoac/tictactoe-go/lobby/api/model"
)


///////////////////////////////////////////////////////////////////////////////////////////
// public

type LobbyDAO struct {
	daoObj *LobbyDAO
	joinedPlayers map[string]*model.Player
	playerChallenges map[string][]string
	challenges map[string]*model.Challenge
}

func (self *LobbyDAO) GetPlayer(playerName string) *model.Player {
	return self.joinedPlayers[playerName]
}

func (self *LobbyDAO) SavePlayer(player *model.Player) {
	self.joinedPlayers[player.PlayerName] = player
	var challenges []string
	self.playerChallenges[player.PlayerName] = challenges
}

func (self *LobbyDAO) RemovePlayer(playerName string) {
	delete(self.joinedPlayers, playerName)
	delete(self.playerChallenges, playerName)
}

func (self *LobbyDAO) GetChallenge(challengeID string) *model.Challenge {
	return self.challenges[challengeID]
}

func (self *LobbyDAO) CreateChallenge(challenge *model.Challenge) {
	self.challenges[challenge.ChallengeID] = challenge

	challengeIDs := self.playerChallenges[challenge.ChallengedPlayer]
	if ( challengeIDs == nil ) {
		challengeIDs = []string{}
	}

	challengeIDs = append( challengeIDs, challenge.ChallengeID )
	self.playerChallenges[challenge.ChallengedPlayer] = challengeIDs
}

func (self *LobbyDAO) RemoveChallenge(challengeID string) {
	// remove challenge
	delete(self.challenges, challengeID)

	// destroy player challenge
	for k, v := range(self.playerChallenges) {
		self.playerChallenges[k] = filter(v, func(value string) bool {
			return challengeID != value
		})
	}
}

func (self *LobbyDAO) GetChallengesFor(challengedPlayer string) []*model.Challenge {
	var thisPlayerChallenges []*model.Challenge
	challengeIDs := self.playerChallenges[challengedPlayer]
	if ( challengeIDs == nil ) {
		return thisPlayerChallenges
	}

	thisPlayerChallenges = _mapChallenges(challengeIDs, func(val string) *model.Challenge {
		return self.challenges[val]
	})

	return thisPlayerChallenges
}

func (self *LobbyDAO) SaveChallenge(challenge *model.Challenge) {
	self.challenges[challenge.ChallengeID] = challenge
}

func (self *LobbyDAO) GetJoinedPlayers() []string {
	keys := []string{}
	for k := range self.joinedPlayers {
		keys = append(keys, k)
	}

	return keys
}

func (self *LobbyDAO) SetChallengeMatchID(challengeID, matchID string) {
	challenge := self.GetChallenge(challengeID)
	if ( challenge != nil ) {
		challenge.MatchID = matchID
	}
}

func (self *LobbyDAO) GetMatchIDForChallenge(challengeID string) string {
	challenge := self.GetChallenge(challengeID)
	if ( challenge == nil ) {
		return ""
	}
	return challenge.MatchID
}

// entry point
func NewLobbyDAO() *LobbyDAO {
	daoObj := &(LobbyDAO{})

	daoObj.joinedPlayers = make(map[string]*model.Player)
	daoObj.playerChallenges = make(map[string][]string)
	daoObj.challenges = make(map[string]*model.Challenge)
	return daoObj
}

///////////////////////////////////////////////////////////////////////////////////////////
// private

// state


// helpers
func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func _mapChallenges(vs []string, f func(string) *model.Challenge) []*model.Challenge {
	vsm := make([]*model.Challenge, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

