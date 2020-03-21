package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/missops/missops-go/api/utils"
)

func sendErrorResponse(w http.ResponseWriter, e utils.ErrorResponse) {
	w.WriteHeader(e.HttpSC)
	res, _ := json.Marshal(e.Error)
	io.WriteString(w, string(res))

}

func sendNormalResponse(w http.ResponseWriter, sc int, resp string) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}