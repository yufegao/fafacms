package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
)

type CreateContentRequest struct {
	Seo          string `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"` // 内容应该有个好听的标志
	Title        string `json:"title" validate:"required,lt=100"`                    // 必须有标题吧
	Status       int    `json:"status" validate:"oneof=0 1"`                         // 隐藏内容
	Describe     string `json:"describe" validate:"omitempty"`                       // 正文
	ImagePath    string `json:"image_path" validate:"omitempty,lt=100"`              // 内容背景图
	NodeId       int    `json:"node_id"`                                             // 内容所属节点，可以没有节点
	Password     string `json:"password"`                                            // 如果非空表示需要密码
	CloseComment int    `json:"close_comment" validate:"oneof=0 1 2"`
}

func CreateContent(c *gin.Context) {
	resp := new(Resp)
	req := new(CreateContentRequest)
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
		flog.Log.Errorf("CreateContent err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("CreateContent err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}

	content := new(model.Content)
	content.UserId = uu.Id
	if req.Seo != "" {
		content.Seo = req.Seo
		exist, err := content.CheckSeoValid()
		if err != nil {
			flog.Log.Errorf("CreateContent err: %s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}
		if exist {
			flog.Log.Errorf("CreateContent err: %s", "seo repeat")
			resp.Error = Error(DbRepeat, "seo repeat")
			return
		}
	}

	if req.NodeId != 0 {
		content.NodeId = req.NodeId
		contentNode := new(model.ContentNode)
		contentNode.Id = req.NodeId
		contentNode.UserId = uu.Id
		exist, err := contentNode.Exist()
		if err != nil {
			flog.Log.Errorf("CreateContent err: %s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}
		if !exist {
			flog.Log.Errorf("CreateContent err: %s", "node not found")
			resp.Error = Error(DbNotFound, "node_id")
			return
		}
	}

	if req.ImagePath != "" {
		p := new(model.File)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			flog.Log.Errorf("CreateContent err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			flog.Log.Errorf("CreateContent err: image not exist")
			resp.Error = Error(ParasError, "image url not exist")
			return
		}

		content.ImagePath = req.ImagePath
	}

	content.Status = req.Status
	content.PreDescribe = req.Describe
	content.Title = req.Title
	content.Password = req.Password
	content.CloseComment = req.CloseComment
	_, err = content.Insert()
	if err != nil {
		flog.Log.Errorf("CreateContent err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Data = content
	resp.Flag = true
}

type UpdateContentRequest struct {
	Id           int    `json:"id" validate:"required"`
	Seo          string `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	Title        string `json:"title" validate:"required,lt=100"`
	Status       int    `json:"status" validate:"oneof=0 1"`
	Describe     string `json:"describe" validate:"omitempty"`
	ImagePath    string `json:"image_path" validate:"omitempty,lt=100"`
	NodeId       int    `json:"node_id"`
	Password     string `json:"password"`
	CloseComment int    `json:"close_comment" validate:"oneof=0 1 2"`
}

func UpdateContent(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateContentRequest)
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
		flog.Log.Errorf("UpdateContent err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateContent err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}

	contentBefore := new(model.Content)
	contentBefore.Id = req.Id
	contentBefore.UserId = uu.Id
	exist, err := contentBefore.Get()
	if err != nil {
		flog.Log.Errorf("UpdateContent err: %s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}

	if !exist {
		flog.Log.Errorf("UpdateContent err: %s", "content not found")
		resp.Error = Error(DbNotFound, "content not found")
		return
	}

	content := new(model.Content)
	content.Id = req.Id
	content.UserId = uu.Id
	if req.Seo != "" && req.Seo != contentBefore.Seo {
		content.Seo = req.Seo
		exist, err := content.CheckSeoValid()
		if err != nil {
			flog.Log.Errorf("UpdateContent err: %s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}
		if exist {
			flog.Log.Errorf("UpdateContent err: %s", "seo repeat")
			resp.Error = Error(DbRepeat, "seo repeat")
			return
		}
	}

	if req.NodeId != 0 && req.NodeId != contentBefore.NodeId {
		content.NodeId = req.NodeId
		contentNode := new(model.ContentNode)
		contentNode.Id = req.NodeId
		contentNode.UserId = uu.Id
		exist, err := contentNode.Exist()
		if err != nil {
			flog.Log.Errorf("UpdateContent err: %s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}
		if !exist {
			flog.Log.Errorf("UpdateContent err: %s", "node not found")
			resp.Error = Error(DbNotFound, "node_id")
			return
		}
	}

	if req.ImagePath != "" && req.ImagePath != contentBefore.ImagePath {
		p := new(model.File)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			flog.Log.Errorf("UpdateContent err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			flog.Log.Errorf("UpdateContent err: image not exist")
			resp.Error = Error(ParasError, "image url not exist")
			return
		}

		content.ImagePath = req.ImagePath
	}

	// 只可以修改0-1状态的内容，即正常和不显示的内容
	if contentBefore.Status <= 1 && req.Status != contentBefore.Status {
		content.Status = req.Status
	}

	// 已经刷新，状态保留
	content.PreFlush = contentBefore.PreFlush

	//  如果内容更新，重置
	if content.PreDescribe != req.Describe {
		content.PreFlush = 0
		content.PreDescribe = req.Describe
	}

	if content.Title != contentBefore.Title {
		content.Title = req.Title
	}

	if req.Password != contentBefore.Password {
		content.Password = req.Password
	}

	if req.CloseComment != contentBefore.CloseComment {
		content.CloseComment = req.CloseComment
	}
	_, err = content.Update()
	if err != nil {
		flog.Log.Errorf("UpdateContent err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
}

type PublishContentRequest struct {
	Id int `json:"id" validate:"required"`
}

func PublishContent(c *gin.Context) {
	resp := new(Resp)
	req := new(PublishContentRequest)
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
		flog.Log.Errorf("PublishContent err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("PublishContent err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}

	content := new(model.Content)
	content.Id = req.Id
	content.UserId = uu.Id
	exist, err := content.Get()
	if err != nil {
		flog.Log.Errorf("PublishContent err: %s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}

	if !exist {
		flog.Log.Errorf("PublishContent err: %s", "content not found")
		resp.Error = Error(DbNotFound, "content not found")
		return
	}

	if content.PreFlush == 1 {
		resp.Flag = true
		return
	}

	content.Describe =content.PreDescribe
	content.UpdateDescribe()

}
func DeleteContent(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func TakeContent(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func ListContent(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}
