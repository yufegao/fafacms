package model

import (
	"github.com/hunterhug/fafacms/core/config"
	"testing"
)

func TestGroup_GetById(t *testing.T) {
	Testxx()

	g := new(Group)
	g.Id = 2
	config.FafaRdb.Client.Id(g.Id).Get(g)
}
