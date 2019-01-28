package model

type Resource struct {
	Id       int    `json:"id" xorm:"bigint pk autoincr"`
	Name     string `json:"name,omitempty"`
	Url      string `json:"url" xorm:"TEXT`
	Describe string `json:"describe,omitempty" xorm:"TEXT`
}

type GroupResource struct {
	Id         int `json:"id" xorm:"bigint pk autoincr"`
	GroupId    int `json:"group_id"`
	ResourceId int `json:"resource_id"`
}
