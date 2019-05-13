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
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	n := new(model.ContentNode)
	n.UserId = uu.Id

	// 如果SEO非空，检查是否已经存在
	if req.Seo != "" {
		n.Seo = req.Seo
		exist, err := n.CheckSeoValid()
		if err != nil {
			flog.Log.Errorf("CreateNode err: %s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}
		if exist {
			// 存在报错
			flog.Log.Errorf("CreateNode err: %s", "node seo already be use")
			resp.Error = Error(ContentNodeSeoAlreadyBeUsed, "")
			return
		}
	}

	// 如果指定了父亲节点
	if req.ParentNodeId != 0 {
		n.ParentNodeId = req.ParentNodeId
		exist, err := n.CheckParentValid()
		if err != nil {
			flog.Log.Errorf("CreateNode err: %s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}
		if !exist {
			// 父亲节点不存在，报错
			flog.Log.Errorf("CreateNode err: %s", "parent content node not found")
			resp.Error = Error(ContentParentNodeNotFound, "")
			return
		}

		n.Level = 1
	}

	// if image not empty
	if req.ImagePath != "" {
		n.ImagePath = req.ImagePath
		p := new(model.File)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			flog.Log.Errorf("CreateNode err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			flog.Log.Errorf("CreateNode err: image not exist")
			resp.Error = Error(FileCanNotBeFound, "image url not exist")
			return
		}
	}
	n.Name = req.Name
	n.Describe = req.Describe
	n.ParentNodeId = req.ParentNodeId
	n.UserName = uu.Name
	n.SortNum, _ = n.CountNodeNum()
	err = n.InsertOne()
	if err != nil {
		flog.Log.Errorf("CreateNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
	resp.Data = n
}

type UpdateInfoOfNodeRequest struct {
	Id        int    `json:"id" validate:"required"`
	Name      string `json:"name" validate:"omitempty,lt=100"`
	Describe  string `json:"describe" validate:"omitempty,lt=200"`
	ImagePath string `json:"image_path" validate:"omitempty,lt=100"`
}

type UpdateStatusOfNodeRequest struct {
	Id     int `json:"id" validate:"required"`
	Status int `json:"status" validate:"oneof=0 1"`
}

type UpdateSeoOfNodeRequest struct {
	Id  int    `json:"id" validate:"required"`
	Seo string `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
}

type UpdateParentOfNodeRequest struct {
	Id           int  `json:"id" validate:"required"`
	ToBeRoot     bool `json:"to_be_root"` // 升级为最上层节点
	ParentNodeId int  `json:"parent_node_id"`
}

func UpdateSeoOfNode(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateSeoOfNodeRequest)
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
		flog.Log.Errorf("UpdateSeoOfNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateSeoOfNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("UpdateSeoOfNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("UpdateSeoOfNode err: %s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	after := new(model.ContentNode)
	after.UserId = n.UserId
	after.Id = n.Id

	seoChange := false
	// SEO不为空
	if req.Seo != "" {
		// 和之前的SEO不一样
		if req.Seo != n.Seo {
			after.Seo = req.Seo
			seoChange = true
			// 检查是否存在SEO
			exist, err := after.CheckSeoValid()
			if err != nil {
				flog.Log.Errorf("UpdateSeoOfNode err: %s", err.Error())
				resp.Error = Error(DBError, err.Error())
				return
			}
			if exist {
				// SEO存在了，报错
				flog.Log.Errorf("UpdateSeoOfNode err: %s", err.Error())
				resp.Error = Error(ContentNodeSeoAlreadyBeUsed, "")
				return
			}
		}
	}

	if seoChange {
		// 更新
		err = after.UpdateSeo()
		if err != nil {
			flog.Log.Errorf("UpdateSeoOfNode err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}
	resp.Flag = true
}
func UpdateInfoOfNode(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateInfoOfNodeRequest)
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
		flog.Log.Errorf("UpdateInfoOfNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateInfoOfNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, "")
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("UpdateInfoOfNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("UpdateInfoOfNode err: %s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	after := new(model.ContentNode)
	after.UserId = n.UserId
	after.Id = n.Id

	// if image not empty
	if req.ImagePath != "" {
		if req.ImagePath != n.ImagePath {
			after.ImagePath = req.ImagePath
			p := new(model.File)
			p.Url = req.ImagePath
			ok, err := p.Exist()
			if err != nil {
				flog.Log.Errorf("UpdateInfoOfNode err:%s", err.Error())
				resp.Error = Error(DBError, err.Error())
				return
			}

			if !ok {
				flog.Log.Errorf("UpdateInfoOfNode err: image not exist")
				resp.Error = Error(FileCanNotBeFound, "")
				return
			}
		}
	}

	// 以下只要存在不一致性才替换
	if req.Name != "" {
		if req.Name != n.Name {
			after.Name = req.Name
		}
	}

	after.Describe = req.Describe

	// 更新
	err = after.UpdateInfo()
	if err != nil {
		flog.Log.Errorf("UpdateNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
}

func UpdateStatusOfNode(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateStatusOfNodeRequest)
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
		flog.Log.Errorf("UpdateStatusOfNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateStatusOfNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("UpdateStatusOfNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("UpdateStatusOfNode err: %s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	after := new(model.ContentNode)
	after.UserId = n.UserId
	after.Id = n.Id
	after.Status = req.Status

	// 更新
	err = after.UpdateStatus()
	if err != nil {
		flog.Log.Errorf("UpdateStatusOfNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
}

func UpdateParentOfNode(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateParentOfNodeRequest)
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
		flog.Log.Errorf("UpdateParentOfNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	if req.ParentNodeId == req.Id {
		flog.Log.Errorf("UpdateParentOfNode err: %s", "self can not be parent")
		resp.Error = Error(ParasError, "self can not be parent")
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateParentOfNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("UpdateParentOfNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("UpdateParentOfNode err: %s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	after := new(model.ContentNode)
	after.UserId = n.UserId
	after.Id = n.Id

	if req.ToBeRoot {
		if n.ParentNodeId == 0 {
			resp.Flag = true
			return
		}
		// 没有指定父亲节点，归零
		after.Level = 0
		after.ParentNodeId = 0
	} else {
		if n.ParentNodeId == req.ParentNodeId {
			resp.Flag = true
			return
		}

		after.ParentNodeId = req.ParentNodeId

		// 检查该父亲节点是否存在
		exist, err := after.CheckParentValid()
		if err != nil {
			flog.Log.Errorf("UpdateParentOfNode err: %s", err.Error())
			resp.Error = Error(DBError, "")
			return
		}
		if !exist {
			// 不存在父亲节点，报错
			flog.Log.Errorf("UpdateParentOfNode err: %s", "parent content node not found")
			resp.Error = Error(ContentParentNodeNotFound, "")
			return
		}
		after.Level = 1

	}

	// 更新
	err = after.UpdateParent()
	if err != nil {
		flog.Log.Errorf("UpdateParentOfNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
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
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("DeleteNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("DeleteNode err: %s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

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
		resp.Error = Error(ContentNodeHasChildren, "")
		return
	}

	content := new(model.Content)
	content.UserId = uu.Id
	content.NodeId = n.Id

	// 删除节点时，节点下不能有内容
	normalContentNum, err := content.CountNumUnderNode()
	if err != nil {
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if normalContentNum >= 1 {
		// 有内容，不能删除
		flog.Log.Errorf("DeleteNode err:%s", "has content child")
		resp.Error = Error(ContentNodeHasContentCanNotDelete, "")
		return
	}

	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	_, err = session.Exec("update fafacms_content_node SET sort_num=sort_num-1 where sort_num > ? and user_id = ?", n.SortNum, n.UserId)
	if err != nil {
		session.Rollback()
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	_, err = session.Where("id=?", n.Id).Delete(new(model.ContentNode))
	if err != nil {
		session.Rollback()
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
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
		resp.Error = Error(GetUserSessionError, err.Error())
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
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	resp.Data = n
	resp.Flag = true
}

type ListNodeRequest struct {
	Id              int      `json:"id"`
	Seo             string   `json:"seo" validate:"omitempty,alphanumunicode,gt=3,lt=30"`
	ParentNodeId    int      `json:"parent_node_id"`

	Level           int      `json:"level" validate:"oneof=-1 0 1"`
	UserId          int      `json:"user_id"`

	Sort            []string `json:"sort" validate:"dive,lt=100"`
	PageHelp
}

//
//type NodesRequest struct {
//	Id       int      `json:"id"`
//	UserId   int      `json:"user_id"`
//	UserName string   `json:"user_name"`
//	Seo      string   `json:"seo"`
//	Status          int      `json:"status" validate:"oneof=-1 0 1"`
//	CreateTimeBegin int64    `json:"create_time_begin"`
//	CreateTimeEnd   int64    `json:"create_time_end"`
//	UpdateTimeBegin int64    `json:"update_time_begin"`
//	UpdateTimeEnd   int64    `json:"update_time_end"`
//	Sort     []string `json:"sort" validate:"dive,lt=100"`
//}


type ListNodeResponse struct {
	Nodes []model.ContentNode `json:"nodes"`
	PageHelp
}

func ListNode(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("ListNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
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

// x->y
// x<y (x+1,y)-> -1
// x>y (y,x-1)-> +1
type SortNodeRequest struct {
	XID int `json:"xid"`
	YID int `json:"yid"`
	x   int `json:"x"`
	y   int `json:"y"`
}

//  拖曳排序超级函数
func SortNode(c *gin.Context) {
	resp := new(Resp)
	req := new(SortNodeRequest)
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
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	x := new(model.ContentNode)
	x.Id = req.XID
	x.UserId = uu.Id
	exist, err := x.GetSortOneNode()
	if err != nil {
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("SortNode err: %s", "x node not found")
		resp.Error = Error(ContentNodeNotFound, "x node not found")
		return
	}

	y := new(model.ContentNode)
	y.Id = req.YID
	y.UserId = uu.Id
	exist, err = y.GetSortOneNode()
	if err != nil {
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("SortNode err: %s", "y node not found")
		resp.Error = Error(ContentNodeNotFound, "y node not found")
		return
	}
	if y.ParentNodeId == x.Id {
		flog.Log.Errorf("SortNode err: %s", "can not move node to be his child's brother")
		resp.Error = Error(ContentNodeSortConflict, "can not move node to be his child's brother")
		return
	}

	req.x = x.SortNum
	req.y = y.SortNum

	children, err := x.CheckChildrenNum()
	if err != nil {
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if y.Level == 1 && children > 0 {
		flog.Log.Errorf("SortNode err: %s", "x has child can not move to be other's child's brother")
		resp.Error = Error(ContentNodeSortConflict, "x has child can not move to be other's child's brother")
		return
	}

	// x->y
	// x<y (x+1,y)-> -1
	// x>y (y,x-1)-> +1
	if req.x < req.y {
		session := config.FafaRdb.Client.NewSession()
		defer session.Close()

		err = session.Begin()
		if err != nil {
			flog.Log.Errorf("SortNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		_, err = session.Exec("update fafacms_content_node SET sort_num=sort_num-1 where sort_num > ? and sort_num <= ? and user_id = ?", req.x, req.y, uu.Id)
		if err != nil {
			session.Rollback()
			flog.Log.Errorf("SortNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		_, err = session.Exec("update fafacms_content_node SET sort_num=?,level=?,parent_node_id=? where id = ?", req.y, y.Level, y.ParentNodeId, x.Id)
		if err != nil {
			session.Rollback()
			flog.Log.Errorf("SortNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		err = session.Commit()
		if err != nil {
			session.Rollback()
			flog.Log.Errorf("SortNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	// x->y
	// x<y (x+1,y)-> -1
	// x>y (y,x-1)-> +1
	if req.x > req.y {
		session := config.FafaRdb.Client.NewSession()
		defer session.Close()

		err = session.Begin()
		if err != nil {
			flog.Log.Errorf("SortNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		_, err = session.Exec("update fafacms_content_node SET sort_num=sort_num+1 where sort_num < ? and sort_num >= ? and user_id = ?", req.x, req.y, uu.Id)
		if err != nil {
			session.Rollback()
			flog.Log.Errorf("SortNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		_, err = session.Exec("update fafacms_content_node SET sort_num=?,level=?,parent_node_id=? where id = ?", req.x, y.Level, y.ParentNodeId, x.Id)
		if err != nil {
			session.Rollback()
			flog.Log.Errorf("SortNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		err = session.Commit()
		if err != nil {
			session.Rollback()
			flog.Log.Errorf("SortNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	resp.Flag = true
	return
}
