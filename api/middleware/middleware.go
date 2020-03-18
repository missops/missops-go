package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middlewareHandler struct {
	r *httprouter.Router
}

//NewMiddlewareHandler : new a middleware
func NewMiddlewareHandler(r *httprouter.Router) http.Handler {
	m := middlewareHandler{}
	m.r = r
	return m
}

func (m middlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//user auth
	ValidateUserSession(r)
	//user auth
	//log
	m.r.ServeHTTP(w, r)
}
