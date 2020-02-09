package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taiko-web/api"
	"taiko-web/config"
)

func NewRouter(conf *config.TaikoWebConfig) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	router.LoadHTMLGlob("templates/*.tmpl")
	router.Static("/src", "./static/src")
	router.Static("/assets", "./static/assets")
	router.Static("/songs", "./static/songs")

	// Default View
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Version": conf.Web.Version,
			"Config":  conf.Web,
		})
	})

	// API
	router.GET("/api/config", api.ConfigHandler)
	router.GET("/api/songs", api.SongsHandler)
	router.GET("/p2", api.MultiplayerHandler)

	return router
}

func Init(conf *config.TaikoWebConfig) {
	gin.SetMode(conf.Mode)
	r := NewRouter(conf)
	r.Run(":" + conf.Port)
}
