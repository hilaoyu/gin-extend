package engine

import (
	"github.com/gin-gonic/gin"
)

type GinEngine struct {
	*gin.Engine
}

func NewGinEngine() (e *GinEngine) {
	e = &GinEngine{gin.New()}
	return
}

func (e *GinEngine) UseDefault() *GinEngine {
	e.Use(gin.Logger(), gin.Recovery())
	return e
}
func (e *GinEngine) Debug(debug bool) *GinEngine {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	return e
}
