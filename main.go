package main

import (
	_ "MyWebApi/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"./logs/MyWebApi.log","separate":["error"]}`)
	beego.Run()
}
