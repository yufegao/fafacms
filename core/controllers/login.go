package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/parrot/util"
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

	// paras not empty
	if req.UserName == "" || req.PassWd == "" {
		flog.Log.Errorf("login err:%s", "paras wrong")
		resp.Error = &ErrorResp{
			ErrorID:  ParasError,
			ErrorMsg: ErrorMap[ParasError],
		}
		return
	}
	// check session
	userInfo, _ := GetUserSession(c)
	if userInfo != nil {
		//c.Set("skipLog", true)
		c.Set("uid", userInfo.Id)
		resp.Flag = true
		return
	}

	// check cookie
	success, userInfo := CheckCookie(c)
	if success {
		err := SetUserSession(c, userInfo)
		if err != nil {
			flog.Log.Errorf("login err:%s", err.Error())
			resp.Error = &ErrorResp{
				ErrorID:  I500,
				ErrorMsg: ErrorMap[I500],
			}
			return
		}

		c.Set("uid", userInfo.Id)
		resp.Flag = true
		return
	}

	// super root user login
	if req.UserName == "hunterhug" && req.PassWd == "hunterhug" {
		u := new(model.User)
		u.Id = -1
		err := SetUserSession(c, u)
		if err != nil {
			flog.Log.Errorf("login err:%s", err.Error())
		}

		c.Set("uid", u.Id)
		resp.Flag = true
		return
	}

	// common people login
	uu := new(model.User)
	uu.Name = req.UserName
	uu.Password = req.PassWd
	ok, err := config.FafaRdb.Client.Get(uu)
	if err != nil {
		flog.Log.Errorf("login err:%s", err.Error())
		resp.Error = &ErrorResp{
			ErrorID:  I500,
			ErrorMsg: ErrorMap[I500],
		}
		return
	}

	if !ok {
		flog.Log.Errorf("login err:%s", "user or password wrong")
		resp.Error = &ErrorResp{
			ErrorID:  LoginWrong,
			ErrorMsg: ErrorMap[LoginWrong],
		}
		return
	}

	c.Set("uid", uu.Id)

	err = SetUserSession(c, uu)
	if err != nil {
		flog.Log.Errorf("login err:%s", err.Error())
		resp.Error = &ErrorResp{
			ErrorID:  I500,
			ErrorMsg: ErrorMap[I500],
		}
		return

	}

	resp.Flag = true
	
	if req.Remember {
		authKey := util.Md5(c.ClientIP() + "|" + uu.Password)
		secretKey := util.IS(uu.Id) + "|" + authKey
		c.SetCookie("auth", secretKey, 3600*24*7, "/", "", false, true)
	}
}

func Logout(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSON(c, 200, resp)
	}()
	user, _ := GetUserSession(c)
	if user != nil {
		DeleteUserSession(c)
	}
	c.SetCookie("auth", "", 0, "", "", false, false)
	resp.Flag = true
}
