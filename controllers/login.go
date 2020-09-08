package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

//   res.redirect(`http://localhost:${config.fusionAuthPort}/oauth2/authorize?client_id=${config.clientID}&redirect_uri=${config.redirectURI}&response_type=code`);

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.Redirect(fmt.Sprintf("%s/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=openid offline_access", authHost, clientID, redirectUrl), 302)
}
