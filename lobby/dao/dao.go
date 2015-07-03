package dao

import (
	"github.com/rachoac/tictactoe-go/lobby/api/model"
)


///////////////////////////////////////////////////////////////////////////////////////////
// public

type LobbyDAO struct {
}


func (dao *LobbyDAO) GetPlayer(playerName string) *model.Player {
	return joinedPlayers[playerName]
}

func (dao *LobbyDAO) SavePlayer(player *model.Player) {
	joinedPlayers[player.PlayerName] = player
	var challenges []string
	playerChallenges[player.PlayerName] = challenges
}

func (dao *LobbyDAO) RemovePlayer(playerName string) {
	delete(joinedPlayers, playerName)
	delete(playerChallenges, playerName)
}

func (dao *LobbyDAO) GetChallenge(challengeID string) *model.Challenge {
	return challenges[challengeID]
}

func (dao *LobbyDAO) CreateChallenge(challenge *model.Challenge) {
	challenges[challenge.ChallengeID] = challenge

	challengeIDs := playerChallenges[challenge.ChallengedPlayer]
	if ( challengeIDs == nil ) {
		challengeIDs = []string{}
	}

	challengeIDs = append( challengeIDs, challenge.ChallengeID )
	playerChallenges[challenge.ChallengedPlayer] = challengeIDs
}

func (dao *LobbyDAO) RemoveChallenge(challengeID string) {
	// remove challenge
	delete(challenges, challengeID)

	// destroy player challenge
	for k, v := range(playerChallenges) {
		playerChallenges[k] = filter(v, func(value string) bool {
			return challengeID != value
		})
	}
}

func (dao *LobbyDAO) GetChallengesFor(challengedPlayer string) []*model.Challenge {
	var thisPlayerChallenges []*model.Challenge
	challengeIDs := playerChallenges[challengedPlayer]
	if ( challengeIDs == nil ) {
		return thisPlayerChallenges
	}

	thisPlayerChallenges = _mapChallenges(challengeIDs, func(val string) *model.Challenge {
		return challenges[val]
	})

	return thisPlayerChallenges
}

func (dao *LobbyDAO) SaveChallenge(challenge *model.Challenge) {
	challenges[challenge.ChallengeID] = challenge
}

func (dao *LobbyDAO) GetJoinedPlayers() []string {
	keys := []string{}
	for k := range joinedPlayers {
		keys = append(keys, k)
	}

	return keys
}

func (dao *LobbyDAO) SetChallengeMatchID(challengeID, matchID string) {
	challenge := dao.GetChallenge(challengeID)
	if ( challenge != nil ) {
		challenge.MatchID = matchID
	}
}

func (dao *LobbyDAO) GetMatchIDForChallenge(challengeID string) string {
	challenge := dao.GetChallenge(challengeID)
	if ( challenge == nil ) {
		return ""
	}
	return challenge.MatchID
}

// entry point
func NewLobbyDAO() *LobbyDAO {
	joinedPlayers = make(map[string]*model.Player)
	playerChallenges = make(map[string][]string)
	challenges = make(map[string]*model.Challenge)
	daoObj = &(LobbyDAO{})
	return daoObj
}

///////////////////////////////////////////////////////////////////////////////////////////
// private

// state
var daoObj *LobbyDAO
var joinedPlayers map[string]*model.Player
var playerChallenges map[string][]string
var challenges map[string]*model.Challenge

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

