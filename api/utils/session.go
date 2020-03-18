package utils

import (
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

//Session :  session struct
type Session struct {
	Name string
	TTL  int64
}

var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

//deleteExpiredSession : delete expired session id
func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
}

//GeneraterNewSessionID : make new session id
func GeneraterNewSessionID(uname string) string {
	id := uuid.NewV4().String()
	ct := time.Now().UnixNano() / 1000000
	ttl := ct + 30*60*1000 //30 min

	ss := &Session{
		Name: uname,
		TTL:  ttl,
	}
	sessionMap.Store(id, ss)
	return id
}

//IsSessionExpired : session id and ttl is not Expired
func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		ct := time.Now().UnixNano() / 1000000
		if ss.(*Session).TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}
		return ss.(*Session).Name, false

	}
	return "", true
}
