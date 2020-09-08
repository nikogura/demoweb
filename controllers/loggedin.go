package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const AUTH_LANDING_TXT = "This is the authenticated landing page."

type LoggedInController struct {
	beego.Controller
}

func (c *LoggedInController) Get() {
	fatoken := c.GetSession("fatoken")
	logs.Debug("Time is %d", time.Now().Unix())

	baseurl, err := url.Parse(beego.AppConfig.String("authHost"))
	if err != nil {
		logs.Error("failed to parse %s as a URL:%s", beego.AppConfig.String("authHost"), err)
		c.Abort("500")
	}
	httpClient := http.Client{}
	var auth = fusionauth.NewClient(&httpClient, baseurl, beego.AppConfig.String("apiKey"))

	tok, ok := fatoken.(*fusionauth.AccessToken)
	if ok {
		logs.Debug(fmt.Sprintf("User ID: %s", tok.UserId))

		resp, err := auth.ValidateJWT(tok.AccessToken)
		if err != nil {
			logs.Error("Error validating JWT: %s", err)
			c.Redirect("http://localhost:9000/", 302)
		}

		if resp.StatusCode != 200 {
			logs.Debug("Token expired.  Attempting refresh.")
			accessToken, oauthErr, err := auth.ExchangeRefreshTokenForAccessToken(tok.RefreshToken, beego.AppConfig.String("clientId"), beego.AppConfig.String("clientSecret"), "", "")

			if err != nil {
				logs.Error("Error exchanging access code for token: %s", err)
				// kill the app session
				c.DelSession("fatoken")
				c.Abort("500")
			}
			if oauthErr != nil {
				logs.Error("Oauth Error in refresh attempt: %s", oauthErr.Error)
				spew.Dump(oauthErr)
				// kill the app session
				c.DelSession("fatoken")
				c.Redirect("http://localhost:9000/", 302)
				//c.Abort("500")
			}

			logs.Debug("Token Refreshed")
			// set token object as the new token
			tok = accessToken
			// persist it to the app session
			c.SetSession("fatoken", accessToken)
		}

		res, err := http.PostForm(fmt.Sprintf("%s/oauth2/introspect", beego.AppConfig.String("authHost")),
			url.Values{
				"client_id": {beego.AppConfig.String("clientId")},
				"token":     {tok.AccessToken},
			})

		defer res.Body.Close()

		resBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			err = errors.Wrapf(err, "failed reading response body")
			logs.Error(err)
			c.Abort("500")
			return
		}

		responseData := make(map[string]interface{})

		err = json.Unmarshal(resBytes, &responseData)
		if err != nil {
			err = errors.Wrapf(err, "failed parsing response body")
			logs.Error(err)
			c.Abort("500")
			return
		}

		content := fmt.Sprintf("%s is logged in - %s", responseData["email"], AUTH_LANDING_TXT)

		c.Data["Content"] = content
		c.Data["Url"] = "http://localhost:9000/logout"
		c.Data["Label"] = "Logout"
		c.Data["Title"] = "Logged In Page"

		c.TplName = "index.tpl"
		return

	}

	// if we have no token at all, redirect to landing page
	c.Redirect("http://localhost:9000/", 302)
}
