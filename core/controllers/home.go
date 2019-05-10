package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/parrot/util"
	"math"
	"time"
)

func GetSecond2DateTimes(secord int64) string {
	tm := time.Unix(secord, 0)
	return tm.Format("2006-01-02 15:04:05")

}

func Home(c *gin.Context) {
	resp := new(Resp)
	resp.Flag = true
	resp.Data = "FaFa CMS: https://github.com/hunterhug/fafacms"
	defer func() {
		c.JSON(200, resp)
	}()
}

type People struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	Name       string `json:"name" xorm:"varchar(100) notnull unique"`  // 独一无二的标志
	NickName   string `json:"nick_name" xorm:"varchar(100) notnull"`    // 昵称，如小花花，随便改
	Email      string `json:"email" xorm:"varchar(100) notnull unique"` // 邮箱，独一无二
	WeChat     string `json:"wechat" xorm:"varchar(100)"`
	WeiBo      string `json:"weibo" xorm:"TEXT"`
	Github     string `json:"github" xorm:"TEXT"`
	QQ         string `json:"qq" xorm:"varchar(100)"`
	Gender     int    `json:"gender" xorm:"not null comment('0 unknow,1 boy,2 girl') TINYINT(1)"`
	Describe   string `json:"describe" xorm:"TEXT"`
	HeadPhoto  string `json:"head_photo" xorm:"varchar(700)"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time,omitempty"`
}

type PeoplesRequest struct {
	Sort []string `json:"sort" validate:"dive,lt=100"`
	PageHelp
}

type PeoplesResponse struct {
	Users []People `json:"users"`
	PageHelp
}

func Peoples(c *gin.Context) {
	resp := new(Resp)

	defer func() {
		JSON(c, 200, resp)
	}()

	respResult := new(PeoplesResponse)
	req := new(PeoplesRequest)
	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	session.Table(new(model.User)).Where("1=1").And("status=?", 1)

	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		flog.Log.Errorf("ListUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	users := make([]model.User, 0)
	p := &req.PageHelp
	if total == 0 {
	} else {
		p.build(session, req.Sort, model.UserSortName)
		err = session.Find(&users)
		if err != nil {
			flog.Log.Errorf("ListUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	peoples := make([]People, 0, len(users))
	for _, v := range users {
		p := People{}
		p.Id = v.Id
		p.Describe = v.Describe
		p.CreateTime = GetSecond2DateTimes(v.CreateTime)

		if v.UpdateTime > 0 {
			p.UpdateTime = GetSecond2DateTimes(v.UpdateTime)
		}

		p.Email = v.Email
		p.Github = v.Github
		p.Name = v.Name
		p.NickName = v.NickName
		p.HeadPhoto = v.HeadPhoto
		p.QQ = v.QQ
		p.WeChat = v.WeChat
		p.WeiBo = v.WeiBo
		p.Gender = v.Gender
		peoples = append(peoples, p)
	}
	respResult.Users = peoples
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

type Node struct {
	Id         int    `json:"id"`
	Seo        string `json:"seo"`
	Name       string `json:"name"`
	Describe   string `json:"describe"`
	ImagePath  string `json:"image_path"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time,omitempty"`
	UserId     int    `json:"user_id"`
	UserName   string `json:"user_name"`
	SortNum    int    `json:"sort_num"`
	Son        []Node
}

type NodesRequest struct {
	Id       int      `json:"id"`
	UserId   int      `json:"user_id"`
	UserName string   `json:"user_name"`
	Seo      string   `json:"seo"`
	Sort     []string `json:"sort" validate:"dive,lt=100"`
}

type NodesResponse struct {
	Nodes []Node `json:"nodes"`
}

func Nodes(c *gin.Context) {
	resp := new(Resp)

	defer func() {
		JSON(c, 200, resp)
	}()

	respResult := new(NodesResponse)
	req := new(NodesRequest)
	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.UserId == 0 && req.UserName == "" {
		flog.Log.Errorf("ListNode err:%s", "")
		resp.Error = Error(ParasError, "where is empty")
		return
	}

	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	session.Table(new(model.ContentNode)).Where("1=1").And("status=?", 0).Asc("sort_num").Asc("create_time")

	if req.UserId != 0 {
		session.And("user_id=?", req.UserId)
	}

	if req.UserName != "" {
		session.And("user_name=?", req.UserName)
	}

	if req.Id != 0 {
		session.And("id=?", req.Id)
	}

	if req.Seo != "" {
		session.And("seo=?", req.Seo)
	}

	nodes := make([]model.ContentNode, 0)
	Build(session, req.Sort, model.ContentNodeSortName)
	err := session.Find(&nodes)
	if err != nil {
		flog.Log.Errorf("ListNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	father := make([]model.ContentNode, 0)
	son := make([]model.ContentNode, 0)
	for _, v := range nodes {
		if v.Level == 0 {
			father = append(father, v)
		} else {
			son = append(son, v)
		}
	}

	n := make([]Node, 0)
	for _, v := range father {
		f := Node{}
		f.Id = v.Id
		f.Seo = v.Seo
		f.Describe = v.Describe
		f.ImagePath = v.ImagePath
		f.Name = v.Name
		if v.UpdateTime > 0 {
			f.UpdateTime = GetSecond2DateTimes(v.UpdateTime)
		}
		f.CreateTime = GetSecond2DateTimes(v.CreateTime)
		f.SortNum = v.SortNum
		f.UserName = v.UserName
		f.UserId = v.UserId
		for _, vv := range son {
			if vv.ParentNodeId == f.Id {
				s := Node{}
				s.Id = vv.Id
				s.Seo = vv.Seo
				s.Describe = vv.Describe
				s.ImagePath = vv.ImagePath
				s.Name = vv.Name
				if vv.UpdateTime > 0 {
					s.UpdateTime = GetSecond2DateTimes(vv.UpdateTime)
				}
				s.CreateTime = GetSecond2DateTimes(vv.CreateTime)
				s.SortNum = vv.SortNum
				s.UserId = vv.UserId
				s.UserName = v.UserName
				f.Son = append(f.Son, s)
			}
		}

		n = append(n, f)

	}

	respResult.Nodes = n
	resp.Flag = true
	resp.Data = respResult
}

type UserInfoRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func UserInfo(c *gin.Context) {
	resp := new(Resp)

	defer func() {
		JSON(c, 200, resp)
	}()

	req := new(UserInfoRequest)
	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.Id == 0 && req.Name == "" {
		resp.Error = Error(ParasError, "where is empty")
		return
	}

	user := new(model.User)
	user.Id = req.Id
	user.Name = req.Name
	exist, err := config.FafaRdb.Client.Where("status=?", 1).Get(user)
	if err != nil {
		flog.Log.Errorf("UserInfo err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("UserInfo err:%s", "Not exist")
		resp.Error = Error(DbNotFound, "not exist")
		return
	}

	v := user
	p := People{}
	p.Id = v.Id
	p.Describe = v.Describe
	p.CreateTime = GetSecond2DateTimes(v.CreateTime)

	if v.UpdateTime > 0 {
		p.UpdateTime = GetSecond2DateTimes(v.UpdateTime)
	}

	p.Email = v.Email
	p.Github = v.Github
	p.Name = v.Name
	p.NickName = v.NickName
	p.HeadPhoto = v.HeadPhoto
	p.QQ = v.QQ
	p.WeChat = v.WeChat
	p.WeiBo = v.WeiBo
	p.Gender = v.Gender

	resp.Flag = true
	resp.Data = p

}

type UserCountRequest struct {
	UserId   int    `json:"id"`
	UserName string `json:"user_name"`
}

type UserCountX struct {
	Count           int    `json:"count"`
	Days            string `json:"days"`
	CreateTimeBegin int64  `json:"create_time_begin"`
	CreateTimeEnd   int64  `json:"create_time_end"`
}
type UserCountResponse struct {
	Info     []UserCountX `json:"info"`
	UserId   int          `json:"user_id"`
	UserName string       `json:"user_name"`
}

func UserCount(c *gin.Context) {
	resp := new(Resp)

	defer func() {
		JSON(c, 200, resp)
	}()

	req := new(UserCountRequest)
	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.UserId == 0 && req.UserName == "" {
		resp.Error = Error(ParasError, "where is empty")
		return
	}

	user := new(model.User)
	user.Id = req.UserId
	user.Name = req.UserName
	user.Status = 1
	err := user.Get()
	if err != nil {
		flog.Log.Errorf("UserCount err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if user.Status != 1 {
		flog.Log.Errorf("UserCount err:%s", "not activate")
		resp.Error = Error(DBError, "not activate")
		return
	}

	req.UserId = user.Id

	sql := "SELECT DATE_FORMAT(from_unixtime(create_time),'%Y%m%d') days,count(id) count FROM `fafacms_content` WHERE user_id=? and version>0 and status=0 group by days;"
	result, err := config.FafaRdb.Client.QueryString(sql, req.UserId)
	if err != nil {
		flog.Log.Errorf("UserCount err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	back := make([]UserCountX, 0)
	for _, v := range result {
		t := UserCountX{}
		t.Count, _ = util.SI(v["count"])
		t.Days = v["days"]
		begin, _ := time.Parse("20060102", t.Days)
		end := begin.AddDate(0, 0, 1)
		t.CreateTimeBegin = begin.Unix()
		t.CreateTimeEnd = end.Unix()
		back = append(back, t)
	}

	resp.Flag = true
	resp.Data = UserCountResponse{
		Info:     back,
		UserId:   user.Id,
		UserName: user.Name,
	}
}

type ContentsRequest struct {
	NodeId          int      `json:"node_id"`
	NodeSeo         string   `json:"node_seo"`
	UserId          int      `json:"user_id"`
	UserName        string   `json:"user_name"`
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	Sort            []string `json:"sort" validate:"dive,lt=100"`
	PageHelp
}

type ContentsX struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	Seo        string `json:"seo" xorm:"index"`
	Title      string `json:"title" xorm:"varchar(200) notnull"`
	UserId     int    `json:"user_id" xorm:"bigint index"` // 内容所属用户
	UserName   string `json:"user_name" xorm:"index"`
	NodeId     int    `json:"node_id" xorm:"bigint index"`                                     // 节点ID
	NodeSeo    string `json:"node_seo" xorm:"index"`                                           // 节点ID SEO
	Top        int    `json:"top" xorm:"not null comment('0 normal, 1 top') TINYINT(1) index"` // 置顶
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time,omitempty"`
	ImagePath  string `json:"image_path" xorm:"varchar(700)"`
	Views      int    `json:"views"` // 被点击多少次，弱化
	IsLock     bool   `json:"is_lock"`
	Describe   string `json:"describe"`
}

type ContentsResponse struct {
	Contents []ContentsX `json:"contents"`
	PageHelp
}

func Contents(c *gin.Context) {
	resp := new(Resp)

	respResult := new(ContentsResponse)
	req := new(ContentsRequest)
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
		flog.Log.Errorf("Contents err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.Content)).Where("1=1")

	if req.UserId != 0 {
		session.And("user_id=?", req.UserId)
	}

	if req.UserName != "" {
		session.And("user_name=?", req.UserName)
	}

	session.And("status=?", 0).And("version>?", 0)

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

	// count num
	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		flog.Log.Errorf("Contents err:%s", err.Error())
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
			flog.Log.Errorf("Contents err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	// result
	bcs := make([]ContentsX, 0, len(cs))
	for _, c := range cs {
		temp := ContentsX{}
		temp.UserId = c.UserId
		temp.Seo = c.Seo
		temp.NodeSeo = c.NodeSeo
		temp.UserName = c.UserName
		temp.Id = c.Id
		temp.Top = c.Top
		temp.Title = c.Title
		temp.NodeId = c.NodeId
		temp.Views = c.Views
		temp.CreateTime = GetSecond2DateTimes(c.CreateTime)
		temp.UpdateTime = GetSecond2DateTimes(c.UpdateTime)
		temp.ImagePath = c.ImagePath

		if c.Password != "" {
			temp.IsLock = true
		}
		bcs = append(bcs, temp)
	}

	respResult.Contents = bcs
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

type ContentRequest struct {
	Id       int    `json:"id"`
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Seo      string `json:"seo"`
	Password string `json:"password"`
}

func Content(c *gin.Context) {
	resp := new(Resp)
	req := new(ContentRequest)
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
	content.UserId = req.UserId
	content.UserName = req.UserName
	content.Seo = req.Seo
	exist, err := content.GetByRaw()
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

	if content.Status == 0 {

	} else if content.Status == 2 {
		flog.Log.Errorf("TakeContent err: %s", "content ban")
		resp.Error = Error(DbNotFound, "content ban")
		return
	} else {
		flog.Log.Errorf("TakeContent err: %s", "content not found")
		resp.Error = Error(DbNotFound, "content not found")
		return
	}

	if content.Password != "" && content.Password != req.Password {
		flog.Log.Errorf("TakeContent err: %s", "content password")
		resp.Error = Error(DbNotFound, "content password")
		return
	}

	cx := content
	temp := ContentsX{}
	temp.UserId = cx.UserId
	temp.Seo = cx.Seo
	temp.NodeSeo = cx.NodeSeo
	temp.UserName = cx.UserName
	temp.Id = cx.Id
	temp.Top = cx.Top
	temp.Title = cx.Title
	temp.NodeId = cx.NodeId
	temp.Views = cx.Views
	temp.CreateTime = GetSecond2DateTimes(cx.CreateTime)
	temp.UpdateTime = GetSecond2DateTimes(cx.UpdateTime)
	temp.ImagePath = cx.ImagePath

	if cx.Password != "" {
		temp.IsLock = true
	}

	temp.Describe = cx.Describe

	cx.UpdateView()

	resp.Flag = true
	resp.Data = temp

}
