package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/config"
)

func Home(c *gin.Context) {
	resp := new(config.Resp)
	resp.Flag = true
	resp.Data = "FaFa CMS: https://github.com/hunterhug/fafacms"
	defer func() {
		c.JSON(200, resp)
	}()
}
