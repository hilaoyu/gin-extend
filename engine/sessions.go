package engine

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	gorillaSessions "github.com/gorilla/sessions"
	"net/http"
)

const (
	GinExtendSessionKey = "github.com/hilaoyu/gin-extend/engine"
)

func (e *GinEngine) UseSessions(store sessions.Store, name string, options sessions.Options, regStructs ...any) *GinEngine {

	if len(regStructs) > 0 {
		for s, _ := range regStructs {
			gob.Register(s)
		}
	}
	store.Options(options)
	e.Use(func(c *gin.Context) {
		s := &ginExtendSession{name, c.Request, store, nil, false, c.Writer}
		c.Set(GinExtendSessionKey, s)
		defer context.Clear(c.Request)
		c.Next()
	})
	return e
}

func GetSession(c *gin.Context) (session sessions.Session, err error) {
	session = c.MustGet(GinExtendSessionKey).(sessions.Session)
	if nil == session {
		err = fmt.Errorf("session not fond")
		return
	}
	return
}
func SaveSession(c *gin.Context) (err error) {
	session, _ := GetSession(c)
	if nil == session {
		return
	}
	session.Save()
	return
}
func ClearSession(c *gin.Context) (err error) {
	session, _ := GetSession(c)
	if nil == session {
		return
	}
	session.Clear()
	session.Save()
	return
}

type ginExtendSession struct {
	name    string
	request *http.Request
	store   sessions.Store
	session *gorillaSessions.Session
	written bool
	writer  http.ResponseWriter
}

func (s *ginExtendSession) ID() string {
	gorillaSession, _ := s.Session()
	if nil != gorillaSession {
		return gorillaSession.ID
	}
	return ""
}

func (s *ginExtendSession) Get(key interface{}) interface{} {
	gorillaSession, _ := s.Session()
	if nil != gorillaSession {
		return gorillaSession.Values[key]
	}
	return nil
}

func (s *ginExtendSession) Set(key interface{}, val interface{}) {
	gorillaSession, _ := s.Session()
	if nil != gorillaSession {
		gorillaSession.Values[key] = val
	}
	s.written = true
}

func (s *ginExtendSession) Delete(key interface{}) {
	gorillaSession, _ := s.Session()
	if nil != gorillaSession {
		delete(gorillaSession.Values, key)
	}
	s.written = true
}

func (s *ginExtendSession) Clear() {
	gorillaSession, _ := s.Session()
	if nil != gorillaSession {
		for key := range gorillaSession.Values {
			s.Delete(key)
		}
	}

}

func (s *ginExtendSession) AddFlash(value interface{}, vars ...string) {
	gorillaSession, _ := s.Session()
	if nil != gorillaSession {
		gorillaSession.AddFlash(value, vars...)
	}
	s.written = true
}

func (s *ginExtendSession) Flashes(vars ...string) []interface{} {
	s.written = true
	gorillaSession, _ := s.Session()
	if nil != gorillaSession {
		return gorillaSession.Flashes(vars...)
	}
	return nil
}

func (s *ginExtendSession) Options(options sessions.Options) {
	s.written = true
	gorillaSession, _ := s.Session()
	if nil != gorillaSession {
		gorillaSession.Options = options.ToGorillaOptions()
	}
}

func (s *ginExtendSession) Save() error {
	if s.Written() {
		gorillaSession, _ := s.Session()
		if nil != gorillaSession {
			e := gorillaSession.Save(s.request, s.writer)
			if e == nil {
				s.written = false
			}
			return e
		}

	}
	return nil
}

func (s *ginExtendSession) Session() (session *gorillaSessions.Session, err error) {
	if s.session == nil {
		s.session, err = s.store.Get(s.request, s.name)
		if nil != err {
			return
		}
	}
	session = s.session
	return
}

func (s *ginExtendSession) Written() bool {
	return s.written
}
