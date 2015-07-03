package api

import "testing"
import "github.com/rachoac/tictactoe-go/lobby/util"
import (
	"github.com/stretchr/testify/assert"
	"github.com/rachoac/tictactoe-go/lobby/api/model"
	"github.com/rachoac/tictactoe-go/lobby/dao"
	"fmt"
)

func TestJoinPlayer(t *testing.T ) {
	api := NewLobbyAPI();

	playerName := util.UUID();
	player := api.JoinPlayer(playerName);

	assert.NotNil(t, player, "Joined player should not be null");
}

func TestRemovePlayer(t *testing.T ) {
	api := NewLobbyAPI();

	playerName := util.UUID();
	api.JoinPlayer(playerName);
	assert.True(t, api.IsPlayerJoined(playerName));

	api.RemovePlayer(playerName);
	assert.False(t, api.IsPlayerJoined(playerName));
}

func TestCreateChallenge(t *testing.T) {
	api := NewLobbyAPI();

	challenger := util.UUID();
	challenged := util.UUID();

	api.JoinPlayer(challenger);
	api.JoinPlayer(challenged);

	challenge, failure := api.CreateChallenge(challenger, challenged);
	assert.NotNil(t, challenge);
	assert.Nil(t, failure, "Did not expect an error on first challenge");

	challenge, failure = api.CreateChallenge(challenger, challenged);
	assert.NotNil(t, failure, "Rechallenge should have failed");
}

func TestGetChallenge(t *testing.T) {
	api := NewLobbyAPI();

	challenger := util.UUID();
	challenged := util.UUID();

	api.JoinPlayer(challenger);
	api.JoinPlayer(challenged);

	challenge, failure := api.CreateChallenge(challenger, challenged);
	assert.Nil(t, failure);
	assert.NotNil(t, challenge);
	assert.NotNil(t, challenge.ChallengeID);
	fmt.Println(challenge)

	resolvedChallenge, failure := api.GetChallengeFor(challenged);
	assert.Nil(t, failure);
	assert.NotNil(t, resolvedChallenge);

	assert.NotNil(t, resolvedChallenge.ChallengeID);
	assert.Equal(t, challenge.ChallengeID, resolvedChallenge.ChallengeID );
}

func TestGetExpiredChallenge(t *testing.T) {
	dao := dao.NewLobbyDAO();
	api := NewLobbyAPI(dao);

	challenger := util.UUID();
	challenged := util.UUID();

	api.JoinPlayer(challenger);
	api.JoinPlayer(challenged);

	challenge := model.NewChallenge()
	challenge.ChallengeID = util.UUID()
	challenge.ChallengerPlayer = challenger
	challenge.ChallengedPlayer = challenged
	challenge.ChallengeTs = 100;

	dao.CreateChallenge( challenge );

	resolvedChallenge, failure := api.GetChallengeFor(challenged);
	assert.Nil(t, failure);
	assert.Nil(t, resolvedChallenge);
}

func TestGetFirstNonExpiredChallenge(t *testing.T) {
	dao := dao.NewLobbyDAO();
	api := NewLobbyAPI(dao);

	challenger := util.UUID();
	challenged := util.UUID();

	api.JoinPlayer(challenger);
	api.JoinPlayer(challenged);

	var firstOne *model.Challenge
	{
		challenge := model.NewChallenge()
		challenge.ChallengeID = util.UUID()
		challenge.ChallengerPlayer = challenger
		challenge.ChallengedPlayer = challenged
		challenge.ChallengeTs = 100;

		dao.CreateChallenge( challenge );

	}
	{
		challenge := model.NewChallenge()
		challenge.ChallengeID = util.UUID()
		challenge.ChallengerPlayer = challenger
		challenge.ChallengedPlayer = challenged
		challenge.ChallengeTs = util.Now()

		dao.CreateChallenge( challenge );
		firstOne = challenge;
	}
	{
		challenge := model.NewChallenge()
		challenge.ChallengeID = util.UUID()
		challenge.ChallengerPlayer = challenger
		challenge.ChallengedPlayer = challenged
		challenge.ChallengeTs = util.Now()

		dao.CreateChallenge( challenge );

	}
	{
		challenge := model.NewChallenge()
		challenge.ChallengeID = util.UUID()
		challenge.ChallengerPlayer = challenger
		challenge.ChallengedPlayer = challenged
		challenge.ChallengeTs = 100;

		dao.CreateChallenge( challenge );
	}

	resolvedChallenge, failure := api.GetChallengeFor(challenged);

	assert.NotNil(t, resolvedChallenge);
	assert.Nil(t, failure);
	assert.Equal(t, firstOne.ChallengeID, resolvedChallenge.ChallengeID)
}

func TestAcceptChallenge(t *testing.T) {
	api := NewLobbyAPI();

	challenger := util.UUID();
	challenged := util.UUID();

	api.JoinPlayer(challenger);
	api.JoinPlayer(challenged);

	challenge, failure := api.CreateChallenge(challenger, challenged);
	assert.False(t, api.IsChallengeAccepted(challenge.ChallengeID));
	assert.Nil(t, failure);

	api.AcceptChallenge(challenge.ChallengeID);

	assert.True(t, api.IsChallengeAccepted(challenge.ChallengeID));
}

func TestRejectChallenge(t *testing.T) {
	api := NewLobbyAPI();

	challenger := util.UUID();
	challenged := util.UUID();

	api.JoinPlayer(challenger);
	api.JoinPlayer(challenged);

	challenge, failure := api.CreateChallenge(challenger, challenged);
	assert.False(t, api.IsChallengeAccepted(challenge.ChallengeID));
	assert.Nil(t, failure);

	api.RejectChallenge(challenge.ChallengeID);

	assert.False(t, api.IsChallengeAccepted(challenge.ChallengeID));
}

