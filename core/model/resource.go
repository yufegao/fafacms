package model

type Resource struct {
	Id       int    `json:"id" xorm:"bigint pk autoincr"`
	Name     string `json:"name,omitempty"`
	Url      string `json:"url" xorm:"TEXT`
	Describe string `json:"describe,omitempty" xorm:"TEXT`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}

// extern: Group --> Resource
type GroupResource struct {
	Id         int `json:"id" xorm:"bigint pk autoincr"`
	GroupId    int `json:"group_id"`
	ResourceId int `json:"resource_id"`
}
