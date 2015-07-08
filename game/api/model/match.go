package model

type Match struct {
	MatchID          string `json:"matchID"`
	ChallengerPlayer string `json:"challengerPlayer"`
	ChallengedPlayer string `json:"challengedPlayer"`
	MatchStartTs     int64 `json:"matchStartTs"`
	State       	 string `json:"state"`
}

func NewMatch(matchID, challengerPlayer, challengedPlayer string ) *Match {
	match := Match{}
	match.MatchID = matchID
	match.ChallengerPlayer = challengerPlayer
	match.ChallengedPlayer = challengedPlayer
	return &(match)
}

