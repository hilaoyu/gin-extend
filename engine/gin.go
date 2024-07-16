package engine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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

func (e *GinEngine) UseRequestLogger(logger io.Writer, notLogged ...string) *GinEngine {
	e.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: gin.LogFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%v | %3d | %13v | %15s | %-7s %#v\n%s",
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
		}),
		Output:    logger,
		SkipPaths: notLogged,
	}))
	return e
}
