package engine

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
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

func (e *GinEngine) UseSessions(store sessions.Store, name string, regStructs ...any) *GinEngine {
	e.Use(sessions.Sessions(name, store))
	if len(regStructs) > 0 {
		for s, _ := range regStructs {
			gob.Register(s)
		}
	}
	return e
}

func GetSessions(c *gin.Context) (session sessions.Session, err error) {
	ginSession, exists := c.Get(sessions.DefaultKey)
	if !exists {
		err = fmt.Errorf("session not enabled")
		return
	}
	session, ok := ginSession.(sessions.Session)
	if !ok {
		err = fmt.Errorf("session type error")
	}
	return
}
