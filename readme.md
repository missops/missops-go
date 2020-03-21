# Go从零开发一个web系列
[TOC]

## 准备Go环境
安装golang就不写了，这里使用到了gomod管理依赖包，需要Go的版本大于1.12,设置GOPATH和GOPROXY,并打开GO111MODULE
```
# 启用 Go Modules 功能
export GO111MODULE=on
# 配置 GOPROXY 环境变量
export GOPROXY=https://goproxy.io

go mod init
```

## web目录结构
![](http://img.hixuxu.com/2020-03-21-043439.png)

## 实现主路由
编译main.go,这里使用了[httprouter](https://github.com/julienschmidt/httprouter)这个库。
```
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/missops/missops-go/api/handler"
)

//RegisterHandlers is httprouter.Router
func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	return router
}

func main() {
	r := RegisterHandlers()
	http.ListenAndServe(":8080", r)

}
```

## 实现注册路由
在main.go添加user路由
```
//RegisterHandlers is httprouter.Router
func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.POST("/user", handler.CreateUserHandler)
	return router
}
```
在hanlder文件下创建文件user.go，里面添加CreateUserHandler
```
package handler

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//CreateUserHandler : handler for  user add
func CreateUserHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	io.WriteString(w, "hello missops ! ")
}
}
```
## 传参数路由
在main.go内添加带user_name参数的路由
```
	router.POST("/user/:user_name", handler.LoginHandler)
```
在handler内添加LoginHandler
```
//LoginHandler ： login handler
func LoginHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := p.ByName("user_name")
	io.WriteString(w, uname)
}
```
使用Postman带参数POST请求接口
![](http://img.hixuxu.com/2020-03-14-042314.jpg)

## 错误处理
校验错误的情况下需要进行error处理，准备error.go文件，我们创建一个Err和ErrorResponse，结构化错误。
```
type ErrorResponse struct {
	HttpSC  int    //http status code
	Error Err      // error message and error code 
}
```
其中的Err结构体如下
```
type Err struct {
	Error string `json:"error"`
	ErrorCode string `json:"error_code"`
}
```
定义错误类型的response变量
```
var (
	ErrorRquestBodyParseFailed = ErrorResponse{HttpSC: 400, Error: Err{Error: "request body parse failed.", ErrorCode: "001"}}
	ErrorAuthFailed            = ErrorResponse{HttpSC: 401, Error: Err{Error: "auth failed.", ErrorCode: "002"}}
)
```
## 连接mysql数据库
准备数据库，创建数据库表。
```
#创建数据库
CREATE DATABASE IF NOT EXISTS missops default charset utf8 COLLATE utf8_general_ci;

#创建表，执行doc里面的missops.sql文件
```
创建db_connect.go文件,适应init函数
```
package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	dbConn, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/missops?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
}
```
## 数据库操作
创建数据库操作,以增删改查用户表为例
```
package models

//AddUserCredential : insert user to databases
func AddUserCredential(userName string, pwd string) error {
	stmtIn, err := dbConn.Prepare("INSERT INTO missops_user (user_name,user_pwd) VALUES (?,?)")
	if err != nil {
		return err
	}
	defer stmtIn.Close()
	stmtIn.Exec(userName, pwd)

	return nil
}

//GetUserCredential : select pwd from databases
func GetUserCredential(userName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT user_pwd FROM missops_user WHERE user_name = ?  ")
	if err != nil {
		return "", err
	}
	defer stmtOut.Close()
	var pwd string
	stmtOut.QueryRow(userName).Scan(&pwd)
	return pwd, nil
}

//DeleteUser : delete user
func DeleteUser(userName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM  missops_user WHERE user_name = ? and user_pwd = ?")
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	stmtDel.Exec(userName, pwd)
	return nil
}

```
## 测试用例
创建针对上面用户增删改查的测试用例
```
package models

import (
	"testing"
)

func clearTables() {
	dbConn.Exec("truncate missops_user")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}
func TestUserWork(t *testing.T) {
	t.Run("ADD", testAddUser)
	t.Run("GET", testGetUser)
	t.Run("DEL", testDelUser)
	t.Run("REGET", testRegetUser)
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("missops", "123456")
	if err != nil {
		t.Errorf("Error of user add: %v ", err)
	}
}

func testGetUser(t *testing.T) {
	pwd, err := GetUserCredential("missops")
	if pwd != "123456" || err != nil {
		t.Errorf("Error of user get: %v ", err)
	}

}

func testDelUser(t *testing.T) {
	err := DeleteUser("missops", "123456")
	if err != nil {
		t.Errorf("Error of user del: %v ", err)
	}
}

func testRegetUser(t *testing.T) {
	_, err := GetUserCredential("missops")
	if err != nil {
		t.Errorf("Error of user get: %v ", err)
	}
}

```
## session设置
先定义一个session结构，包含name和ttl
```
//Session :  session struct
type Session struct {
	Name string
	TTL   int64
}
```
使用sync.Map存入内存
```
var sessionMap *sync.Map

func init(){
	sessionMap = &sync.Map{}
}
```
增加创建session和检查session是否过期的方法
```
//deleteExpiredSession : delete expired session id
func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
}

//GeneraterNewSessionID : make new session id
func GeneraterNewSessionID(uname string) string {
	id := uuid.NewV4().String()
	ct := time.Now().UnixNano() / 1000000
	ttl := ct + 30*60*1000 //30 min

	ss := &Session{
		Name: uname,
		TTL:  ttl,
	}
	sessionMap.Store(id, ss)
	return id
}

//IsSessionExpired : session id and ttl is not Expired
func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		ct := time.Now().UnixNano() / 1000000
		if ss.(*Session).TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}
		return ss.(*Session).Name, false

	}
	return "", true
}
```
## middleware设置
创建middlewarehandler，需要实现http.Handler的ServerHTTP方法
```
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
	//Check session
	//log
	m.r.ServeHTTP(w, r)
}

```
将middlewarehandler注册到ListenAndServe
```
	r := RegisterHandlers()
	mh := middleware.NewMiddlewareHandler(r)
	http.ListenAndServe(":8080", mh)
```
创建auth.go丰富middleware处理内容
```
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
```
## response定义
定义sendErrorResponse和sendNormalResponse 
```
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
```
## handler补全
以CreateUserHandler为例，先定义输入输出struct
```
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

```
CreateUserHandler
```
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

```
Postman请求接口

![](http://img.hixuxu.com/2020-03-21-043258.png)

