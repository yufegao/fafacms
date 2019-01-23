package model

type User struct {
	Id        string `json:"id,omitempty" xorm:"bigint pk autoincr"`
	Name      string `json:"name,omitempty" xorm:"varchar(100) notnull"`
	NickName  string `json:"nick_name,omitempty" xorm:"varchar(100) notnull"`
	Email     string `json:"email,omitempty" xorm:"varchar(100)"`
	WeChat    string `json:"wechat,omitempty" xorm:"varchar(100)"`
	WeiBo     string `json:"weibo,omitempty" xorm:"TEXT`
	QQ        string `json:"qq,omitempty" xorm:"varchar(100)"`
	Password  string `json:"password,omitempty" xorm:"varchar(100)`
	Gender    int    `json:"gender,omitempty" xorm:"not null comment('1 boy，2 girl') TINYINT(1)"`
	Says      string `json:"says,omitempty" xorm:"TEXT`
	ImagePath string `json:"image_path" xorm:"TEXT`
}
