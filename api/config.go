package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taiko-web/config"
)

func ConfigHandler(c *gin.Context) {
	conf := config.GetConfig()
	c.JSON(http.StatusOK, &conf.Web)
}
