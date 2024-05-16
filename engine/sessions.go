package engine

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (e *GinEngine) UseSessions(store sessions.Store, name string, options sessions.Options, regStructs ...any) *GinEngine {

	if len(regStructs) > 0 {
		for s, _ := range regStructs {
			gob.Register(s)
		}
	}
	store.Options(options)
	e.Use(sessions.Sessions(name, store))
	return e
}

func GetSession(c *gin.Context) (session sessions.Session, err error) {
	session = sessions.Default(c)
	if nil == session {
		err = fmt.Errorf("session not fond")
		return
	}
	return
}
func SaveSession(c *gin.Context) (err error) {
	session := sessions.Default(c)
	if nil == session {
		return
	}
	session.Save()
	return
}
func ClearSession(c *gin.Context) (err error) {
	session := sessions.Default(c)
	if nil == session {
		return
	}
	session.Clear()
	session.Save()
	return
}
