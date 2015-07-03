package gameclient

import (
    "net/http"
    "net/url"
    "io/ioutil"

    "github.com/rachoac/tictactoe-go/lobby/api/model"
    "encoding/json"
)

type GameClient struct {
}

func (client *GameClient) CreateMatch( challengerPlayer, challengedPlayer string ) (string, error) {
    resp, err := http.PostForm(gameServiceURL + "/game/match/create",
        url.Values{"challengerPlayer" : {challengerPlayer}, "challengedPlayer" : {challengedPlayer}})

    if err != nil {
        return "", err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        return "", err
    }

    var match model.Match
    err = json.Unmarshal(body, &match)
    if err != nil {
        return "", err
    }

    return match.MatchID, nil
}

var gameServiceURL string
func NewGameClient(_gameServiceURL string) *GameClient {
    gameServiceURL = _gameServiceURL;
    return &GameClient{}
}