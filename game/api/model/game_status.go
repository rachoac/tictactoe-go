package model

import "encoding/json"
import "encoding/base64"

type GameStatus struct {
	State string `json:"state"`
	Winner string `json:"winner"`
	TurnOwner string `json:"turnOwner"`
	WinningCells []string `json:"winningCells"`
	Match *Match `json:"match"`
	BoardDataMap map[string]string `json:"boardData"`
}

func setBoardDataMap(self* GameStatus, boardData []byte) map[string]string {
	var m map[string]string

	json.Unmarshal(boardData,  &m)
	b, err := base64.StdEncoding.DecodeString(m["board"])

	if ( err != nil ) {}
//	fmt.Println( "XXX: ", b, err )

	var x map[string]string
	json.Unmarshal(b,  &x)

	return x
}

func NewGameStatus( state, winner, turnOwner string, match *Match, boardData []byte) *GameStatus {
	obj := GameStatus{}
	obj.State = state
	obj.Winner = winner
	obj.TurnOwner = turnOwner
	obj.Match = match
	obj.BoardDataMap = setBoardDataMap(&obj, boardData)

	return &(obj)
}