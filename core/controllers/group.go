package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"time"
)

type CreateGroupRequest struct {
	Name      string `json:"name"`
	Describe  string `json:"describe"`
	ImagePath string `json:"image_path"`
}

func CreateGroup(c *gin.Context) {
	resp := new(Resp)
	req := new(CreateGroupRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.Name == "" {
		flog.Log.Errorf("CreateGroup err: name can not empty")
		resp.Error = Error(ParasError, "name can not empty")
		return
	}

	g := new(model.Group)
	g.Name = req.Name
	ok, err := config.FafaRdb.Client.Exist(g)
	if err != nil {
		flog.Log.Errorf("CreateGroup err:%s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}

	if ok {
		flog.Log.Errorf("CreateGroup err: group name exist")
		resp.Error = Error(ParasError, "group name exist")
		return
	}

	g.ImagePath = req.ImagePath

	if g.ImagePath != "" {
		p := new(model.Picture)
		p.Url = g.ImagePath

		ok, err := config.FafaRdb.Client.Exist(p)
		if err != nil {
			flog.Log.Errorf("CreateGroup err:%s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}

		if !ok {
			flog.Log.Errorf("CreateGroup err: image not exist")
			resp.Error = Error(ParasError, "image not exist")
			return
		}

	}

	g.Describe = req.Describe

	g.CreateTime = time.Now().Unix()
	_, err = config.FafaRdb.InsertOne(g)
	if err != nil {
		flog.Log.Errorf("CreateGroup err:%s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}
	resp.Flag = true
	resp.Data = g
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
