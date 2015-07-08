package api

import (
	"errors"
	log "github.com/Sirupsen/logrus"

	"github.com/rachoac/tictactoe-go/lobby/util"
	"github.com/rachoac/tictactoe-go/lobby/dao"
	"github.com/rachoac/tictactoe-go/lobby/api/model"
	"github.com/rachoac/tictactoe-go/lobby/api/gameclient"
//	"fmt"
)

type LobbyAPI struct {
	lobbyDAO *dao.LobbyDAO
}

/**
 * Join a named player to the lobby, indicating that they are ready for a game challenge
 * @param playerName
 * @return
 */
func (self *LobbyAPI) JoinPlayer(playerName string) *model.Player {
	player := self.lobbyDAO.GetPlayer(playerName)
	if ( player == nil ) {
		player = model.NewPlayer()
		player.PlayerName = playerName;
		player.JoinedTs = util.Now()
		self.lobbyDAO.SavePlayer(player)
	}

	return player
}

/**
 * Confirm whether or not the named player is in the lobby and is ready for a game challenge
 * @param playerName
 * @return
 */
func (self *LobbyAPI) IsPlayerJoined(playerName string) bool {
	player := self.lobbyDAO.GetPlayer(playerName)
	return player != nil
}

/**
 * Remove a named player from the lobby. This destroys any pending game challenges
 * @param playerName
 */
func (self *LobbyAPI) RemovePlayer(playerName string) {
	self.lobbyDAO.RemovePlayer(playerName)
}

/**
 * Creates a challenge between two players. If the challenge already exists, then an
 * exception is thrown
 * @param challengerPlayer
 * @param challengedPlayer
 * @return
 */
func (self *LobbyAPI) CreateChallenge(challengerPlayer, challengedPlayer string) (*model.Challenge, error) {
	failure := validatePlayer(self, challengerPlayer)
	if ( failure != nil ) {
		return nil, failure
	}

	failure = validatePlayer(self, challengedPlayer)
	if ( failure != nil ) {
		return nil, failure
	}

	challengeID := challengerPlayer + "_" + challengedPlayer
	challenge := self.lobbyDAO.GetChallenge(challengeID)
	if ( challenge != nil ) {
		// challenge already exists
		return nil, errors.New("Challenge already exists: " + challengeID)
	}

	challenge = model.NewChallenge()
	challenge.ChallengeID = challengeID
	challenge.ChallengedPlayer = challengedPlayer
	challenge.ChallengerPlayer = challengerPlayer
	challenge.ChallengeTs = util.Now()

	self.lobbyDAO.CreateChallenge(challenge)

	return challenge, nil
}

/**
 * Returns the first non expired game challenge offered to a player
 * @param challengedPlayer
 * @return
 */
func (self *LobbyAPI) GetChallengeFor(challengedPlayer string) (*model.Challenge, error) {
	failure := validatePlayer(self, challengedPlayer)
	if ( failure != nil ) {
		return nil, failure
	}

	challenges := self.lobbyDAO.GetChallengesFor( challengedPlayer )
	if ( len(challenges) < 1 ) {
		return nil, nil
	}

	for _, challenge := range challenges {

		if ( challenge == nil ) {
			continue
		}

		if ( challenge.IsChallengeExpired() ) {
			self.RemoveExpiredChallenge(challenge.ChallengeID)
			continue
		}

		return challenge, nil
	}

	return nil, nil
}

func (self *LobbyAPI) GetChallenge(challengeID string) *model.Challenge {
	return self.lobbyDAO.GetChallenge(challengeID)
}

/**
 * Removes the named challenge if its expired
 * @param challengeID
 * @return
 */
func (self *LobbyAPI) RemoveExpiredChallenge(challengeID string) bool {
	challenge := self.lobbyDAO.GetChallenge( challengeID )
	if ( challenge == nil ) {
		return false
	}

	if ( challenge.IsChallengeExpired() ) {
		log.Info("Challenge " + challengeID + " expired, removing")
		self.lobbyDAO.RemoveChallenge( challengeID )
	}

	return true
}

/**
 * Confirms that a challenge has been accepted. If it doesn't exist, returns
 * false.
 * @param challengeID
 * @return
 */
func (self *LobbyAPI) IsChallengeAccepted(challengeID string) bool {
	challenge := self.lobbyDAO.GetChallenge( challengeID )
	if ( challenge == nil ) {
		return false
	}

	return challenge.ChallengeStatus == "accepted"
}

/**
 * Indicates that the challenge has been accepted
 * @param challengeID
 * @return String matchID if successful match is created
 */
func (self *LobbyAPI) AcceptChallenge(challengeID string) (string, error) {
	challenge := self.lobbyDAO.GetChallenge( challengeID )
	if ( challenge == nil ) {
		return "", errors.New("Challenge " + challengeID + " not found")
	}

	challenge.ChallengeStatus = "accepted"

	self.lobbyDAO.SaveChallenge(challenge)

	client := gameclient.NewGameClient("http://localhost:9092");

	// call game service and prepare a match
	matchID, failure := client.CreateMatch( challenge.ChallengerPlayer, challenge.ChallengedPlayer );
	if ( failure != nil ) {
		return "", failure;
	}

	self.lobbyDAO.SetChallengeMatchID( challengeID, matchID )

	return matchID, nil
}

/**
 * Indicates that the challenge has been rejected
 * @param challengeID
 * @return
 */
func (self *LobbyAPI) RejectChallenge(challengeID string) bool {
	challenge := self.lobbyDAO.GetChallenge( challengeID )
	if ( challenge == nil ) {
		return true
	}

	challenge.ChallengeStatus = "rejected"
	self.lobbyDAO.SaveChallenge(challenge)
	return true
}

/**
 * Destroys the named challenge
 * @param challengeID
 */
func (self *LobbyAPI) RemoveChallenge(challengeID string) {
	self.lobbyDAO.RemoveChallenge(challengeID)
}

/**
 * Provides a list of players currently in the lobby and accepting games
 * @return
 */
func (self *LobbyAPI) GetJoinedPlayers() []*model.Player {
	var players []*model.Player = []*model.Player{}

	for _, playerID := range self.lobbyDAO.GetJoinedPlayers() {
		players = append( players, self.lobbyDAO.GetPlayer(playerID) )
	}
	return players
}

func (self *LobbyAPI) GetMatchIDForChallenge(challengeID string) string {
	return self.lobbyDAO.GetMatchIDForChallenge(challengeID)
}

// entrypoint
func NewLobbyAPI(_dao ... *dao.LobbyDAO) *LobbyAPI {
	apiObj := &(LobbyAPI{})
	if ( len(_dao) > 0 ) {
		apiObj.lobbyDAO = _dao[0]
	} else {
		apiObj.lobbyDAO = dao.NewLobbyDAO();
	}

	return apiObj
}

///////////////////////////////////////////////////////////////////////////////////////////
// private

// helpers
func validatePlayer(api *LobbyAPI, playerName string) error {
	player := api.lobbyDAO.GetPlayer(playerName)
	if (player == nil) {
		msg := "Challenged " + playerName + " not found"
		return errors.New(msg)
	}
	return nil
}
