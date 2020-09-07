package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

const UNAUTH_LANDING_TXT = "This is the unauthenticated landing page."

func (c *MainController) Get() {
	c.Data["Content"] = UNAUTH_LANDING_TXT
	c.Data["Url"] = "http://localhost:9000/login"
	c.Data["Label"] = "Login"
	c.Data["Title"] = "LandingPage"

	c.TplName = "index.tpl"
}
