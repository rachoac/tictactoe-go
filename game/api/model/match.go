package model

type Match struct {
	MatchID          string
	ChallengerPlayer string
	ChallengedPlayer string
	MatchStartTs     int64
	MatchState       string
}

