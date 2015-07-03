package dao

import "testing"
import "fmt"
import "github.com/rachoac/tictactoe-go/lobby/api/model"
import "github.com/stretchr/testify/assert"

func TestGetSavePlayer(t *testing.T ) {
	dao := NewLobbyDAO();

	player := model.NewPlayer()
	player.PlayerName = "aron"
	player.JoinedTs = 1234

	dao.SavePlayer(player)
	savedPlayer := dao.GetPlayer(player.PlayerName)

	assert.NotEmpty(t, savedPlayer, "The original player should be rehydratable")
	assert.Equal(t, player.PlayerName, savedPlayer.PlayerName, "The rehydrated player name should be the same as the originally saved player")
	assert.Equal(t, player.JoinedTs, savedPlayer.JoinedTs, "The rehydrated player joined ts should be the same as the originally saved player")
}

func TestRemovePlayer(t *testing.T ) {
	dao := NewLobbyDAO();

	player := model.NewPlayer()
	player.PlayerName = "aron"
	player.JoinedTs = 1234

	dao.SavePlayer(player)
	dao.RemovePlayer(player.PlayerName)
	savedPlayer := dao.GetPlayer(player.PlayerName)

	assert.Nil(t, savedPlayer, "The original player should be destroyed")
}

func TestCreateSaveChallenge(t *testing.T) {
	dao := NewLobbyDAO();

	challenge := model.NewChallenge()
	challenge.ChallengeID = "abc123"

	dao.CreateChallenge(challenge)

	savedChallenge := dao.GetChallenge(challenge.ChallengeID)
	assert.NotNil(t, savedChallenge, "The original challenge should be rehydratable")

	assert.Equal(t, savedChallenge.ChallengeID, savedChallenge.ChallengeID, "The rehydrated challenge should have the same ID")
}

func TestRemoveChallenge(t *testing.T) {
	dao := NewLobbyDAO();

	challenge := model.NewChallenge()
	challenge.ChallengeID = "abc123"

	dao.CreateChallenge(challenge)
	dao.RemoveChallenge(challenge.ChallengeID)

	savedChallenge := dao.GetChallenge(challenge.ChallengeID)

	assert.Nil(t, savedChallenge, "The original challenge should be destroyed")
}

func TestGetPlayerChallenges(t *testing.T) {
	dao := NewLobbyDAO();

	playerName := "aron"

	{
		challenge := model.NewChallenge()
		challenge.ChallengeID = "abc123"
		challenge.ChallengedPlayer = playerName
		dao.CreateChallenge(challenge)
	}

	challenges := dao.GetChallengesFor(playerName)
	fmt.Println(challenges)
	assert.True(t, len(challenges) == 1, "Should be 1 challenge")

	{
		challenge := model.NewChallenge()
		challenge.ChallengeID = "abc456"
		challenge.ChallengedPlayer = playerName
		dao.CreateChallenge(challenge)
	}
	challenges = dao.GetChallengesFor(playerName)
	assert.True(t, len(challenges) == 2, "Should be 2 challenges")
	fmt.Println(len(challenges), challenges)
}

func TestRemovePlayerChallenges(t *testing.T) {
	dao := NewLobbyDAO();

	playerName := "aron"

	{
		challenge := model.NewChallenge()
		challenge.ChallengeID = "abc123"
		challenge.ChallengedPlayer = playerName
		dao.CreateChallenge(challenge)
	}
	{
		challenge := model.NewChallenge()
		challenge.ChallengeID = "abc456"
		challenge.ChallengedPlayer = playerName
		dao.CreateChallenge(challenge)
	}
	challenges := dao.GetChallengesFor(playerName)
	assert.True(t, len(challenges) == 2, "Should be 2 challenges")
	fmt.Println(len(challenges), challenges)

	dao.RemoveChallenge("abc123")

	challenges = dao.GetChallengesFor(playerName)
	fmt.Println(challenges)
	assert.True(t, len(challenges) == 1, "Should be 1 challenge")
}

func TestGetJoinedPlayers(t *testing.T) {
	dao := NewLobbyDAO();

	joinedPlayers := dao.GetJoinedPlayers();
	assert.True(t, len(joinedPlayers) == 0, "Should be 0 players")

	for i := 0; i < 3; i++ {
		{
			player := model.NewPlayer()
			player.PlayerName = "aron" + string(i)
			player.JoinedTs = 1234

			dao.SavePlayer(player)
		}
	}

	joinedPlayers = dao.GetJoinedPlayers();
	assert.True(t, len(joinedPlayers) == 3, "Should be 3 players")

	for i := 0; i < 3; i++ {
		{
			dao.RemovePlayer("aron" + string(i))
		}
	}
	joinedPlayers = dao.GetJoinedPlayers();
	assert.True(t, len(joinedPlayers) == 0, "Should be 0 players")
}

func TestGetSetChallengeMatchID(t *testing.T) {
	dao := NewLobbyDAO();

	challenge := model.NewChallenge()
	challenge.ChallengeID = "abc123"
	challenge.ChallengedPlayer = "aron"
	dao.CreateChallenge(challenge)

	dao.SetChallengeMatchID(challenge.ChallengeID, "match123");
	persistedMatchID := dao.GetMatchIDForChallenge(challenge.ChallengeID)
	fmt.Println(persistedMatchID)
	assert.Equal(t, "match123", persistedMatchID, "Failed to find matchID for challenge")
}
