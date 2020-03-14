package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/missops/missops-go/api/handler"
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

	http.ListenAndServe(":8080", r)

}
