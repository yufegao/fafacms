package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/parrot/util"
	"strconv"
	"strings"
)

var AuthDebug = false

// auth filter
var AuthFilter = func(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		if resp.Error == nil {
			return
		}
		c.AbortWithStatusJSON(403, resp)
	}()

	// get session
	u, _ := GetUserSession(c)
	if u == nil {
		// if not exist session check cookie
		success, userInfo := CheckCookie(c)
		if success {
			// set session
			err := SetUserSession(c, userInfo)
			if err != nil {
				flog.Log.Errorf("filter err:%s", err.Error())
				resp.Error = Error(I500, "")
				return
			}
			u = userInfo
		} else {
			// cookie and session not exist, nologin
			flog.Log.Errorf("filter err: %s", "no cookie")
			resp.Error = Error(NoLogin, "")
			return
		}
	}

	// record log will need uid, monitor who op
	c.Set("uid", u.Id)

	if AuthDebug {
		return
	}

	// root user can ignore auth
	if u.Id == -1 {
		return
	}

	//  get groupId by user
	nowUser := new(model.User)
	nowUser.Id = u.Id
	err := nowUser.Get()
	if err != nil {
		flog.Log.Errorf("filter err:%s", err.Error())
		resp.Error = Error(AuthPermit, "")
		return
	}

	if nowUser.Status == 0 {
		flog.Log.Errorf("filter err: not active")
		resp.Error = Error(AuthPermit, "not active")
		return
	}

	// resource is exist
	r := new(model.Resource)
	url := c.Request.URL.Path
	r.Url = url
	r.Admin = true

	// resource not found can skip auth
	if err := r.Get(); err != nil {
		flog.Log.Warnf("resource found url:%s, auth err:%s", url, err.Error())
		return
	}

	// if group has this resource
	gr := new(model.GroupResource)
	gr.GroupId = nowUser.GroupId
	gr.ResourceId = r.Id
	exist, err := config.FafaRdb.Client.Exist(gr)
	if err != nil {
		flog.Log.Errorf("filter err:%s", err.Error())
		resp.Error = Error(AuthPermit, "")
		return
	}

	if !exist {
		// not found
		flog.Log.Errorf("filter err:%s", "resource not allow")
		resp.Error = Error(AuthPermit, "")
		return
	}
}

func CheckCookie(c *gin.Context) (success bool, user *model.User) {
	// cookie store a string
	cookieString, err := c.Cookie("auth")
	if err != nil {
		return false, nil
	}

	// cookie string split
	arr := strings.Split(cookieString, "|")
	if len(arr) < 2 {
		// cookie clean
		c.SetCookie("auth", "", -1, "/", "", false, true)
		return
	}

	// userId and md5(ip+password) get
	var userId int64
	str, password := arr[0], arr[1]
	userId, err = strconv.ParseInt(str, 10, 0)
	if err != nil {
		return
	}

	// get user password
	user = &model.User{}
	user.Id = int(userId)
	err = user.Get()
	if err != nil {
		return
	}

	// if the same
	if password == util.Md5(c.ClientIP()+"|"+user.Password) {
		success = true
		return
	} else {
		// cookie clean
		c.SetCookie("auth", "", -1, "/", "", false, true)
		return
	}
}

func GetUserSession(c *gin.Context) (*model.User, error) {
	u := new(model.User)
	s := config.FafaSessionMgr.Load(c.Request)

	// get session from redis..
	err := s.GetObject("user", u)
	if err != nil {
		return nil, err
	}

	// not found
	if u.Id == 0 {
		return nil, errors.New("no session")
	}
	return u, err
}

func SetUserSession(c *gin.Context, user *model.User) error {
	s := config.FafaSessionMgr.Load(c.Request)
	user.Password = ""
	user.ActivateExpired = 0
	user.ActivateMd5 = ""
	err := s.PutObject(c.Writer, "user", user)
	return err
}

func DeleteUserSession(c *gin.Context) error {
	s := config.FafaSessionMgr.Load(c.Request)
	err := s.Destroy(c.Writer)
	return err
}
