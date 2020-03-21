package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/missops/missops-go/api/models"
	"github.com/missops/missops-go/api/utils"
)

//UserCredential : request
type userCredential struct {
	Uname string `json:"user_name"`
	Pwd   string `json:"user_password"`
}

//createUserResponse : reponse
type createUserResponse struct {
	Success   bool   `json:"success"`
	Sessionid string `json:"session_id"`
}

//CreateUserHandler : handler for  user add
func CreateUserHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)

	ubody := &userCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, utils.ErrorRquestBodyParseFailed)
		return
	}
	if err := models.AddUserCredential(ubody.Uname, ubody.Pwd); err != nil {
		sendErrorResponse(w, utils.ErrorDBFailed)
		return
	}
	id := utils.GeneraterNewSessionID(ubody.Uname)
	resp := &createUserResponse{Success: true, Sessionid: id}

	if res, err := json.Marshal(resp); err != nil {
		sendErrorResponse(w, utils.ErrorInternalFault)
	} else {
		sendNormalResponse(w, 201, string(res))
	}

}

//LoginHandler ： login handler
func LoginHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := p.ByName("user_name")
	io.WriteString(w, uname)
}
