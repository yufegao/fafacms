package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"math"
	"time"
)

type CreateGroupRequest struct {
	Name      string `json:"name" validate:"required,gt=1,lt=100"`
	Describe  string `json:"describe" validate:"lt=100"`
	ImagePath string `json:"image_path" validate:"lt=100"`
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

	// validate
	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("CreateGroup err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// if exist group
	g := new(model.Group)
	g.Name = req.Name
	ok, err := g.Exist()
	if err != nil {
		// db err
		flog.Log.Errorf("CreateGroup err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if ok {
		// found
		flog.Log.Errorf("CreateGroup err: group name exist")
		resp.Error = Error(ParasError, "group name exist")
		return
	}

	// if image not empty
	if req.ImagePath != "" {
		// picture table exist
		g.ImagePath = req.ImagePath
		p := new(model.File)
		p.Url = g.ImagePath
		ok, err = p.Exist()
		if err != nil {
			// db err
			flog.Log.Errorf("CreateGroup err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			// not found
			flog.Log.Errorf("CreateGroup err: image not exist")
			resp.Error = Error(ParasError, "image url not exist")
			return
		}

	}

	// insert now
	g.Describe = req.Describe
	g.CreateTime = time.Now().Unix()
	_, err = config.FafaRdb.InsertOne(g)
	if err != nil {
		// db err
		flog.Log.Errorf("CreateGroup err:%s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}
	resp.Flag = true
	resp.Data = g
}

type UpdateGroupRequest struct {
	Id        int    `json:"id" validate:"required,gt=0"`
	Name      string `json:"name" validate:"lt=100"`
	Describe  string `json:"describe" validate:"lt=100"`
	ImagePath string `json:"image_path" validate:"lt=100"`
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

	// validate
	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("UpdateGroup err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// if group exist
	g := new(model.Group)
	g.Id = req.Id
	ok, err := g.Exist()
	if err != nil {
		// db err
		flog.Log.Errorf("UpdateGroup err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !ok {
		// not found
		resp.Error = Error(DbNotFound, "")
		return
	}

	// if image not empty
	if req.ImagePath != "" {
		g.ImagePath = req.ImagePath
		p := new(model.File)
		p.Url = g.ImagePath
		// find picture table
		ok, err := p.Exist()
		if err != nil {
			// db err
			flog.Log.Errorf("UpdateGroup err:%s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}

		if !ok {
			// not found
			flog.Log.Errorf("UpdateGroup err: image not exist")
			resp.Error = Error(ParasError, "image url not exist")
			return
		}
	}

	// if group name change repeat
	if req.Name != "" {
		temp := new(model.Group)
		temp.Name = req.Name
		// exist the same name
		ok, err := temp.Exist()
		if err != nil {
			// db err
			flog.Log.Errorf("UpdateGroup err:%s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}
		if ok {
			// found
			resp.Error = Error(DbRepeat, "group name")
			return
		}
		g.Name = req.Name
	}

	if req.Describe != "" {
		g.Describe = req.Describe
	}

	err = g.Update()
	if err != nil {
		// db err
		flog.Log.Errorf("UpdateGroup err:%s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}

	resp.Flag = true
	resp.Data = g
}

type DeleteGroupRequest struct {
	Id   int    `json:"id" `
	Name string `json:"name" validate:"lt=100"`
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

	// validate
	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("DeleteGroup err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// take group info
	temp := new(model.Group)
	temp.Id = req.Id
	temp.Name = req.Name
	ok, err := temp.Take()
	if err != nil {
		// db err
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !ok {
		// not found
		resp.Error = Error(DbNotFound, "")
		return
	}

	// resource exist under group
	gr := new(model.GroupResource)
	gr.GroupId = temp.Id
	ok, err = gr.Exist()
	if err != nil {
		// db err
		resp.Error = Error(DBError, err.Error())
		return
	}
	if ok {
		// found can not delete
		resp.Error = Error(DbHookIn, "exist resource")
		return
	}

	// user exist under group
	u := new(model.User)
	u.GroupId = temp.Id
	ok, err = u.Exist()
	if err != nil {
		// db err
		resp.Error = Error(DBError, err.Error())
		return
	}
	if ok {
		// found can not delete
		resp.Error = Error(DbHookIn, "exist user")
		return
	}

	// delete group
	g := new(model.Group)
	g.Id = temp.Id
	err = g.Delete()
	if err != nil {
		// db err
		flog.Log.Errorf("DeleteGroup err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
}

type TakeGroupRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name" validate:"lt=100"`
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

	// validate
	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("TakeGroup err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// take group info
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
	Name            string   `json:"name" validate:"lt=100"`
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	UpdateTimeBegin int64    `json:"update_time_begin"`
	UpdateTimeEnd   int64    `json:"update_time_end"`
	Sort            []string `json:"sort" validate:"dive,lt=100"`
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
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	// validate
	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("ListGroup err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.Group)).Where("1=1")

	// query prepare
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

	// count num
	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		// db err
		flog.Log.Errorf("ListGroup err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if count>0 start list
	groups := make([]model.Group, 0)
	p := &req.PageHelp
	if total == 0 {
	} else {
		// sql build
		p.build(session, req.Sort, model.GroupSortName)
		// do query
		err = session.Find(&groups)
		if err != nil {
			// db err
			flog.Log.Errorf("ListGroup err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	// result
	respResult.Groups = groups
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}
