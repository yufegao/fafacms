package server

import "github.com/gin-gonic/gin"

func Server() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	return r
}
