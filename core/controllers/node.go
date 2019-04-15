package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"math"
)

type CreateNodeRequest struct {
	Type         int    `json:"type"  validate:"oneof=0 1"`
	Seo          string `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	Name         string `json:"name" validate:"required,lt=100"`
	Describe     string `json:"describe" validate:"omitempty,lt=200"`
	ImagePath    string `json:"image_path" validate:"omitempty,lt=100"`
	ParentNodeId int    `json:"parent_node_id"`
}

func CreateNode(c *gin.Context) {
	resp := new(Resp)
	req := new(CreateNodeRequest)
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
		flog.Log.Errorf("CreateNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("CreateNode err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}

	n := new(model.ContentNode)
	n.UserId = uu.Id
	n.Type = req.Type
	if req.Seo != "" {
		n.Seo = req.Seo
		exist, err := n.CheckSeoValid()
		if err != nil {
			resp.Error = Error(DBError, "")
			return
		}
		if exist {
			resp.Error = Error(DbRepeat, "field seo")
			return
		}
	}

	if req.ParentNodeId != 0 {
		n.ParentNodeId = req.ParentNodeId
		exist, err := n.CheckParentValid()
		if err != nil {
			resp.Error = Error(DBError, "")
			return
		}
		if !exist {
			resp.Error = Error(DbNotFound, "field parent node")
			return
		}

		n.Level = 1
	}

	// if image not empty
	if req.ImagePath != "" {
		p := new(model.File)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			// db err
			flog.Log.Errorf("CreateNode err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			// not found
			flog.Log.Errorf("CreateNode err: image not exist")
			resp.Error = Error(ParasError, "image url not exist")
			return
		}

		n.ImagePath = req.ImagePath
	}
	n.Name = req.Name
	n.Describe = req.Describe
	n.ParentNodeId = req.ParentNodeId
	err = n.InsertOne()
	if err != nil {
		// db err
		flog.Log.Errorf("CreateNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
	resp.Data = n

}

type UpdateNodeRequest struct {
	Id           int    `json:"id" validate:"required"`
	Seo          string `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	Name         string `json:"name" validate:"omitempty,lt=100"`
	Describe     string `json:"describe" validate:"omitempty,lt=200"`
	ImagePath    string `json:"image_path" validate:"omitempty,lt=100"`
	ParentNodeId int    `json:"parent_node_id"`
	Status       int    `json:"status" validate:"oneof=0 1"`
}

func UpdateNode(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateNodeRequest)
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
		flog.Log.Errorf("UpdateNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateNode err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("UpdateNode err: %s", err.Error())
		resp.Error = Error(DBError, "id empty")
		return
	}
	if !exist {
		flog.Log.Errorf("UpdateNode err: %s", "field id not found")
		resp.Error = Error(DbNotFound, "field id not found")
		return
	}

	if n.Id == req.ParentNodeId {
		flog.Log.Errorf("UpdateNode err: %s", "loop err")
		resp.Error = Error(DbNotFound, "loop err")
		return
	}

	if req.Seo != "" {
		if req.Seo != n.Seo {
			n.Seo = req.Seo
			exist, err := n.CheckSeoValid()
			if err != nil {
				resp.Error = Error(DBError, "")
				return
			}
			if exist {
				resp.Error = Error(DbRepeat, "field seo")
				return
			}
		}
	}

	if req.ParentNodeId != 0 {
		if req.ParentNodeId != n.ParentNodeId {
			n.ParentNodeId = req.ParentNodeId
			exist, err := n.CheckParentValid()
			if err != nil {
				resp.Error = Error(DBError, "")
				return
			}
			if !exist {
				resp.Error = Error(DbNotFound, "field parent node")
				return
			}
			n.Level = 1
		}
	} else {
		n.Level = 0
		n.ParentNodeId = 0
	}

	// if image not empty
	if req.ImagePath != "" {
		if req.ImagePath != n.ImagePath {
			p := new(model.File)
			p.Url = req.ImagePath
			ok, err := p.Exist()
			if err != nil {
				// db err
				flog.Log.Errorf("UpdateNode err:%s", err.Error())
				resp.Error = Error(DBError, err.Error())
				return
			}

			if !ok {
				// not found
				flog.Log.Errorf("UpdateNode err: image not exist")
				resp.Error = Error(ParasError, "image url not exist")
				return
			}

			n.ImagePath = req.ImagePath
		}
	}

	if req.Name != n.Name {
		n.Name = req.Name
	}

	if req.Describe != n.Describe {
		n.Describe = req.Describe
	}

	n.Status = req.Status

	err = n.Update()
	if err != nil {
		// db err
		flog.Log.Errorf("UpdateNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
	resp.Data = n
}

type DeleteNodeRequest struct {
	Id int `json:"id" validate:"required"`
}

func DeleteNode(c *gin.Context) {
	resp := new(Resp)
	req := new(DeleteNodeRequest)
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
		flog.Log.Errorf("DeleteNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("DeleteNode err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// todo
	//n.Get()
}

type TakeNodeRequest struct {
	Id  int    `json:"id"`
	Seo string `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
}

func TakeNode(c *gin.Context) {
	resp := new(Resp)
	req := new(TakeNodeRequest)
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
		flog.Log.Errorf("TakeNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("TakeNode err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id
	n.Seo = req.Seo
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("TakeNode err: %s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}

	if !exist {
		flog.Log.Errorf("TakeNode err: %s", "node not found")
		resp.Error = Error(DbNotFound, "node not found")
		return
	}

	resp.Data = n
	resp.Flag = true
}

type ListNodeRequest struct {
	Id           int    `json:"id"`
	Type         int    `json:"type" validate:"oneof=0 1"`
	Seo          string `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	ParentNodeId int    `json:"parent_node_id"`
	Status       int    `json:"status" validate:"oneof=-1 0 1"`
	Level        int    `json:"status" validate:"oneof=-1 0 1"`

	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	UpdateTimeBegin int64    `json:"update_time_begin"`
	UpdateTimeEnd   int64    `json:"update_time_end"`
	Sort            []string `json:"sort" validate:"dive,lt=100"`
	PageHelp
}

type ListNodeResponse struct {
	Nodes []model.ContentNode `json:"nodes"`
	PageHelp
}

func ListNode(c *gin.Context) {
	resp := new(Resp)

	respResult := new(ListNodeResponse)
	req := new(ListNodeRequest)
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
		flog.Log.Errorf("ListNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("ListNode err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}

	// new query list session
	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.ContentNode)).Where("1=1")

	// query prepare
	if req.Id != 0 {
		session.And("id=?", req.Id)
	}

	session.And("type=?", req.Type)
	session.And("user_id=?", uu.Id)

	if req.Status != -1 {
		session.And("status=?", req.Status)
	}

	if req.Seo != "" {
		session.And("seo=?", req.Seo)
	}

	if req.Level != -1 {
		session.And("level=?", req.Level)
	}

	session.And("parent_node_id=?", req.ParentNodeId)

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
		flog.Log.Errorf("ListNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if count>0 start list
	nodes := make([]model.ContentNode, 0)
	p := &req.PageHelp
	if total == 0 {
	} else {
		// sql build
		p.build(session, req.Sort, model.ContentNodeSortName)
		// do query
		err = session.Find(&nodes)
		if err != nil {
			// db err
			flog.Log.Errorf("ListNode err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	// result
	respResult.Nodes = nodes
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}
