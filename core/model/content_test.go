package model

import (
	"fmt"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/util/rdb"
	"testing"
)

func Testxx() {
	var err error
	c := rdb.MyDbConfig{}
	c.User = "root"
	c.Host = "127.0.0.1"
	c.Pass = "123456789"
	c.DriverName = "mysql"
	c.Prefix = "fafacms_"
	c.Name = "fafa"
	c.Debug = true
	db, err := rdb.NewDb(c)
	if err != nil {
		panic(err)
	}

	config.FafaRdb = db
}

func TestContentNode_CheckSeoValid(t *testing.T) {
	Testxx()
	n := new(ContentNode)
	n.Seo = "sss"
	n.UserId = 2
	n.Level = 1

	exist, err := n.CheckSeoValid()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(exist)

}

func TestContentNode_InsertOne(t *testing.T) {
	Testxx()

	n := new(ContentNode)
	n.UserId = 1
	n.Seo = ""
	n.Status = 1
	c, err := config.FafaRdb.Client.Cols("status").Cols("seo").Where("user_id=?", n.UserId).Update(n)
	fmt.Println(c, err)

}

func TestGroup_Delete(t *testing.T) {
	Testxx()

	g := new(Group)
	g.Id = 1000
	err := g.Delete()
	fmt.Println(err)
}

func TestResource_Get(t *testing.T) {
	Testxx()

	g := new(Resource)
	g.Id = 1000
	err := g.Get()
	fmt.Println(err)
}
