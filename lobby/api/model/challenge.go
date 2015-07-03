package model

import "github.com/rachoac/tictactoe-go/lobby/util"

type Challenge struct {
	MatchID string `json:"matchID"`
	ChallengeStatus string `json:"challengeStatus"`
	ChallengeID string `json:"challengeID" binding:"required"`
	ChallengedPlayer string `json:"challengedPlayer" binding:"required"`
	ChallengerPlayer string `json:"challengerPlayer" binding:"required"`
	ChallengeTs int64 `json:"challengeTs" binding:"required"`
}

func (self *Challenge) IsChallengeExpired() bool {
	return util.Elapsed(self.ChallengeTs, 10000);
}

func NewChallenge() *Challenge {
	return &(Challenge{})
}
