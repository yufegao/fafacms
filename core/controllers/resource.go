package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"math"
)

type ListResourceRequest struct {
	Id    int      `json:"id"`
	Name  string   `json:"name" validate:"omitempty,lt=100"`
	Url   string   `json:"url"`
	Admin int      `json:"admin" validate:"required,oneof=-1 0 1"`
	Sort  []string `json:"sort" validate:"dive,lt=100"`

	PageHelp
}

type ListResourceResponse struct {
	Resources []model.Resource `json:"resources"`
	PageHelp
}

func ListResource(c *gin.Context) {
	resp := new(Resp)

	respResult := new(ListResourceResponse)
	req := new(ListResourceRequest)
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
		flog.Log.Errorf("ListResource err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.Resource)).Where("1=1")

	// query prepare
	if req.Id != 0 {
		session.And("id=?", req.Id)
	}
	if req.Name != "" {
		session.And("name=?", req.Name)
	}

	if req.Admin != -1 {
		session.And("admin=?", req.Admin)
	}

	if req.Url != "" {
		session.And("url=?", req.Url)
	}

	// count num
	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		// db err
		flog.Log.Errorf("ListResource err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if count>0 start list
	r := make([]model.Resource, 0)
	p := &req.PageHelp
	if total == 0 {
	} else {
		// sql build
		p.build(session, req.Sort, model.ResourceSortName)
		// do query
		err = session.Find(&r)
		if err != nil {
			// db err
			flog.Log.Errorf("ListResource err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	// result
	respResult.Resources = r
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

type AssignGroupAndResourceRequest struct {
	GroupId         int   `json:"group_id"`
	ResourceRelease int   `json:"resource_release"`
	Resources       []int `json:"resources"`
}

func AssignGroupAndResource(c *gin.Context) {
	resp := new(Resp)
	req := new(AssignGroupAndResourceRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	resourceNums := len(req.Resources)
	if resourceNums == 0 && req.ResourceRelease != 1 {
		flog.Log.Errorf("AssignGroupAndResource err:%s", "resources empty")
		resp.Error = Error(ParasError, "resources")
		return
	}
	if req.GroupId == 0 {
		flog.Log.Errorf("AssignGroupAndResource err:%s", "group id empty")
		resp.Error = Error(ParasError, "group_id")
		return
	}

	g := new(model.Group)
	g.Id = req.GroupId
	exist, err := g.GetById()
	if err != nil {
		flog.Log.Errorf("AssignGroupAndResource err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("AssignGroupAndResource err:%s", "group not found")
		resp.Error = Error(DbNotFound, "group")
		return
	}

	if resourceNums > 0 {
		num, err := config.FafaRdb.Client.In("id", req.Resources).Count(new(model.Resource))
		if err != nil {
			flog.Log.Errorf("AssignGroupAndResource err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if int(num) != resourceNums {
			flog.Log.Errorf("AssignGroupAndResource err:%s", "resource wrong")
			resp.Error = Error(ParasError, fmt.Sprintf("resource wrong:%d!=%d", num, resourceNums))
			return
		}
	}
	session := config.FafaRdb.Client.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		flog.Log.Errorf("AssignGroupAndResource err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if len(req.Resources) > 0 {
		session.In("resource_id", req.Resources)
	}
	gr := new(model.GroupResource)
	gr.GroupId = req.GroupId
	_, err = session.Cols("group_id").Delete(gr)
	if err != nil {
		session.Rollback()
		flog.Log.Errorf("AssignGroupAndResource err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	rs := make([]model.GroupResource, 0, resourceNums)
	for _, r := range req.Resources {
		rs = append(rs, model.GroupResource{GroupId: req.GroupId, ResourceId: r})
	}
	_, err = session.Insert(rs)
	if err != nil {
		session.Rollback()
		flog.Log.Errorf("AssignGroupAndResource err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	err = session.Commit()
	if err != nil {
		flog.Log.Errorf("AssignGroupAndResource err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
}
