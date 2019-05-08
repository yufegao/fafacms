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
	Seo          string `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	Name         string `json:"name" validate:"required,lt=100"`
	Describe     string `json:"describe" validate:"omitempty,lt=200"`
	ImagePath    string `json:"image_path" validate:"omitempty,lt=100"`
	ParentNodeId int    `json:"parent_node_id"`
	SortNum      int    `json:"sort_num"`
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

	// 如果SEO非空，检查是否已经存在
	if req.Seo != "" {
		n.Seo = req.Seo
		exist, err := n.CheckSeoValid()
		if err != nil {
			resp.Error = Error(DBError, "")
			return
		}
		if exist {
			// 存在报错
			resp.Error = Error(DbRepeat, "field seo")
			return
		}
	}

	// 如果指定了父亲节点
	if req.ParentNodeId != 0 {
		n.ParentNodeId = req.ParentNodeId
		exist, err := n.CheckParentValid()
		if err != nil {
			resp.Error = Error(DBError, "")
			return
		}
		if !exist {
			// 父亲节点不存在，报错
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
	n.SortNum = req.SortNum
	err = n.InsertOne()
	if err != nil {
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
	Status       int    `json:"status" validate:"oneof=-1 0 1"`
	SortNum      int    `json:"sort_num"`
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

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("UpdateNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	if req.ParentNodeId == req.Id {
		flog.Log.Errorf("UpdateNode err: %s", "self can not be parent")
		resp.Error = Error(ParasError, "self can not be parent")
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

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("UpdateNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("UpdateNode err: %s", "field id not found")
		resp.Error = Error(DbNotFound, "field id not found")
		return
	}

	// 不能将自己作为自己的父亲
	if n.Id == req.ParentNodeId {
		flog.Log.Errorf("UpdateNode err: %s", "loop err")
		resp.Error = Error(DbNotFound, "loop err")
		return
	}

	// SEO不为空
	if req.Seo != "" {
		// 和之前的SEO不一样
		if req.Seo != n.Seo {
			n.Seo = req.Seo
			// 检查是否存在SEO
			exist, err := n.CheckSeoValid()
			if err != nil {
				resp.Error = Error(DBError, "")
				return
			}
			if exist {
				// SEO存在了，报错
				resp.Error = Error(DbRepeat, "field seo")
				return
			}
		}
	}

	// 指定了父亲节点
	if req.ParentNodeId > 0 {
		// 和之前的父亲节点不一样
		if req.ParentNodeId != n.ParentNodeId {
			n.ParentNodeId = req.ParentNodeId
			// 检查该父亲节点是否存在
			exist, err := n.CheckParentValid()
			if err != nil {
				resp.Error = Error(DBError, "")
				return
			}
			if !exist {
				// 不存在父亲节点，报错
				resp.Error = Error(DbNotFound, "field parent node")
				return
			}
			// 有了父亲节点，级别为1
			n.Level = 1
		}
	} else if req.ParentNodeId == -1 {
	} else {
		// 没有指定父亲节点，归零
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
				flog.Log.Errorf("UpdateNode err:%s", err.Error())
				resp.Error = Error(DBError, err.Error())
				return
			}

			if !ok {
				flog.Log.Errorf("UpdateNode err: image not exist")
				resp.Error = Error(ParasError, "image url not exist")
				return
			}

			n.ImagePath = req.ImagePath
		}
	}

	// 以下只要存在不一致性才替换
	if req.Name != "" {
		if req.Name != n.Name {
			n.Name = req.Name
		}
	}

	if req.Describe != "" {
		if req.Describe != n.Describe {
			n.Describe = req.Describe
		}
	}

	if n.Status != -1 {
		if n.Status != req.Status {
			n.Status = req.Status
		}
	}

	if n.SortNum != -1 {
		if n.SortNum != req.SortNum {
			n.SortNum = req.SortNum
		}
	}

	// 更新
	err = n.Update()
	if err != nil {
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

	// 删除节点时节点下不能有节点
	childNum, err := n.CheckChildrenNum()
	if err != nil {
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if childNum >= 1 {
		// 不能删除
		flog.Log.Errorf("DeleteNode err:%s", "has node child")
		resp.Error = Error(DbHookIn, "has node child")
		return
	}

	content := new(model.Content)
	content.UserId = uu.Id
	content.NodeId = n.Id

	// 删除节点时，节点下不能有内容
	allContentNum, normalContentNum, err := content.CountNumOfNode()
	if err != nil {
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if normalContentNum >= 1 {
		// 有内容，不能删除
		flog.Log.Errorf("DeleteNode err:%s", "has content child")
		resp.Error = Error(DbHookIn, "has content child")
		return
	}

	// 如果从来没有删除过，请直接删除
	if allContentNum == 0 {
		err = n.Delete()
	} else {
		//逻辑删除
		err = n.LogicDelete()
	}

	if err != nil {
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
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
		resp.Error = Error(DBError, err.Error())
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
	Id              int      `json:"id"`
	Seo             string   `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	ParentNodeId    int      `json:"parent_node_id"`
	Status          int      `json:"status" validate:"oneof=-1 0 1 2"`
	Level           int      `json:"level" validate:"oneof=-1 0 1"`
	UserId          int      `json:"user_id"`
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
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("ListNode err: %s", err.Error())
		resp.Error = Error(I500, "")
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	ListNodeHelper(c, uid)
}

func ListNodeAdmin(c *gin.Context) {
	ListNodeHelper(c, 0)
}

func ListNodeHelper(c *gin.Context, userId int) {
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

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("ListNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
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

	if userId != 0 {
		session.And("user_id=?", userId)
		if req.Status > 1 {
			// 用户不能让他查找到逻辑删除的节点
			req.Status = 0
		}
	} else {
		if req.UserId != 0 {
			session.And("user_id=?", req.UserId)
		}
	}

	if req.Status != -1 {
		session.And("status=?", req.Status)
	}

	if req.Seo != "" {
		session.And("seo=?", req.Seo)
	}

	if req.Level != -1 {
		session.And("level=?", req.Level)
	}

	if req.ParentNodeId != -1 {
		session.And("parent_node_id=?", req.ParentNodeId)
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
