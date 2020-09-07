package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const AUTH_LANDING_TXT = "This is the authenticated landing page."

type LoggedInController struct {
	beego.Controller
}

func (c *LoggedInController) Get() {

	token := c.GetSession("token")

	logs.Debug("Token:")
	spew.Dump(token)

	if token != nil {
		tok, ok := token.(string)
		if ok {
			logs.Debug(fmt.Sprintf("Token: %s", tok))

			res, err := http.PostForm(fmt.Sprintf("%s/oauth2/introspect", authHost),
				url.Values{
					"client_id": {clientID},
					"token":     {tok},
				})

			defer res.Body.Close()

			resBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				err = errors.Wrapf(err, "failed reading response body")
				logs.Error(err)
				c.Abort("500")
				return
			}

			logs.Debug(fmt.Sprintf("%s", resBytes))

			responseData := make(map[string]interface{})

			err = json.Unmarshal(resBytes, &responseData)
			if err != nil {
				err = errors.Wrapf(err, "failed parsing response body")
				logs.Error(err)
				c.Abort("500")
				return
			}

			logs.Debug("introspection result:")
			spew.Dump(responseData)

			content := fmt.Sprintf("%s is logged in - %s", responseData["email"], AUTH_LANDING_TXT)

			c.Data["Content"] = content
			c.Data["Url"] = "http://localhost:9000/logout"
			c.Data["Label"] = "Logout"
			c.Data["Title"] = "Logged In Page"

			c.TplName = "index.tpl"
			return
		}
	}

	c.Redirect("http://localhost:9000/", 302)
}
