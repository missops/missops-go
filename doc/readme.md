# Go从零开发一个web系列
## 准备Go环境
安装golang就不写了，这里使用到了gomod管理依赖包，需要Go的版本大于1.12,设置GOPATH和GOPROXY,并打开GO111MODULE
```
# 启用 Go Modules 功能
export GO111MODULE=on
# 配置 GOPROXY 环境变量
export GOPROXY=https://goproxy.io

go mod init
```

## 创建web目录结构
├── api
│   ├── handler
│   ├── main.go
├── doc
│   └── readme.md
└── go.mod
## 实现主路由
编译main.go,这里使用了httprouter这个库：https://github.com/julienschmidt/httprouter
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

## 实现用户登录功能
#### 准备数据库
```
#创建数据库
CREATE DATABASE IF NOT EXISTS missops default charset utf8 COLLATE utf8_general_ci;

#创建表，执行doc里面的missops.sql文件
```

#### 实现注册路由
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
#### 传参数路由
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

