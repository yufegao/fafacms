package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/util/mail"
	"github.com/hunterhug/parrot/util"
	"math"
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

	// send email
	mm := new(mail.Message)
	mm.Sender = config.FafaConfig.MailConfig
	mm.To = u.Email
	mm.ToName = u.NickName
	mm.Body = fmt.Sprintf(mm.Body, config.FafaConfig.Domain+"/verify?code="+u.ActivateMd5)
	err = mm.Sent()
	if err != nil {
		flog.Log.Errorf("RegisterUser err:%s", err.Error())
		resp.Error = Error(EmailError, err.Error())
		return
	}
	err = u.InsertOne()
	if err != nil {
		// db err
		flog.Log.Errorf("RegisterUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
	resp.Data = u
}

func ActivateUser(c *gin.Context) {
	resp := new(Resp)
	code := c.Query("code")
	defer func() {
		LogAlone(c, nil, resp)
	}()

	if code == "" {
		flog.Log.Errorf("ActivateUser err:%s", "code empty")
		resp.Error = Error(ParasError, "code empty")
		c.String(200, "code empty")
		return
	}

	u := new(model.User)
	u.ActivateMd5 = code

	exist, err := u.IsActivateCodeExist()
	if err != nil {
		flog.Log.Errorf("ActivateUser err:%s", err.Error())
		resp.Error = Error(ParasError, "db err")
		c.String(200, "db err")
		return
	}

	if !exist {
		flog.Log.Errorf("ActivateUser err:%s", "not exist code")
		resp.Error = Error(LazyError, "code not found")
		c.String(200, "code not found")
		return
	}

	if u.Status != 0 {
		c.Redirect(302, config.FafaConfig.Domain)
		return
	}

	if u.ActivateExpired < time.Now().Unix() {
		flog.Log.Errorf("ActivateUser err:%s", "code expired")
		resp.Error = Error(LazyError, "code expired")
		c.String(200, "code expired, resent email:<a href='%s/activate/code?code=%s'>Here</a>", config.FafaConfig.Domain, code)
		return
	} else {
		u.Status = 1
		err = u.UpdateStatus()
		if err != nil {
			flog.Log.Errorf("ActivateUser err:%s", err.Error())
			resp.Error = Error(ParasError, "db err")
			c.String(200, "db err")
			return
		}
		err = SetUserSession(c, u)
		if err != nil {
			flog.Log.Errorf("ActivateUser err:%s", err.Error())
			resp.Error = Error(I500, ErrorMap[I500])
			c.String(200, ErrorMap[I500])
			return
		}
	}

	c.Redirect(302, config.FafaConfig.Domain)

}

func ResendActivateCodeToUser(c *gin.Context) {
	resp := new(Resp)
	code := c.Query("code")
	defer func() {
		LogAlone(c, nil, resp)
	}()

	if code == "" {
		resp.Error = Error(ParasError, "code empty")
		c.String(200, "code empty")
		return
	}

	u := new(model.User)
	u.ActivateMd5 = code

	exist, err := u.IsActivateCodeExist()
	if err != nil {
		flog.Log.Errorf("ResendUser err:%s", err.Error())
		resp.Error = Error(ParasError, "db err")
		c.String(200, "db err")
		return
	}
	if !exist {
		flog.Log.Errorf("ResendUser err:%s", "not exist code")
		resp.Error = Error(LazyError, "code not found")
		c.String(200, "code not found")
		return
	}

	if u.Status != 0 {
	} else if u.ActivateExpired > time.Now().Unix() {
		flog.Log.Errorf("ResendUser err:%s", "code not expired")
		resp.Error = Error(LazyError, "code not expired")
		c.String(200, "code not expired")
		return
	}

	err = u.UpdateActivateCode()
	if err != nil {
		flog.Log.Errorf("ResendUser err:%s", err.Error())
		resp.Error = Error(ParasError, "db err")
		c.String(200, "db err")
		return
	}

	// send email
	mm := new(mail.Message)
	mm.Sender = config.FafaConfig.MailConfig
	mm.To = u.Email
	mm.ToName = u.NickName
	mm.Body = fmt.Sprintf(mm.Body, config.FafaConfig.Domain+"/activate?code="+u.ActivateMd5)
	err = mm.Sent()
	if err != nil {
		flog.Log.Errorf("ResendUser err:%s", err.Error())
		resp.Error = Error(EmailError, err.Error())
		return
	}

	c.String(200, "email code reset")
}

type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func ForgetPasswordOfUser(c *gin.Context) {
	resp := new(Resp)
	req := new(ForgetPasswordRequest)
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
	u.Email = req.Email
	ok, err := u.GetUserByEmail()
	if err != nil {
		flog.Log.Errorf("ForgetPassword err:%s", err.Error())
		resp.Error = Error(DBError, ErrorMap[DBError])
		return
	}
	if !ok {
		flog.Log.Errorf("ForgetPassword err:%s", "not found")
		resp.Error = Error(DbNotFound, ErrorMap[DbNotFound])
		return
	}

	if u.CodeExpired < time.Now().Unix() {
		err = u.UpdateCode()
		if err != nil {
			flog.Log.Errorf("ForgetPassword err:%s", err.Error())
			resp.Error = Error(DBError, ErrorMap[DBError])
			return
		}

		// send email
		mm := new(mail.Message)
		mm.Sender = config.FafaConfig.MailConfig
		mm.To = u.Email
		mm.ToName = u.NickName
		mm.Body = "code is: " + u.Code
		err = mm.Sent()
		if err != nil {
			flog.Log.Errorf("ForgetPassword err:%s", err.Error())
			resp.Error = Error(EmailError, err.Error())
			return
		}

	} else {
		flog.Log.Errorf("ForgetPassword err:%s", "time not reach")
		resp.Error = Error(TimeNotReachError, ErrorMap[TimeNotReachError])
		return
	}

	resp.Flag = true
}

type ChangePasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Code     string `json:"code" validate:"required,lt=9,gt=5"`
	Password string `json:"password" validate:"alphanumunicode,gt=5,lt=17"`
}

func ChangePasswordOfUser(c *gin.Context) {
	resp := new(Resp)
	req := new(ChangePasswordRequest)
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
		flog.Log.Errorf("ChangePassword err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	u := new(model.User)
	u.Email = req.Email
	ok, err := u.GetUserByEmail()
	if err != nil {
		flog.Log.Errorf("ChangePassword err:%s", err.Error())
		resp.Error = Error(DBError, ErrorMap[DBError])
		return
	}
	if !ok {
		flog.Log.Errorf("ChangePassword err:%s", "not found")
		resp.Error = Error(DbNotFound, "email")
		return
	}

	if u.Code == req.Code {
		u.Password = req.Password
		err = u.UpdatePassword()
		if err != nil {
			flog.Log.Errorf("ChangePassword err:%s", err.Error())
			resp.Error = Error(DBError, ErrorMap[DBError])
			return
		}
	} else {
		flog.Log.Errorf("ChangePassword err:%s", "code wrong")
		resp.Error = Error(CodeWrong, "not valid")
		return
	}

	DeleteUserSession(c)
	c.SetCookie("auth", "", 0, "", "", false, false)
	resp.Flag = true
}

type UpdateUserRequest struct {
	NickName  string `json:"nick_name" validate:"omitempty,gt=1,lt=50"`
	WeChat    string `json:"wechat" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	WeiBo     string `json:"weibo" validate:"omitempty,url"`
	Github    string `json:"github" validate:"omitempty,url"`
	QQ        string `json:"qq" validate:"omitempty,numeric,gt=6,lt=12"`
	Gender    int    `json:"gender" validate:"oneof=0 1 2"`
	Describe  string `json:"describe" validate:"omitempty,lt=200"`
	ImagePath string `json:"image_path" validate:"omitempty,lt=100"`
}

func UpdateUser(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateUserRequest)
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
		flog.Log.Errorf("UpdateUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateUser err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}

	if uu == nil {
		flog.Log.Errorf("UpdateUser err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}

	u := new(model.User)
	u.Id = uu.Id

	// if image not empty
	if req.ImagePath != "" {
		p := new(model.Picture)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			// db err
			flog.Log.Errorf("UpdateUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			// not found
			flog.Log.Errorf("UpdateUser err: image not exist")
			resp.Error = Error(ParasError, "image url not exist")
			return
		}

		u.HeadPhoto = req.ImagePath
	}

	u.Describe = req.Describe
	u.NickName = req.NickName
	u.Gender = req.Gender
	u.WeChat = req.WeChat
	u.QQ = req.QQ
	u.Github = req.Github
	u.WeiBo = req.WeiBo
	err = u.UpdateInfo()
	if err != nil {
		// db err
		flog.Log.Errorf("UpdateUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
	resp.Data = u

}

func TakeUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()

	u, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("TakeUser err:%s", "session not found")
		resp.Error = Error(LazyError, "session not found")
		return
	}
	resp.Flag = true
	resp.Data = u
}

type ListUserRequest struct {
	Id              int      `json:"id"`
	Name            string   `json:"name" validate:"lt=100"`
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	UpdateTimeBegin int64    `json:"update_time_begin"`
	UpdateTimeEnd   int64    `json:"update_time_end"`
	Sort            []string `json:"sort" validate:"dive,lt=100"`

	Email  string `json:"email" validate:"omitempty,email"`
	WeChat string `json:"wechat" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	WeiBo  string `json:"weibo" validate:"omitempty,url"`
	Github string `json:"github" validate:"omitempty,url"`
	QQ     string `json:"qq" validate:"omitempty,numeric,gt=6,lt=12"`
	Gender int    `json:"gender" validate:"oneof=-1 0 1 2"`
	Status int    `json:"status" validate:"oneof=-1 0 1 2"`

	PageHelp
}

type ListUserResponse struct {
	Users []model.User `json:"users"`
	PageHelp
}

func ListUser(c *gin.Context) {
	resp := new(Resp)

	respResult := new(ListUserResponse)
	req := new(ListUserRequest)
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
		flog.Log.Errorf("ListUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.User)).Where("1=1")

	// query prepare
	if req.Id != 0 {
		session.And("id=?", req.Id)
	}
	if req.Name != "" {
		session.And("name=?", req.Name)
	}

	if req.Status != -1 {
		session.And("status=?", req.Status)
	}

	if req.Gender != -1 {
		session.And("gender=?", req.Gender)
	}

	if req.QQ != "" {
		session.And("q_q=?", req.QQ)
	}

	if req.Email != "" {
		session.And("email=?", req.Email)
	}

	if req.Github != "" {
		session.And("github=?", req.Github)
	}

	if req.WeiBo != "" {
		session.And("wei_bo=?", req.WeiBo)
	}
	if req.WeChat != "" {
		session.And("we_chat=?", req.WeChat)
	}

	if req.CreateTimeBegin > 0 {
		session.And("create_time>=?", req.CreateTimeBegin)
	}

	if req.CreateTimeEnd > 0 {
		session.And("create_time<?", req.CreateTimeBegin)
	}

	if req.UpdateTimeBegin > 0 {
		session.And("update_time>=?", req.UpdateTimeBegin)
	}

	if req.UpdateTimeEnd > 0 {
		session.And("update_time<?", req.UpdateTimeEnd)
	}

	// count num
	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		// db err
		flog.Log.Errorf("ListUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if count>0 start list
	users := make([]model.User, 0)
	p := &req.PageHelp
	if total == 0 {
	} else {
		// sql build
		p.build(session, req.Sort, model.UserSortName)
		// do query
		err = session.Find(&users)
		if err != nil {
			// db err
			flog.Log.Errorf("ListUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	// result
	respResult.Users = users
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

type AssignGroupRequest struct {
	GroupId      int   `json:"group_id"`
	GroupRelease int   `json:"group_release"`
	Users        []int `json:"users"`
}

func AssignGroupToUser(c *gin.Context) {
	resp := new(Resp)
	req := new(AssignGroupRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if len(req.Users) == 0 {
		flog.Log.Errorf("AssignGroupToUser err:%s", "users empty")
		resp.Error = Error(ParasError, "users")
		return
	}

	if req.GroupRelease == 1 {
		u := new(model.User)
		num, err := config.FafaRdb.Client.Table(new(model.User)).Cols("group_id").In("id", req.Users).Update(u)
		if err != nil {
			flog.Log.Errorf("AssignGroupToUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
		resp.Data = num
	} else {
		if req.GroupId == 0 {
			flog.Log.Errorf("AssignGroupToUser err:%s", "group id empty")
			resp.Error = Error(ParasError, "group_id")
			return
		}

		g := new(model.Group)
		g.Id = req.GroupId
		exist, err := g.GetById()
		if err != nil {
			flog.Log.Errorf("AssignGroupToUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !exist {
			flog.Log.Errorf("AssignGroupToUser err:%s", "group not found")
			resp.Error = Error(DbNotFound, "group")
			return
		}

		u := new(model.User)
		u.GroupId = req.GroupId
		num, err := config.FafaRdb.Client.Table(new(model.User)).Cols("group_id").In("id", req.Users).Update(u)
		if err != nil {
			flog.Log.Errorf("AssignGroupToUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
		resp.Data = num
	}

	resp.Flag = true

}
