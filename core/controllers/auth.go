package controllers

import (
	"github.com/gin-gonic/gin"
	"sync"
)

var (
	AuthResource = sync.Map{}
)

// url===>resource_id
func InitAuthResource() {

}

// filter
var AuthFilter = func(c *gin.Context) {
}
