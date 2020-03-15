package utils

//Err : error message and error code
type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

//ErrorResponse :  http statuscode and error
type ErrorResponse struct {
	HttpSC int
	Error  Err
}

var (
	ErrorRquestBodyParseFailed = ErrorResponse{HttpSC: 400, Error: Err{Error: "request body parse failed.", ErrorCode: "001"}}
	ErrorAuthFailed            = ErrorResponse{HttpSC: 401, Error: Err{Error: "auth failed.", ErrorCode: "002"}}
)
