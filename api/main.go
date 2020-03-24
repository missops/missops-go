package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/missops/missops-go/api/handler"
	"github.com/missops/missops-go/api/middleware"
)

//RegisterHandlers is httprouter.Router
func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.POST("/user", handler.CreateUserHandler)
	router.POST("/user/:user_name", handler.LoginHandler)
	return router
}

func main() {
	r := RegisterHandlers()
	mh := middleware.NewMiddlewareHandler(r)
	http.ListenAndServe(":8080", mh)
}
