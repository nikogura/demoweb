package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

// Get Responds to GET on /login.  Redirects to FA.  Note the 'scope' parameters.  Without these we don't get ID or refresh tokens from FA.
func (c *LoginController) Get() {
	c.Redirect(fmt.Sprintf("%s/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=openid offline_access", authHost, clientID, redirectUrl), 302)
}
