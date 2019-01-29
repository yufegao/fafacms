package controllers

import (
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	UserName string `json:"user_name"`
	PassWd   string `json:"pass_wd"`
	Remember string `json:"remember"`
}

func Login(c *gin.Context) {
	resp := new(Resp)
	req := new(LoginRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}
}
