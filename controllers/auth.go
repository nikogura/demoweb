package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/session"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const clientID = "b843a55c-4f3f-43d0-b176-674abcbeffb5"
const authHost = "http://localhost:9011"
const apiKey = "ToPgrT9sJcfzSALDIlepaFZi6rtjl3TZaRmhkDqqZ4M"
const redirectUrl = "http://localhost:9000/oauth-callback"
const clientSecret = "LoYgAuihAfUPuTYmYQOMHR1nOpxSgiJhtEGt5GWS9Cs"

var globalSessions *session.Manager

type AuthController struct {
	beego.Controller
}

func (c *AuthController) Get() {
	logs.Debug("Auth controller firing")

	res, err := http.PostForm(fmt.Sprintf("%s/oauth2/token", authHost),
		url.Values{
			"client_id":     {clientID},
			"client_secret": {clientSecret},
			"code":          {c.GetString("code")},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {redirectUrl},
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

	// save token to session
	c.SetSession("token", responseData["access_token"])

	fmt.Printf("Response from FA:")
	spew.Dump(responseData)

	c.Redirect("/loggedin", 302)
}
