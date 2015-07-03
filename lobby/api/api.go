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

type LobbyAPI struct {}

/**
 * Join a named player to the lobby, indicating that they are ready for a game challenge
 * @param playerName
 * @return
 */
func (api *LobbyAPI) JoinPlayer(playerName string) *model.Player {
	player := lobbyDAO.GetPlayer(playerName)
	if ( player == nil ) {
		player = model.NewPlayer()
		player.PlayerName = playerName;
		player.JoinedTs = util.Now()
		lobbyDAO.SavePlayer(player)
	}

	return player
}

/**
 * Confirm whether or not the named player is in the lobby and is ready for a game challenge
 * @param playerName
 * @return
 */
func (api *LobbyAPI) IsPlayerJoined(playerName string) bool {
	player := lobbyDAO.GetPlayer(playerName)
	return player != nil
}

/**
 * Remove a named player from the lobby. This destroys any pending game challenges
 * @param playerName
 */
func (api *LobbyAPI) RemovePlayer(playerName string) {
	lobbyDAO.RemovePlayer(playerName)
}

/**
 * Creates a challenge between two players. If the challenge already exists, then an
 * exception is thrown
 * @param challengerPlayer
 * @param challengedPlayer
 * @return
 */
func (api *LobbyAPI) CreateChallenge(challengerPlayer, challengedPlayer string) (*model.Challenge, error) {
	failure := validatePlayer(challengerPlayer)
	if ( failure != nil ) {
		return nil, failure
	}

	failure = validatePlayer(challengedPlayer)
	if ( failure != nil ) {
		return nil, failure
	}

	challengeID := challengerPlayer + "_" + challengedPlayer
	challenge := lobbyDAO.GetChallenge(challengeID)
	if ( challenge != nil ) {
		// challenge already exists
		return nil, errors.New("Challenge already exists: " + challengeID)
	}

	challenge = model.NewChallenge()
	challenge.ChallengeID = challengeID
	challenge.ChallengedPlayer = challengedPlayer
	challenge.ChallengerPlayer = challengerPlayer
	challenge.ChallengeTs = util.Now()

	lobbyDAO.CreateChallenge(challenge)

	return challenge, nil
}

/**
 * Returns the first non expired game challenge offered to a player
 * @param challengedPlayer
 * @return
 */
func (api *LobbyAPI) GetChallengeFor(challengedPlayer string) (*model.Challenge, error) {
	failure := validatePlayer(challengedPlayer)
	if ( failure != nil ) {
		return nil, failure
	}

	challenges := lobbyDAO.GetChallengesFor( challengedPlayer )
	if ( len(challenges) < 1 ) {
		return nil, nil
	}

	for _, challenge := range challenges {

		if ( challenge == nil ) {
			continue
		}

		if ( challenge.IsChallengeExpired() ) {
			apiObj.RemoveExpiredChallenge(challenge.ChallengeID)
			continue
		}

		return challenge, nil
	}

	return nil, nil
}

func (api *LobbyAPI) GetChallenge(challengeID string) *model.Challenge {
	return lobbyDAO.GetChallenge(challengeID)
}

/**
 * Removes the named challenge if its expired
 * @param challengeID
 * @return
 */
func (api *LobbyAPI) RemoveExpiredChallenge(challengeID string) bool {
	challenge := lobbyDAO.GetChallenge( challengeID )
	if ( challenge == nil ) {
		return false
	}

	if ( challenge.IsChallengeExpired() ) {
		log.Info("Challenge " + challengeID + " expired, removing")
		lobbyDAO.RemoveChallenge( challengeID )
	}

	return true
}

/**
 * Confirms that a challenge has been accepted. If it doesn't exist, returns
 * false.
 * @param challengeID
 * @return
 */
func (api *LobbyAPI) IsChallengeAccepted(challengeID string) bool {
	challenge := lobbyDAO.GetChallenge( challengeID )
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
func (api *LobbyAPI) AcceptChallenge(challengeID string) (string, error) {
	challenge := lobbyDAO.GetChallenge( challengeID )
	if ( challenge == nil ) {
		return "", errors.New("Challenge " + challengeID + " not found")
	}

	challenge.ChallengeStatus = "accepted"

	lobbyDAO.SaveChallenge(challenge)

	client := gameclient.NewGameClient("http://localhost:9092");

	// call game service and prepare a match
	matchID, failure := client.CreateMatch( challenge.ChallengerPlayer, challenge.ChallengedPlayer );
	if ( failure != nil ) {
		return "", failure;
	}

	lobbyDAO.SetChallengeMatchID( challengeID, matchID )

	return matchID, nil
}

/**
 * Indicates that the challenge has been rejected
 * @param challengeID
 * @return
 */
func (api *LobbyAPI) RejectChallenge(challengeID string) bool {
	challenge := lobbyDAO.GetChallenge( challengeID )
	if ( challenge == nil ) {
		return true
	}

	challenge.ChallengeStatus = "rejected"
	lobbyDAO.SaveChallenge(challenge)
	return true
}

/**
 * Destroys the named challenge
 * @param challengeID
 */
func (api *LobbyAPI) RemoveChallenge(challengeID string) {
	lobbyDAO.RemoveChallenge(challengeID)
}

/**
 * Provides a list of players currently in the lobby and accepting games
 * @return
 */
func (api *LobbyAPI) GetJoinedPlayers() []*model.Player {
	var players []*model.Player = []*model.Player{}

	for _, playerID := range lobbyDAO.GetJoinedPlayers() {
		players = append( players, lobbyDAO.GetPlayer(playerID) )
	}
	return players
}

func (api *LobbyAPI) GetMatchIDForChallenge(challengeID string) string {
	return lobbyDAO.GetMatchIDForChallenge(challengeID)
}

// entrypoint
func NewLobbyAPI(_dao ... *dao.LobbyDAO) *LobbyAPI {
	apiObj = &(LobbyAPI{})
	if ( len(_dao) > 0 ) {
		lobbyDAO = _dao[0]
	} else {
		lobbyDAO = dao.NewLobbyDAO();
	}

	return apiObj
}

///////////////////////////////////////////////////////////////////////////////////////////
// private

// state
var apiObj *LobbyAPI
var lobbyDAO *dao.LobbyDAO

// helpers
func validatePlayer(playerName string) error {
	player := lobbyDAO.GetPlayer(playerName)
	if (player == nil) {
		msg := "Challenged " + playerName + " not found"
		return errors.New(msg)
	}
	return nil
}
