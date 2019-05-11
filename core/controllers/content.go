package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"math"
)

type CreateContentRequest struct {
	Seo          string `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"` // 内容应该有个好听的标志
	Title        string `json:"title" validate:"required,lt=100"`                    // 必须有标题吧
	Status       int    `json:"status" validate:"oneof=0 1"`                         // 隐藏内容
	Top          int    `json:"top" validate:"oneof=0 1"`                            // 置顶
	Describe     string `json:"describe" validate:"omitempty"`                       // 正文
	ImagePath    string `json:"image_path" validate:"omitempty,lt=100"`              // 内容背景图
	NodeId       int    `json:"node_id"`                                             // 内容所属节点，可以没有节点
	Password     string `json:"password"`                                            // 如果非空表示需要密码
	CloseComment int    `json:"close_comment" validate:"oneof=0 1 2"`                // 评论设置
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
		exist, err := contentNode.Get()
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

		content.NodeSeo = contentNode.Seo
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
	content.Top = req.Top
	content.UserName = uu.Name
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
	Top          int    `json:"top" validate:"oneof=0 1"`
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

	content.NodeSeo = contentBefore.NodeSeo
	content.NodeId = contentBefore.NodeId

	if req.NodeId != 0 && req.NodeId != contentBefore.NodeId {
		content.NodeId = req.NodeId
		contentNode := new(model.ContentNode)
		contentNode.Id = req.NodeId
		contentNode.UserId = uu.Id
		exist, err := contentNode.Get()
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

		content.NodeSeo = contentNode.Seo
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
	if contentBefore.Status <= 1 {
		content.Status = req.Status
	} else {
		content.Status = contentBefore.Status
	}

	// 已经刷新，状态保留
	content.PreFlush = contentBefore.PreFlush

	//  如果内容更新，重置
	if contentBefore.PreDescribe != req.Describe {
		content.PreFlush = 0
		content.PreDescribe = req.Describe
	}

	if contentBefore.Title != req.Title {
		content.Title = req.Title
	}

	content.Password = req.Password
	content.CloseComment = req.CloseComment
	content.Top = req.Top
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

	content.Describe = content.PreDescribe
	err = content.UpdateDescribe()
	if err != nil {
		flog.Log.Errorf("PublishContent err: %s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}
	resp.Flag = true
}

type CancelContentRequest struct {
	Id int `json:"id" validate:"required"`
}

func CancelContent(c *gin.Context) {
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
		flog.Log.Errorf("CancelContent err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("CancelContent err: %s", err.Error())
		resp.Error = Error(I500, "")
		return
	}

	content := new(model.Content)
	content.Id = req.Id
	content.UserId = uu.Id
	exist, err := content.Get()
	if err != nil {
		flog.Log.Errorf("CancelContent err: %s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}

	if !exist {
		flog.Log.Errorf("CancelContent err: %s", "content not found")
		resp.Error = Error(DbNotFound, "content not found")
		return
	}

	if content.PreFlush == 1 {
		resp.Flag = true
		return
	}

	content.PreDescribe = content.Describe
	err = content.ResetDescribe()
	if err != nil {
		flog.Log.Errorf("CancelContent err: %s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}
	resp.Flag = true
}

type ListContentRequest struct {
	Id              int      `json:"id"`
	Seo             string   `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	NodeId          int      `json:"node_id"`
	NodeSeo         string   `json:"node_seo"`
	Top             int      `json:"top" validate:"oneof=-1 0 1"`
	Status          int      `json:"status" validate:"oneof=-1 0 1 2 3 4"`
	CloseComment    int      `json:"close_comment" validate:"oneof=-1 0 1 2"`
	UserId          int      `json:"user_id"`
	UserName        string   `json:"user_name"`
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	UpdateTimeBegin int64    `json:"update_time_begin"`
	UpdateTimeEnd   int64    `json:"update_time_end"`
	Sort            []string `json:"sort" validate:"dive,lt=100"`
	PageHelp
}

type ListContentResponse struct {
	Contents []model.Content `json:"contents"`
	PageHelp
}

func ListContent(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("ListContent err: %s", err.Error())
		resp.Error = Error(I500, "")
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	ListContentHelper(c, uid)
}

func ListContentAdmin(c *gin.Context) {
	ListContentHelper(c, 0)
}

func ListContentHelper(c *gin.Context, userId int) {
	resp := new(Resp)

	respResult := new(ListContentResponse)
	req := new(ListContentRequest)
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
		flog.Log.Errorf("ListContent err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.Content)).Where("1=1")

	// query prepare
	if req.Id != 0 {
		session.And("id=?", req.Id)
	}

	if userId != 0 {
		session.And("user_id=?", userId)
		if req.Status > 3 {
			// 用户不能让他查找到逻辑删除的内容
			req.Status = 0
		}
		session.And("status<?", 4)
	} else {
		if req.UserName != "" {
			session.And("user_name=?", req.UserName)
		}
		if req.UserId != 0 {
			session.And("user_id=?", req.UserId)
		}
	}

	if req.Status != -1 {
		session.And("status=?", req.Status)
	}

	if req.Top != -1 {
		session.And("top=?", req.Top)
	}

	if req.Seo != "" {
		session.And("seo=?", req.Seo)
	}

	if req.CloseComment != -1 {
		session.And("close_comment=?", req.CloseComment)
	}

	if req.NodeId != 0 {
		session.And("node_id=?", req.NodeId)
	}

	if req.NodeSeo != "" {
		session.And("node_seo=?", req.NodeSeo)
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
		flog.Log.Errorf("ListContent err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if count>0 start list
	cs := make([]model.Content, 0)
	p := &req.PageHelp
	if total == 0 {
	} else {
		// sql build
		p.build(session, req.Sort, model.ContentSortName)
		// do query
		err = session.Omit("describe", "pre_describe").Find(&cs)
		if err != nil {
			flog.Log.Errorf("ListContent err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	// result
	respResult.Contents = cs
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

type ListContentHistoryRequest struct {
	Id     int      `json:"id" validate:"required"`
	UserId int      `json:"user_id"`
	Sort   []string `json:"sort" validate:"dive,lt=100"`
	PageHelp
}

type ListContentHistoryResponse struct {
	Contents []model.ContentHistory `json:"contents"`
	PageHelp
}

func ListContentHistory(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("ListContentHistory err: %s", err.Error())
		resp.Error = Error(I500, "")
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	ListContentHistoryHelper(c, uid)
}

func ListContentHistoryAdmin(c *gin.Context) {
	ListContentHistoryHelper(c, 0)
}

func ListContentHistoryHelper(c *gin.Context, userId int) {
	resp := new(Resp)

	respResult := new(ListContentHistoryResponse)
	req := new(ListContentHistoryRequest)
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
		flog.Log.Errorf("ListContentHistory err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.ContentHistory)).Where("1=1")

	session.And("content_id=?", req.Id)

	if userId != 0 {
		session.And("user_id=?", userId)
	} else {
		if req.UserId != 0 {
			session.And("user_id=?", req.UserId)
		}
	}

	// count num
	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		flog.Log.Errorf("ListContentHistory err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if count>0 start list
	cs := make([]model.ContentHistory, 0)
	p := &req.PageHelp
	if total == 0 {
	} else {
		// sql build
		p.build(session, req.Sort, model.ContentHistorySortName)
		// do query
		err = session.Omit("describe").Find(&cs)
		if err != nil {
			flog.Log.Errorf("ListContentHistory err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	// result
	respResult.Contents = cs
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

type TakeContentRequest struct {
	Id int `json:"id" validate:"required"`
}

func TakeContentHelper(c *gin.Context, userId int) {
	resp := new(Resp)
	req := new(TakeContentRequest)
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
		flog.Log.Errorf("TakeContent err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	content := new(model.Content)
	content.Id = req.Id
	content.UserId = userId
	exist, err := content.GetByAdmin()
	if err != nil {
		flog.Log.Errorf("TakeContent err: %s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}

	if !exist {
		flog.Log.Errorf("TakeContent err: %s", "content not found")
		resp.Error = Error(DbNotFound, "content not found")
		return
	}

	resp.Data = content
	resp.Flag = true
}

func TakeContent(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("TakeContent err: %s", err.Error())
		resp.Error = Error(I500, "")
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	TakeContentHelper(c, uid)
}

func TakeContentAdmin(c *gin.Context) {
	TakeContentHelper(c, 0)
}

type TakeContentHistoryRequest struct {
	Id int `json:"id" validate:"required"`
}

func TakeContentHistoryHelper(c *gin.Context, userId int) {
	resp := new(Resp)
	req := new(TakeContentHistoryRequest)
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
		flog.Log.Errorf("TakeContentHistory err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	content := new(model.ContentHistory)
	content.Id = req.Id
	content.UserId = userId
	exist, err := content.GetByAdmin()
	if err != nil {
		flog.Log.Errorf("TakeContentHistory err: %s", err.Error())
		resp.Error = Error(DBError, "")
		return
	}

	if !exist {
		flog.Log.Errorf("TakeContentHistory err: %s", "content not found")
		resp.Error = Error(DbNotFound, "content not found")
		return
	}

	resp.Data = content
	resp.Flag = true
}

func TakeContentHistory(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("TakeContentHistory err: %s", err.Error())
		resp.Error = Error(I500, "")
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	TakeContentHistoryHelper(c, uid)
}

func TakeContentHistoryAdmin(c *gin.Context) {
	TakeContentHistoryHelper(c, 0)
}

type DeleteContentRequest struct {
	Id     int `json:"id" validate:"required"`
	Status int `json:"status" validate:"oneof=0 1 2 3 4"`
}

// 文章状态操作
func DeleteContentHelper(c *gin.Context, userId int, typeDelete int) {
	resp := new(Resp)
	req := new(DeleteContentRequest)
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
		flog.Log.Errorf("DeleteContent err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	content := new(model.Content)
	content.Id = req.Id
	content.UserId = userId
	if typeDelete == 0 {
		// 送去垃圾站
		_, err = content.UpdateStatusTo3()
	} else if typeDelete == 1 {
		// 逻辑删除 已经修改为真正意义上的删除，毕竟空间有限
		err = content.UpdateStatusTo4()
	} else if typeDelete == 2 {
		// 垃圾恢复
		_, err = content.UpdateStatusTo3Reverse()
	} else {
		// 管理员权限
		content.Status = req.Status
		_, err = content.UpdateStatus()
	}
	resp.Flag = true
	return
}

// 垃圾回收，假删除
func DeleteContent(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("DeleteContent err: %s", err.Error())
		resp.Error = Error(I500, "")
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	DeleteContentHelper(c, uid, 0)
}

// 逻辑删除
func ReallyDeleteContent(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("DeleteContent err: %s", err.Error())
		resp.Error = Error(I500, "")
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	DeleteContentHelper(c, uid, 1)
}

// 垃圾恢复
func DeleteContentRedo(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("DeleteContent err: %s", err.Error())
		resp.Error = Error(I500, "")
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	DeleteContentHelper(c, uid, 2)
}

// 管理员超大权限
func DeleteContentAdmin(c *gin.Context) {
	DeleteContentHelper(c, 0, 3)
}
