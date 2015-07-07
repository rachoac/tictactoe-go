package service

import (
//	api "github.com/rachoac/tictactoe-go/game/api"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
//	"fmt"
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
//	game  := api.NewGameAPI()

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
		c.JSON(200, "")
	})

	//
	r.GET("/game/match/:matchID/turnOwner", func(c *gin.Context) {
		c.JSON(200, "")
	})

	// @QueryParam("player")
	// @QueryParam("x")
	// @QueryParam("y")
	r.PUT("/game/match/:matchID/move", func(c *gin.Context) {
		c.JSON(200, "")
	})

	r.GET("/game/match/:matchID", func(c *gin.Context) {
		c.JSON(200, "")
	})

	r.DELETE("/game/match/:matchID", func(c *gin.Context) {
		c.JSON(200, "")
	})

	r.Run(":9092") // listen and serve on 0.0.0.0:9092
}
