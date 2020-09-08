package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type LogoutController struct {
	beego.Controller
}

func (c *LogoutController) Get() {
	logs.Debug("Logout controller firing")
	// remove app session token
	c.DelSession("fatoken")
	// perform logout on FA
	c.Redirect(fmt.Sprintf("%s/oauth2/logout?client_id=%s", beego.AppConfig.String("authHost"), beego.AppConfig.String("clientId")), 302)
}
