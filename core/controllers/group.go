package controllers

import (
	"github.com/gin-gonic/gin"
)

func CreateGroup(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func UpdateGroup(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func DeleteGroup(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func TakeGroup(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func ListGroup(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}
