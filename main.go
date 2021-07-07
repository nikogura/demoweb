package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
)

const APP_HOST = "localhost"
const APP_PORT = 9999

const AUTH_HOST = "http://localhost:9011"
const AUTH_CLIENT_ID = "323407ed-d014-4562-9f24-aa4375e6bb06"
const AUTH_CLIENT_SECRET = "JVqEJkB-PWKSTpc3uFws7nXxhx0djyjzUXdV47BrM7M"
const AUTH_REDIRECT_URL = "http://localhost:9999/oauth-callback"
const AUTH_API_KEY = ""

//go:embed content
var content embed.FS

var FAConfig *oauth2.Config
var oauthStateString = randstr.Hex(16)
var CodeVerifier, _ = cv.CreateCodeVerifier()
var codeChallenge = CodeVerifier.CodeChallengeS256()

func init() {
	log.SetLevel(log.DebugLevel)

	FAConfig = &oauth2.Config{
		ClientID:     AUTH_CLIENT_ID,
		ClientSecret: AUTH_CLIENT_SECRET,
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("%s/oauth2/authorize", AUTH_HOST),
			TokenURL:  fmt.Sprintf("%s/oauth2/token", AUTH_HOST),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		RedirectURL: AUTH_REDIRECT_URL,
		Scopes:      []string{"openid", "offline_access"},
	}
}

func main() {
	router := gin.Default()

	router.Use(ginsession.New())

	router.GET("/", Main)
	router.GET("/login", Login)
	router.GET("/logout", Logout)
	router.GET("/oauth-callback", AuthCallback)

	//router.Use(Serve("/", content))

	addr := fmt.Sprintf("%s:%d", APP_HOST, APP_PORT)
	fmt.Printf("Server starting on %s.\n", addr)

	err := router.Run(addr)
	if err != nil {
		log.Fatalf("failed running server: %s", err)
	}
}

//func Serve(urlPrefix string, efs embed.FS) gin.HandlerFunc {
//	// the embedded filesystem has a 'content/' at the top level.  We wanna strip this so we can treat the root of the views directory as the web root.
//	fsys, err := fs.Sub(efs, "content")
//	if err != nil {
//		log.Fatalf(err.Error())
//	}
//
//	fileserver := http.FileServer(http.FS(fsys))
//	if urlPrefix != "" {
//		fileserver = http.StripPrefix(urlPrefix, fileserver)
//	}
//
//	return func(c *gin.Context) {
//		fileserver.ServeHTTP(c.Writer, c.Request)
//		c.Abort()
//	}
//}

// Boring static html returned on login
func Main(c *gin.Context) {
	content := `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
		<title>Web Tech Demo</title>
		<script type="application/javascript" src="https://unpkg.com/babel-standalone@6.26.0/babel.js"></script>
		<link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" rel="stylesheet">
	</head>
	<body>
		<div className="container">
				<div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
						<h1>Web Technology Demonstrator</h1>
						<a href="/login" className="btn btn-primary btn-lg btn-login btn-block">Sign In</a>
				</div>
		</div>
	</body>
</html>

`

	fmt.Fprintf(c.Writer, content)

}

// Login gets called when the user clicks the login link/button.  This redirects to FA
func Login(c *gin.Context) {
	url := FAConfig.AuthCodeURL(oauthStateString, oauth2.SetAuthURLParam("code_challenge", codeChallenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"))

	c.Redirect(http.StatusFound, url)
}

// AuthCallback This is where FA sends the user back to after logging in.
func AuthCallback(c *gin.Context) {
	userInfo, err := GetUserInfo(c, c.Query("state"), c.Query("code"))
	if err != nil {
		log.Errorf("Error parsing auth data: %s", err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	userInfoBytes, err := json.MarshalIndent(userInfo, "", "  ")
	if err != nil {
		log.Errorf("Error marshalling user info: %s", err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	content := fmt.Sprintf(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
		<title>Web Tech Demo</title>
		<script type="application/javascript" src="https://unpkg.com/babel-standalone@6.26.0/babel.js"></script>
		<link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" rel="stylesheet">
	</head>
	<body>
		<div className="container">
				<div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
						<h1>Web Technology Demonstrator</h1>
						<h2> Logged In User Info:</h2>
						%s <br/>
						<a href="/logout" className="btn btn-primary btn-lg btn-login btn-block">Log Out</a>
				</div>
		</div>
	</body>
</html>
`, string(userInfoBytes))

	fmt.Fprintf(c.Writer, content)
}

func Logout(c *gin.Context) {
	log.Debugf("Logout firing")
	session := ginsession.FromContext(c)
	session.Delete("fatoken")

	url := fmt.Sprintf("%s/oauth2/logout?client_id=%s&post_logout_redirect_uri=%s", AUTH_HOST, AUTH_CLIENT_ID, fmt.Sprintf("http://%s:%d/", APP_HOST, APP_PORT))
	log.Debugf("Redirecting user to %s", url)

	c.Redirect(http.StatusFound, url)
}

// GetUserInfo checks the state string and hits FA for a token.  Once it has one, it can look up the user info from FA
func GetUserInfo(c *gin.Context, state string, code string) (userInfo map[string]interface{}, err error) {
	if state != oauthStateString {
		return userInfo, fmt.Errorf("Invalid oauth state: %q != %q", state, oauthStateString)
	}

	token, err := FAConfig.Exchange(context.TODO(), code, oauth2.SetAuthURLParam("code_verifier", CodeVerifier.String()))
	if err != nil {
		err = errors.Wrapf(err, "failed to exchange token")
		return userInfo, err
	}

	// save the token in the session
	session := ginsession.FromContext(c)
	session.Set("fatoken", token)

	// lookup user info
	url := fmt.Sprintf("%s/oauth2/userinfo", AUTH_HOST)

	var bearer = "Bearer " + token.AccessToken

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		err = errors.Wrapf(err, "failed getting user info")
		return userInfo, err
	}

	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.Wrapf(err, "failed reading user info request rewponse body")
		return userInfo, err
	}

	err = json.Unmarshal(b, &userInfo)
	if err != nil {
		err = errors.Wrapf(err, "failed to unmarshal FA response.")
		return userInfo, err
	}

	return userInfo, err
}

func GetUserInfoFaApi(c *gin.Context) {
	if AUTH_API_KEY == "" {
		log.Error("API KEY is blank.  Can't use this method without it")
		c.Abort()
		return
	}

	baseurl, err := url.Parse(AUTH_HOST)
	if err != nil {
		log.Error("failed to parse %s as a URL:%s", AUTH_HOST, err)
		c.Abort()
		return
	}

	httpClient := http.Client{}
	var auth = fusionauth.NewClient(&httpClient, baseurl, AUTH_API_KEY)

	accessToken, oauthErr, err := auth.ExchangeOAuthCodeForAccessToken(c.Query("code"), AUTH_CLIENT_ID, AUTH_CLIENT_SECRET, AUTH_REDIRECT_URL)
	if err != nil {
		log.Error("Error exchanging access code for token: %s", err)
		c.Abort()
		return
	}
	if oauthErr != nil {
		log.Error("Oauth Error: %s", oauthErr.Error)
		c.Abort()
		return
	}

	// save the token in the session
	session := ginsession.FromContext(c)
	session.Set("fatoken", accessToken)

	//c.ManualOauth2()

	c.Redirect(http.StatusFound, "/loggedin")
}

func ManualOauth2(c *gin.Context) {
	res, err := http.PostForm(fmt.Sprintf("%s/oauth2/token", AUTH_HOST),
		url.Values{
			"client_id":     {AUTH_CLIENT_ID},
			"client_secret": {AUTH_CLIENT_SECRET},
			"code":          {c.Query("code")},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {AUTH_REDIRECT_URL},
		})

	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.Wrapf(err, "failed reading response body")
		log.Error(err)
		c.Abort()
		return
	}

	responseData := make(map[string]interface{})

	err = json.Unmarshal(resBytes, &responseData)
	if err != nil {
		err = errors.Wrapf(err, "failed parsing response body")
		log.Error(err)
		c.Abort()
		return
	}

	// save token to session
	session := ginsession.FromContext(c)
	session.Set("token", responseData["access_token"])
}
