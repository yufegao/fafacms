package controllers

import (
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func UpdateUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func DeleteUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func TakeUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func ListUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}