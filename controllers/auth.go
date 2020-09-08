package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/session"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

var globalSessions *session.Manager

type AuthController struct {
	beego.Controller
}

func (c *AuthController) Get() {
	baseurl, err := url.Parse(beego.AppConfig.String("authHost"))
	if err != nil {
		logs.Error("failed to parse %s as a URL:%s", beego.AppConfig.String("authHost"), err)
		c.Abort("500")
	}

	httpClient := http.Client{}
	var auth = fusionauth.NewClient(&httpClient, baseurl, beego.AppConfig.String("apiKey"))

	accessToken, oauthErr, err := auth.ExchangeOAuthCodeForAccessToken(c.GetString("code"), beego.AppConfig.String("clientId"), beego.AppConfig.String("clientSecret"), beego.AppConfig.String("redirectUrl"))
	if err != nil {
		logs.Error("Error exchanging access code for token: %s", err)
		c.Abort("500")
	}
	if oauthErr != nil {
		logs.Error("Oauth Error: %s", oauthErr.Error)
		c.Abort("500")
	}

	c.SetSession("fatoken", accessToken)

	//c.ManualOauth2()

	c.Redirect("/loggedin", 302)
}

func (c *AuthController) ManualOauth2() {
	res, err := http.PostForm(fmt.Sprintf("%s/oauth2/token", beego.AppConfig.String("authHost")),
		url.Values{
			"client_id":     {beego.AppConfig.String("clientId")},
			"client_secret": {beego.AppConfig.String("clientSecret")},
			"code":          {c.GetString("code")},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {beego.AppConfig.String("redirectUrl")},
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

	// save token to session
	c.SetSession("token", responseData["access_token"])
}
