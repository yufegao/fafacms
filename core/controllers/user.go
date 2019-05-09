package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/util"
	"github.com/hunterhug/fafacms/core/util/mail"
	"math"
	"time"
)

type RegisterUserRequest struct {
	Name       string `json:"name" validate:"required,alphanumunicode,gt=1,lt=50"`
	NickName   string `json:"nick_name" validate:"required,gt=1,lt=50"`
	Email      string `json:"email" validate:"required,email"`
	WeChat     string `json:"wechat" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	WeiBo      string `json:"weibo" validate:"omitempty,url"`
	Github     string `json:"github" validate:"omitempty,url"`
	QQ         string `json:"qq" validate:"omitempty,numeric,gt=6,lt=12"`
	Password   string `json:"password" validate:"alphanumunicode,gt=5,lt=17"`
	RePassword string `json:"repassword" validate:"eqfield=Password"`
	Gender     int    `json:"gender" validate:"oneof=0 1 2"`
	Describe   string `json:"describe" validate:"omitempty,lt=200"`
	ImagePath  string `json:"image_path" validate:"omitempty,lt=100"`
}

// 用户注册，任何人可以用唯一邮箱来注册
func RegisterUser(c *gin.Context) {
	resp := new(Resp)
	req := new(RegisterUserRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	// 配置如果关闭注册，那么直接返回
	if config.FafaConfig.DefaultConfig.CloseRegister {
		resp.Error = Error(LazyError, "can not register")
		return
	}

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("RegisterUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// 唯一名字不能重复，作为子域名存在
	u := new(model.User)
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

	// 邮箱不能重复
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
		p := new(model.File)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {

			flog.Log.Errorf("RegisterUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			flog.Log.Errorf("RegisterUser err: image not exist")
			resp.Error = Error(ParasError, "image url not exist")
			return
		}

		u.HeadPhoto = req.ImagePath
	}

	// 激活验证码
	u.ActivateCode = util.GetGUID()
	u.ActivateCodeExpired = time.Now().Add(48 * time.Hour).Unix()

	u.Describe = req.Describe
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
	mm.Body = fmt.Sprintf(mm.Body, u.ActivateCode)
	err = mm.Sent()
	if err != nil {
		flog.Log.Errorf("RegisterUser err:%s", err.Error())
		resp.Error = Error(EmailError, err.Error())
		return
	}

	err = u.InsertOne()
	if err != nil {
		flog.Log.Errorf("RegisterUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true

	// 如果不是调试模式，不应该返回信息
	if AuthDebug {
		resp.Data = u
	}

}

// 创建用户，管理员权限
func CreateUser(c *gin.Context) {
	resp := new(Resp)
	req := new(RegisterUserRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("CreateUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	u := new(model.User)
	u.Name = req.Name
	repeat, err := u.IsNameRepeat()
	if err != nil {
		flog.Log.Errorf("CreateUser err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if repeat {
		flog.Log.Errorf("CreateUser err: %s", "name already use by other")
		resp.Error = Error(ParasError, "name already use by other")
		return
	}

	// email check
	u.Email = req.Email
	repeat, err = u.IsEmailRepeat()
	if err != nil {
		flog.Log.Errorf("CreateUser err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if repeat {
		flog.Log.Errorf("CreateUser err: %s", "email already use by other")
		resp.Error = Error(ParasError, "email already use by other")
		return
	}

	// if image not empty
	if req.ImagePath != "" {
		p := new(model.File)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			flog.Log.Errorf("CreateUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			flog.Log.Errorf("CreateUser err: image not exist")
			resp.Error = Error(ParasError, "image url not exist")
			return
		}

		u.HeadPhoto = req.ImagePath
	}

	u.Describe = req.Describe
	u.NickName = req.NickName
	u.Password = req.Password
	u.Gender = req.Gender
	u.WeChat = req.WeChat
	u.QQ = req.QQ
	u.Github = req.Github
	u.WeiBo = req.WeiBo

	// 默认激活
	u.Status = 1

	err = u.InsertOne()
	if err != nil {
		flog.Log.Errorf("CreateUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
	resp.Data = u
}

// 用户自己激活自己
func ActivateUser(c *gin.Context) {
	resp := new(Resp)
	code := c.Query("code")
	email := c.Query("email")
	defer func() {
		LogAlone(c, nil, resp)
	}()

	if code == "" {
		flog.Log.Errorf("ActivateUser err:%s", "code empty")
		resp.Error = Error(ParasError, "code empty")
		c.String(200, "code empty")
		return
	}

	// 必须邮箱和激活码一起来
	u := new(model.User)
	u.ActivateCode = code
	u.Email = email

	// 判断激活码是否存在
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

	// 如果用户不是未激活状态
	if u.Status != 0 {
		flog.Log.Errorf("ActivateUser err:%s", "already active")
		resp.Error = Error(LazyError, "already active")
		c.String(200, "already active")
		return
	}

	// 验证码过期，要重新生成验证码，需要用户手动请求另外的API
	if u.ActivateCodeExpired < time.Now().Unix() {
		flog.Log.Errorf("ActivateUser err:%s", "code expired")
		resp.Error = Error(LazyError, "code expired")
		c.String(200, "code expired, resent email:<a href='%s/activate/resent?code=%s&email=%s'>Here</a>", config.FafaConfig.Domain, code, email)
		return
	} else {
		// 更新用户的状态
		u.Status = 1
		err = u.UpdateStatus()
		if err != nil {
			flog.Log.Errorf("ActivateUser err:%s", err.Error())
			resp.Error = Error(ParasError, "db err")
			c.String(200, "db err")
			return
		}

		// 激活成功马上为用户设置Session
		err = SetUserSession(c, u)
		if err != nil {
			flog.Log.Errorf("ActivateUser err:%s", err.Error())
			resp.Error = Error(I500, "")
			c.String(200, ErrorMap[I500])
			return
		}
	}

	resp.Flag = true
	c.String(200, "ok")
}

// 用户激活验证码失效了，重新生成并发送邮件
func ResendActivateCodeToUser(c *gin.Context) {
	resp := new(Resp)
	code := c.Query("code")
	email := c.Query("email")
	defer func() {
		LogAlone(c, nil, resp)
	}()

	if code == "" {
		resp.Error = Error(ParasError, "code empty")
		c.String(200, "code empty")
		return
	}

	u := new(model.User)
	u.ActivateCode = code
	u.Email = email

	// 要生成新的验证码必须携带之前的验证码才行
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
		// 用户不是未激活状态
		flog.Log.Errorf("ActivateUser err:%s", "already active")
		resp.Error = Error(LazyError, "already active")
		c.String(200, "already active")
		return
	} else if u.ActivateCodeExpired > time.Now().Unix() {
		// 验证码过期时间还没到，要等一下
		flog.Log.Errorf("ResendUser err:%s", "code not expired")
		resp.Error = Error(LazyError, "code not expired")
		c.String(200, "code not expired")
		return
	}

	// 更新验证码，过期时间48小时
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
	mm.Body = fmt.Sprintf(mm.Body, u.ActivateCode)
	err = mm.Sent()
	if err != nil {
		flog.Log.Errorf("ResendUser err:%s", err.Error())
		resp.Error = Error(EmailError, err.Error())
		return
	}

	resp.Flag = true
	c.String(200, "email code reset")
}

type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// 用户忘记了密码，需要发重置密码验证码
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

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("RegisterUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// 通过用户邮箱获取用户信息
	u := new(model.User)
	u.Email = req.Email
	ok, err := u.GetUserByEmail()
	if err != nil {
		flog.Log.Errorf("ForgetPassword err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !ok {
		flog.Log.Errorf("ForgetPassword err:%s", "email not found")
		resp.Error = Error(DbNotFound, "email not found")
		return
	}

	// 重设密码验证码过期的话重新设置
	if u.ResetCodeExpired < time.Now().Unix() {
		// 验证码300秒内有效
		err = u.UpdateCode()
		if err != nil {
			flog.Log.Errorf("ForgetPassword err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		// send email
		mm := new(mail.Message)
		mm.Sender = config.FafaConfig.MailConfig
		mm.To = u.Email
		mm.ToName = u.NickName
		mm.Body = "reset password code is: " + u.ResetCode
		err = mm.Sent()
		if err != nil {
			flog.Log.Errorf("ForgetPassword err:%s", err.Error())
			resp.Error = Error(EmailError, err.Error())
			return
		}

	} else {
		flog.Log.Errorf("ForgetPassword err:%s", "time not reach")
		resp.Error = Error(TimeNotReachError, "hold on please")
		return
	}

	resp.Flag = true
}

type ChangePasswordRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Code       string `json:"code" validate:"required,lt=9,gt=5"`
	Password   string `json:"password" validate:"alphanumunicode,gt=5,lt=17"`
	RePassword string `json:"repassword" validate:"eqfield=Password"`
}

// 更改密码，需要用到忘记密码的验证码
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

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("ChangePassword err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// 通过用户邮箱获取用户信息
	u := new(model.User)
	u.Email = req.Email
	ok, err := u.GetUserByEmail()
	if err != nil {
		flog.Log.Errorf("ChangePassword err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !ok {
		flog.Log.Errorf("ChangePassword err:%s", "email not found")
		resp.Error = Error(DbNotFound, "email not found")
		return
	}

	// 验证码一致，可以修改
	if u.ResetCode == req.Code {
		u.Password = req.Password
		err = u.UpdatePassword()
		if err != nil {
			flog.Log.Errorf("ChangePassword err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	} else {
		flog.Log.Errorf("ChangePassword err:%s", "code wrong")
		resp.Error = Error(CodeWrong, "not valid")
		return
	}

	// 更改密码后需要删除登录信息
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

// 用户自己修改自己的信息
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

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("UpdateUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// 获取自己的信息
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateUser err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}

	u := new(model.User)
	u.Id = uu.Id

	// if image not empty
	if req.ImagePath != "" {
		p := new(model.File)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			flog.Log.Errorf("UpdateUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
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
		flog.Log.Errorf("UpdateUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
	resp.Data = u

}

// 用户获取自己的信息
func TakeUser(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()

	u, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("TakeUser err:%s", err.Error())
		resp.Error = Error(I500, "")
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
	Email           string   `json:"email" validate:"omitempty,email"`
	WeChat          string   `json:"wechat" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	WeiBo           string   `json:"weibo" validate:"omitempty,url"`
	Github          string   `json:"github" validate:"omitempty,url"`
	QQ              string   `json:"qq" validate:"omitempty,numeric,gt=6,lt=12"`
	Gender          int      `json:"gender" validate:"oneof=-1 0 1 2"`
	Status          int      `json:"status" validate:"oneof=-1 0 1 2"`
	PageHelp
}

type ListUserResponse struct {
	Users []model.User `json:"users"`
	PageHelp
}

// 列出用户列表，超级管理员权限
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

	var validate = validator.New()
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

type ListGroupUserRequest struct {
	GroupId int `json:"group_id" validate:"lt=0"`
}

type ListGroupUserResponse struct {
	Users []model.User `json:"users"`
}

// 列出组下的用户
func ListGroupUser(c *gin.Context) {
	resp := new(Resp)

	respResult := new(ListGroupUserResponse)
	req := new(ListGroupUserRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("ListGroupUser err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	users := make([]model.User, 0)

	// group list where prepare
	err = session.Table(new(model.User)).Where("group_id=?", req.GroupId).Find(&users)
	if err != nil {
		flog.Log.Errorf("ListUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	respResult.Users = users
	resp.Data = respResult
	resp.Flag = true
}

type AssignGroupRequest struct {
	GroupId      int   `json:"group_id"`
	GroupRelease int   `json:"group_release"`
	Users        []int `json:"users"`
}

// 为用户分配组，每个用户只能有一个组，权限相对弱一点
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

	// 为用户移除组
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

type UpdateUserAdminRequest struct {
	Id       int    `json:"id" validate:"required,gt=0"`
	NickName string `json:"nick_name" validate:"omitempty,gt=1,lt=50"`
	Password string `json:"password,omitempty"`
	Status   int    `json:"status" validate:"oneof=0 1 2"`
}

// 更新用户信息，超级管理员，可以修改用户密码，以及将用户加入黑名单，禁止使用等
func UpdateUserAdmin(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateUserAdminRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("UpdateUserAdmin err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	u := new(model.User)
	u.NickName = req.NickName
	u.Id = req.Id
	u.Password = req.Password

	// 可以将用户拉入黑名单或者激活
	u.Status = req.Status
	err = u.UpdateInfo()
	if err != nil {
		flog.Log.Errorf("UpdateUserAdmin err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
	resp.Data = u
}
