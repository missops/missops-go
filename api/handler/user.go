package handler

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//CreateUserHandler : handler for  user add
func CreateUserHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	io.WriteString(w, "ok")
}

//LoginHandler ： login handler
func LoginHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := p.ByName("user_name")
	io.WriteString(w, uname)
}
