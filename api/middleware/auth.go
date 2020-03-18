package middleware

import (
	"net/http"

	"github.com/missops/missops-go/api/utils"
)

//ValidateUserSession : for middleware check session
func ValidateUserSession(r *http.Request) bool {
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
