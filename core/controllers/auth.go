package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/model"
	"sync"
)

var (
	AuthResource = sync.Map{}
)

// url===>resource_id
func InitAuthResource() {

}

// filter
var AuthFilter = func(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		if resp.Error == nil {
			return
		}
		c.AbortWithStatusJSON(403, resp)
	}()

	u, _ := GetUserSeesion(c)
	if u != nil {
		c.Set("uid", u.Id)
		return
	}



}

func GetUserSeesion(c *gin.Context) (*model.User, error) {
	u := new(model.User)
	s := config.FafaSessionMgr.Load(c.Request)
	err := s.GetObject("user", u)
	if err != nil {
		return nil, err
	}

	if u.Id == 0 {
		return nil, errors.New("no session")
	}
	return u, err
}

func SetUserSession(c *gin.Context, user *model.User) error {
	s := config.FafaSessionMgr.Load(c.Request)
	err := s.PutObject(c.Writer, "user", user)
	return err
}
