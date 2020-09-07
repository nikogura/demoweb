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
	c.DelSession("token")
	// perform logout
	c.Redirect(fmt.Sprintf("%s/oauth2/logout?client_id=%s", authHost, clientID), 302)
}
