package controllers

import (
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		c.JSON(200, resp)
	}()
}

func UpdateUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		c.JSON(200, resp)
	}()
}

func DeleteUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		c.JSON(200, resp)
	}()
}

func TakeUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		c.JSON(200, resp)
	}()
}
