package engine

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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
