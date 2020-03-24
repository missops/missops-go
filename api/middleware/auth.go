package middleware

import (
	"net/http"

	"github.com/missops/missops-go/api/utils"
)

//validateUserSession : for middleware check session
func validateUserSession(r *http.Request) bool {
	sid := r.Header.Get("X-Session-ID")
	if len(sid) == 0 {
		return false
	}
	uname, ok := utils.IsSessionExpired(sid)
	if ok {
		return false
	}
	r.Header.Add("X-Session-Name", uname)
	return true
}

//validateUser: user is validate
func validateUser(r *http.Request) bool {
	return false
}
