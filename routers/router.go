package routers

import (
	"github.com/astaxie/beego"
	"github.com/nikogura/demoweb/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/loggedin", &controllers.LoggedInController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/oauth-callback", &controllers.AuthController{})
	beego.Router("/logout", &controllers.LogoutController{})
}
