package server

import (
	"net/http"
	relayDlv "nostr-ex/pkg/app/relay/delivery"
	userDlv "nostr-ex/pkg/app/user/delivery"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(Logger, gin.Recovery())

	router.LoadHTMLGlob("web/templates/*")

	//router.GET("/", SocketHandler)
	router.GET("/watcher", func(c *gin.Context) {
		c.HTML(http.StatusOK, "watcher.html", gin.H{
			"title": "Posts",
		})
	})
	router.POST("/relay", relayDlv.AddRelay)

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
