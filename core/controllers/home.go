package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
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
	SortNum    int    `json:"sort_num"`
	Son        []Node
}

type NodesRequest struct {
	Id     int      `json:"id"`
	UserId int      `json:"user_id"`
	Seo    string   `json:"seo"`
	Sort   []string `json:"sort" validate:"dive,lt=100"`
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

	session := config.FafaRdb.Client.NewSession()
	defer session.Close()

	session.Table(new(model.ContentNode)).Where("1=1").And("status=?", 0)

	if req.UserId != 0 {
		session.And("user_id=?", req.UserId)
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

func UserCount(c *gin.Context) {

}

func Contents(c *gin.Context) {

}

func Content(c *gin.Context) {

}
