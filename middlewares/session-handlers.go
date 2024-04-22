package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hilaoyu/gin-extend/engine"
)

func SessionHandler(before func(s sessions.Session, c *gin.Context), after func(s sessions.Session, c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {

		session, err := engine.GetSessions(c)
		if nil != err {
			c.AbortWithError(500, err)
		}
		before(session, c)
		c.Next()
		// 请求后
		after(session, c)
	}
}
