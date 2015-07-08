package service

import (
	"github.com/rachoac/tictactoe-go/game/api"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/rachoac/tictactoe-go/game/dao"
	"errors"
	"strconv"
	"fmt"
)
func Init() {
	setupRoutes()
}

type Error struct {
	err string `json:"err"`
}

type MatchForm struct {
	ChallengerPlayer string `form:"challengerPlayer" binding:"required"`
	ChallengedPlayer string	`form:"challengedPlayer" binding:"required"`
}

// setup the routes
func setupRoutes() {
	r := gin.Default()
	game  := api.NewGameAPI(dao.NewGameDAO())

	// middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		c.Next()
	})

	log.Info("Setting up routes")

	// @FormParam("challengerPlayer")
	// @FormParam("challengedPlayer")
	r.POST("/game/match/create", func(c *gin.Context) {
		var form MatchForm
		var result = c.Bind(&form)
		fmt.Println(c)
		if result == nil {
			match := game.CreateMatch(form.ChallengerPlayer, form.ChallengedPlayer)
			c.JSON(200, match)
		} else {
			fmt.Println(result)
			c.JSON(500, gin.H{"error" : errors.New("Error parsing game create form post") })
		}
	})

	//
	r.GET("/game/match/:matchID/turnOwner", func(c *gin.Context) {
		matchID := c.Query("matchID")
		c.JSON(200, game.GetMatchTurnOwner(matchID))
	})

	r.OPTIONS("/game/match/:matchID/move", func(c *gin.Context) {
		c.Header("Allow", "GET,PUT,POST,DELETE,OPTIONS")
		c.String(200, "");
	})


	r.OPTIONS("/game/match/:matchID", func(c *gin.Context) {
		c.Header("Allow", "GET,PUT,POST,DELETE,OPTIONS")
		c.String(200, "");
	})

	// @QueryParam("player")
	// @QueryParam("x")
	// @QueryParam("y")
	r.PUT("/game/match/:matchID/move", func(c *gin.Context) {
		matchID := c.Param("matchID")
		player := c.Query("player")
		var failure error
		x, failure := strconv.ParseInt(c.Query("x"), 10, 32)
		if ( failure != nil ) {
			fmt.Println(failure)
			c.JSON(500, gin.H{"error" : failure.Error() })
		}

		y, failure := strconv.ParseInt(c.Query("y"), 10, 32)
		if ( failure != nil ) {
			fmt.Println(failure)
			c.JSON(500, gin.H{"error" : failure.Error() })
		}

		fmt.Println(matchID, player, x, y)

		failure = game.PerformMove(matchID, player, int(x), int(y));
		if ( failure == nil ) {
			c.JSON(200, game.GetMatchStatus(matchID))
		} else {
			fmt.Println(failure)
			c.JSON(500, gin.H{"error" : failure.Error() })
		}
	})

	r.GET("/game/match/:matchID", func(c *gin.Context) {
		matchID := c.Param("matchID")
		c.JSON(200, game.GetMatchStatus(matchID))
	})

	r.DELETE("/game/match/:matchID", func(c *gin.Context) {
		matchID := c.Param("matchID")
		gameStatus := game.StopMatch(matchID)
		c.JSON(200, gameStatus)
	})

	r.Run(":9092") // listen and serve on 0.0.0.0:9092
}
