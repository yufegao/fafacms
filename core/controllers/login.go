package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
)

type LoginRequest struct {
	UserName string `json:"user_name"`
	PassWd   string `json:"pass_wd"`
	Remember bool   `json:"remember"`
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

	userInfo, _ := GetUserSeesion(c)
	if userInfo != nil {
		c.Set("uid", userInfo.Id)
		resp.Flag = true
		return
	}

	// super user
	if req.UserName == "hunterhug" && req.PassWd == "fafa" {
		u := new(model.User)
		u.Id = -1
		err := SetUserSession(c, u)
		if err != nil {
			flog.Log.Errorf("login err:%s", err.Error())
			resp.Error = &ErrorResp{
				ErrorID:  LoginPermit,
				ErrorMsg: fmt.Sprintf("%s:%s", ErrorMap[LoginPermit], err.Error()),
			}
			return
		}

		c.Set("uid", u.Id)
		resp.Flag = true
		return
	}

	return
}
