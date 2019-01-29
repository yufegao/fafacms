package controllers

import (
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	resp := new(Resp)
	resp.Flag = true
	resp.Data = "FaFa CMS: https://github.com/hunterhug/fafacms"
	defer func() {
		c.JSON(200, resp)
	}()
}
