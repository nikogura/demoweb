package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/nikogura/demoweb/routers"
)

func main() {
	l := logs.NewLogger()
	l.SetLogger(logs.AdapterConsole)

	beego.Run()
}
