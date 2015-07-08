package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateMatch(t *testing.T) {
	api := NewGameAPI()

	challenger := "challenger"
	challenged := "challenged"

	match := api.CreateMatch(challenger, challenged)
	assert.NotNil(t, match)

	assert.NotEmpty(t, match.MatchID)
	assert.Equal(t, challenger, match.ChallengerPlayer)
	assert.Equal(t, challenged, match.ChallengedPlayer)
	assert.Equal(t, "", match.State)
	assert.NotEmpty(t, match.MatchStartTs)
}

func TestGetMatchTurnOwner(t *testing.T) {
	api := NewGameAPI()

	challenger := "challenger"
	challenged := "challenged"

	match := api.CreateMatch(challenger, challenged)

	turnOwner := api.GetMatchTurnOwner(match.MatchID)
	assert.Equal(t, challenger, turnOwner)

	api.PerformMove(match.MatchID, turnOwner, 0, 0)

	turnOwner = api.GetMatchTurnOwner(match.MatchID)
	assert.Equal(t, challenged, turnOwner)
}

func TestGetMatchStatusBasic(t *testing.T) {
	api := NewGameAPI()

	challenger := "challenger"
	challenged := "challenged"

	match := api.CreateMatch(challenger, challenged)

	status := api.GetMatchStatus(match.MatchID)
	assert.NotNil(t, status)
	assert.NotNil(t, status.BoardDataMap  )
	assert.Empty(t, status.BoardDataMap["0_0"])

	assert.Equal( t, challenger, status.TurnOwner )
	assert.Equal( t, status.TurnOwner, api.GetMatchTurnOwner(match.MatchID))
	assert.Equal( t, status.Match, match)
	assert.Equal( t, "running", status.State )

	err := api.PerformMove(match.MatchID, challenger, 0, 0)
	assert.Empty(t, err)
	status = api.GetMatchStatus(match.MatchID)
	assert.NotNil(t, status.BoardDataMap  )
	assert.Equal(t, challenger, status.BoardDataMap["0_0"])
}

func TestGetMatchStatusWin(t *testing.T) {
	api := NewGameAPI()

	challenger := "challenger"
	challenged := "challenged"

	match := api.CreateMatch(challenger, challenged)

	api.PerformMove(match.MatchID, challenger, 0, 0)
			api.PerformMove(match.MatchID, challenged, 1, 0)
	status := api.GetMatchStatus(match.MatchID)
	assert.Empty(t, status.WinningCells)
	assert.Empty(t, status.Winner)
	assert.Equal(t, "running", status.State)

	api.PerformMove(match.MatchID, challenger, 0, 1)
			api.PerformMove(match.MatchID, challenged, 1, 1)
	status = api.GetMatchStatus(match.MatchID)
	assert.Empty(t, status.WinningCells)
	assert.Empty(t, status.Winner)
	assert.Equal(t, "running", status.State)

	api.PerformMove(match.MatchID, challenger, 0, 2)
	status = api.GetMatchStatus(match.MatchID)
	assert.NotEmpty(t, status.WinningCells)
	assert.Equal(t, []string{ "0_0", "0_1", "0_2"}, status.WinningCells )
	assert.Equal(t, challenger, status.Winner)
	assert.Equal(t, "won", status.State)
}