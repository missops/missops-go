package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/missops/missops-go/api/utils"
)

//SendErrorResponse :  error response
func SendErrorResponse(w http.ResponseWriter, e utils.ErrorResponse) {
	w.WriteHeader(e.HttpSC)
	res, _ := json.Marshal(e.Error)
	io.WriteString(w, string(res))

}

//SendNormalResponse : normal response
func SendNormalResponse(w http.ResponseWriter, sc int, resp string) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
