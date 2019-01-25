package controllers

import (
	. "github.com/hunterhug/fafacms/core/config"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	resp := &Resp{Flag: false}
	defer func() {
		c.JSON(200, resp)
	}()

}
