package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
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
	n.Type = req.Type
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

func UpdateNode(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func DeleteNode(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func TakeNode(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

type ListNodeRequest struct {
	Id           int    `json:"id"`
	Type         int    `json:"type"  validate:"oneof=0 1"`
	Seo          string `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	ParentNodeId int    `json:"parent_node_id"`
	Status       int    `json:"status" validate:"oneof=0 1"`
}

func ListNode(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}
