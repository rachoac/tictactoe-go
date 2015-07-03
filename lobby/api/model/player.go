package model

type Player struct {
	PlayerName string `json:"playerName" binding:"required"`
	JoinedTs int64 `json:"joinedTs" binding:"required"`
}

func NewPlayer() *Player {
	instance := Player{}
	return &instance
}