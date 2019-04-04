package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"math"
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
		resp.Error = Error(ParasError, "group name can not empty")
		return
	}

	g := new(model.Group)
	g.Name = req.Name
	ok, err := g.Exist()
	if err != nil {
		flog.Log.Errorf("CreateGroup err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if ok {
		flog.Log.Errorf("CreateGroup err: group name exist")
		resp.Error = Error(ParasError, "group name exist")
		return
	}

	if req.ImagePath != "" {
		g.ImagePath = req.ImagePath
		p := new(model.Picture)
		p.Url = g.ImagePath
		ok, err = p.Exist()
		if err != nil {
			flog.Log.Errorf("CreateGroup err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			flog.Log.Errorf("CreateGroup err: image not exist")
			resp.Error = Error(ParasError, "image url not exist")
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

type UpdateGroupRequest struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Describe  string `json:"describe"`
	ImagePath string `json:"image_path"`
}

func UpdateGroup(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateGroupRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.Id == 0 {
		flog.Log.Errorf("UpdateGroup err: id can not empty")
		resp.Error = Error(ParasError, "group id can not empty")
		return
	}

	g := new(model.Group)
	g.Id = req.Id
	ok, err := g.Exist()

	if err != nil {
		flog.Log.Errorf("UpdateGroup err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !ok {
		resp.Error = Error(DbNotFound, "")
		return
	}

	if req.ImagePath != "" {
		g.ImagePath = req.ImagePath

		p := new(model.Picture)
		p.Url = g.ImagePath
		ok, err := p.Exist()
		if err != nil {
			flog.Log.Errorf("UpdateGroup err:%s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}

		if !ok {
			flog.Log.Errorf("UpdateGroup err: image not exist")
			resp.Error = Error(ParasError, "image url not exist")
			return
		}
	}

	if req.Name != "" {
		temp := new(model.Group)
		temp.Name = req.Name
		ok, err := temp.Exist()
		if err != nil {
			flog.Log.Errorf("UpdateGroup err:%s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}
		if ok {
			resp.Error = Error(DbRepeat, "group name")
			return
		}
		g.Name = req.Name
	}
	if req.Describe != "" {
		g.Describe = req.Describe
	}

	g.UpdateTime = time.Now().Unix()
	_, err = config.FafaRdb.Update(g)
	if err != nil {
		flog.Log.Errorf("UpdateGroup err:%s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}

	resp.Flag = true
	resp.Data = g
}

type DeleteGroupRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func DeleteGroup(c *gin.Context) {
	resp := new(Resp)
	req := new(DeleteGroupRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	//
	temp := new(model.Group)
	temp.Id = req.Id
	temp.Name = req.Name
	ok, err := temp.Take()
	if err != nil {
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !ok {
		resp.Error = Error(DbNotFound, "")
		return
	}

	gr := new(model.GroupResource)
	gr.GroupId = temp.Id
	ok, err = gr.Exist()
	if err != nil {
		resp.Error = Error(DBError, err.Error())
		return
	}
	if ok {
		resp.Error = Error(DbHookIn, "exist resource")
		return
	}

	//
	u := new(model.User)
	u.GroupId = temp.Id
	ok, err = u.Exist()
	if err != nil {
		resp.Error = Error(DBError, err.Error())
		return
	}
	if ok {
		resp.Error = Error(DbHookIn, "exist user")
		return
	}

	g := new(model.Group)
	g.Id = temp.Id
	err = g.Delete()
	if err != nil {
		flog.Log.Errorf("DeleteGroup err:%s", err.Error())
		resp.Error = Error(LazyError, err.Error())
		return
	}

	resp.Flag = true
}

type TakeGroupRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func TakeGroup(c *gin.Context) {
	resp := new(Resp)
	req := new(TakeGroupRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	g := new(model.Group)
	g.Id = req.Id
	g.Name = req.Name
	ok, err := g.Take()
	if err != nil {
		flog.Log.Errorf("TakeGroup err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !ok {
		resp.Error = Error(DbNotFound, "")
		return
	}

	resp.Flag = true
	resp.Data = g
}

type ListGroupRequest struct {
	Id              int      `json:"id"`
	Name            string   `json:"name"`
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	UpdateTimeBegin int64    `json:"update_time_begin"`
	UpdateTimeEnd   int64    `json:"update_time_end"`
	Sort            []string `json:"sort"`
	PageHelp
}

type ListGroupResponse struct {
	Groups []model.Group `json:"groups"`
	PageHelp
}

func ListGroup(c *gin.Context) {
	resp := new(Resp)

	respResult := new(ListGroupResponse)
	req := new(ListGroupRequest)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	session := config.FafaRdb.Client.NewSession()
	defer session.Close()
	session.Table(new(model.Group)).Where("1=1")
	if req.Id != 0 {
		session.And("id=?", req.Id)
	}
	if req.Name != "" {
		session.And("name=?", req.Name)
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

	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		flog.Log.Errorf("ListGroup err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	groups := make([]model.Group, 0)
	p := &req.PageHelp
	if total == 0 {
	} else {
		p.build(session, req.Sort, model.GroupSortName)
		err = session.Find(&groups)
		if err != nil {
			flog.Log.Errorf("ListGroup err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	respResult.Groups = groups

	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}
