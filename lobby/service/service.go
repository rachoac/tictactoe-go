package service

import (
	api "github.com/rachoac/tictactoe-go/lobby/api"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"fmt"
)
func Init() {
	setupRoutes()
}

type Error struct {
	err string `json:"err"`
}

// setup the routes
func setupRoutes() {
	r := gin.Default()
	lobby  := api.NewLobbyAPI()

	// middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		c.Next()
	})

	log.Info("Setting up routes")

	r.OPTIONS("/lobby/roster", func(c *gin.Context) {
		c.Header("Allow", "GET,PUT,POST,DELETE,OPTIONS")
		c.String(200, "");
	})

	// Return list of players in the lobby
	r.GET("/lobby/roster", func(c *gin.Context) {
		players := lobby.GetJoinedPlayers()
		c.JSON(200, players)
	})

	// Joins a player to the lobby, making them avaiable for game challenges
	// query param: playerName
	r.POST("/lobby/roster", func(c *gin.Context) {
		playerName := c.Query("playerName")
		player := lobby.JoinPlayer(playerName)
		c.JSON(200, player )
	})

	// Removes a player from the lobby, making them unavailable for game challenges
	// query param: playerName
	r.DELETE("/lobby/roster", func(c *gin.Context) {
		playerName := c.Query("playerName")
		lobby.RemovePlayer(playerName)
	})

	// Challenge a player
	// query param: challengerPlayerName
	r.POST("/lobby/roster/:challengedPlayerName/challenge", func(c *gin.Context) {
		challengedPlayerName := c.Param("challengedPlayerName")
		challengerPlayerName := c.Query("challengerPlayerName")
		challenge, failure := lobby.CreateChallenge(challengerPlayerName, challengedPlayerName)
		if ( failure != nil ) {
			fmt.Println(failure)
			c.JSON(500, gin.H{"error" : failure.Error() })
		} else {
			c.JSON(200, challenge)
		}
	})

	// Retrieves a player challenge (first nonexpired challenge)
	r.GET("/lobby/roster/:challengedPlayerName/challenge", func(c *gin.Context) {
		challengedPlayerName := c.Param("challengedPlayerName")
		challenge, failure := lobby.GetChallengeFor(challengedPlayerName)
		if ( failure != nil ) {
			fmt.Println(failure)
			c.JSON(500, gin.H{"error" : failure.Error() })
		} else {
			if ( challenge == nil ) {
				c.JSON(404, gin.H{"error" : "could not find challenge for player " + challengedPlayerName })
			} else {
				c.JSON(200, challenge)
			}
		}
	})

	r.OPTIONS("/lobby/challenge/:challengedPlayerName", func(c *gin.Context) {
		c.Header("Allow", "GET,PUT,POST,DELETE,OPTIONS")
		c.String(200, "");
	})

	// Retrieves the status of a challenge
	r.GET("/lobby/challenge/:challengeID", func(c *gin.Context) {
		challengeID := c.Param("challengeID")
		challenge := lobby.GetChallenge(challengeID)

		if ( challenge == nil ) {
			c.JSON(500, gin.H{"error" : "ChallengeID " + challengeID + " not found"})
		} else {
			c.JSON(200, challenge)
		}
	})

	// Cancels a challenge
	r.DELETE("/lobby/challenge/:challengeID", func(c *gin.Context) {
		challengeID := c.Param("challengeID")
		challenge := lobby.GetChallenge(challengeID)

		if ( challenge == nil ) {
			c.JSON(500, gin.H{"error" : "ChallengeID " + challengeID + " not found"})
		} else {
			lobby.RemoveChallenge(challengeID)
			c.JSON(200, challenge)
		}
	})

	//
	r.GET("/lobby/challenge/:challengeID/matchID", func(c *gin.Context) {
		challengeID := c.Param("challengeID")

		c.String(200, lobby.GetMatchIDForChallenge(challengeID))
	})

	// Accepts or rejects player challenge
	r.PUT("/lobby/challenge/:challengeID", func(c *gin.Context) {
		challengeID := c.Param("challengeID")
		challenge := lobby.GetChallenge(challengeID)
		if ( challenge == nil ) {
			c.JSON(500, gin.H{"error" : "Invalid challenge"})
			return
		}

		response := c.Query("response")
		if ( "accepted" == response ) {
			matchID, failure := lobby.AcceptChallenge(challengeID)
			challenge.MatchID = matchID
			if ( failure != nil ) {
				c.JSON(500, gin.H{"error" : failure.Error() })
			} else {
				c.JSON(200, challenge )
			}
		} else if ("rejected"  == response) {
			lobby.RejectChallenge(challengeID)
			challenge.MatchID = ""
			c.JSON(200, challenge )
		} else {
			c.JSON(500, gin.H{"error" : "Invalid challenge response " + response })
		}
	})

	r.Run(":9090") // listen and serve on 0.0.0.0:9090
}
