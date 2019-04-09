package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/parrot/util"
	"time"
)

type RegisterUserRequest struct {
	Name      string `json:"name" validate:"required,alphanumunicode,gt=1,lt=50"`
	NickName  string `json:"nick_name" validate:"required,gt=1,lt=50"`
	Email     string `json:"email" validate:"required,email"`
	WeChat    string `json:"wechat" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	WeiBo     string `json:"weibo" validate:"omitempty,url"`
	Github    string `json:"github" validate:"omitempty,url"`
	QQ        string `json:"qq" validate:"omitempty,numeric,gt=6,lt=12"`
	Password  string `json:"password" validate:"alphanumunicode,gt=5,lt=17"`
	Gender    int    `json:"gender" validate:"oneof=0 1 2"`
	Describe  string `json:"describe" validate:"omitempty,lt=200"`
	ImagePath string `json:"image_path" validate:"omitempty,lt=100"`
}

func RegisterUser(c *gin.Context) {
	resp := new(Resp)
	req := new(RegisterUserRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	// validate
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("RegisterUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	u := new(model.User)

	// name check
	u.Name = req.Name
	repeat, err := u.IsNameRepeat()
	if err != nil {
		flog.Log.Errorf("RegisterUser err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if repeat {
		flog.Log.Errorf("RegisterUser err: %s", "name already use by other")
		resp.Error = Error(ParasError, "name already use by other")
		return
	}

	// email check
	u.Email = req.Email
	repeat, err = u.IsEmailRepeat()
	if err != nil {
		flog.Log.Errorf("RegisterUser err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if repeat {
		flog.Log.Errorf("RegisterUser err: %s", "email already use by other")
		resp.Error = Error(ParasError, "email already use by other")
		return
	}

	// if image not empty
	if req.ImagePath != "" {
		p := new(model.Picture)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			// db err
			flog.Log.Errorf("RegisterUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			// not found
			flog.Log.Errorf("RegisterUser err: image not exist")
			resp.Error = Error(ParasError, "image url not exist")
			return
		}

		u.HeadPhoto = req.ImagePath
	}

	u.ActivateMd5 = util.Md5(u.Email)
	u.CreateTime = time.Now().Unix()
	u.Describe = req.Describe
	u.ActivateExpired = time.Now().Add(48 * time.Hour).Unix()
	u.NickName = req.NickName
	u.Password = req.Password
	u.Gender = req.Gender
	u.WeChat = req.WeChat
	u.QQ = req.QQ
	u.Github = req.Github
	u.WeiBo = req.WeiBo

	err = u.InsertOne()
	if err != nil {
		// db err
		flog.Log.Errorf("RegisterUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
}

func VerifyUser(c *gin.Context) {
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
