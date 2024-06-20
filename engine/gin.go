package engine

import (
	"github.com/gin-gonic/gin"
)

type GinEngine struct {
	*gin.Engine
}

func NewGinEngine(debug ...bool) (e *GinEngine) {
	if len(debug) > 0 {
		if debug[0] {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
	}
	e = &GinEngine{gin.New()}
	return
}

func (e *GinEngine) UseDefault() *GinEngine {
	e.Use(gin.Logger(), gin.Recovery())
	return e
}
