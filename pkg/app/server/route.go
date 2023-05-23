package server

import (
	userDlv "nostr-ex/pkg/app/user/delivery"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(Logger, gin.Recovery())

	router.POST("/event", userDlv.PostEvent)
	router.POST("/event/req", userDlv.ReqEvent)
	router.DELETE("/event/req", userDlv.CloseReq)

	router.GET("/exit", exit)
	router.GET("/info", info)

	return router
}

func info(c *gin.Context) {
	c.JSON(200, gin.H{ // response json
		"version": "0.0.0.1",
	})
}

func exit(c *gin.Context) { //TODO:
	c.JSON(200, gin.H{ // response json
		"version": "0.0.0.1",
	})
}
