package oauth2

import (
	"context"
	. "example.com/go-boot/platform/config"
	"example.com/go-boot/platform/initializer"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"io"
	"net/http"
	"time"
)

var (
	oauthConfig *oauth2.Config
)
var ctx = context.Background()

func init() {
	Routes(initializer.Router.Group("/login"))
}
func init() {
	oauthConfig = &oauth2.Config{
		RedirectURL:  AppConfig.Oauth2.RedirectUrl,
		ClientID:     AppConfig.Oauth2.ClientId,
		ClientSecret: AppConfig.Oauth2.ClientSecret,
		Scopes:       AppConfig.Oauth2.Scopes,
		// Todo not only for Azure Endpoint
		Endpoint: endpoints.AzureAD(AppConfig.Oauth2.Tenant),
	}
}

var (
	// TODO: randomize it
	oauthStateString = "pseudo-random"
)

type UserInfo struct {
	accessTokenSub        string
	idTokenName           string
	accessTokenExpiration time.Time
	idTokenExpiration     time.Time
	idTokenValue          string
	accessTokenValue      string
}

type WebSSO struct {
	State        string `form:"state"`
	Code         string `form:"code"`
	SessionState string `form:"session_state"`
}

var token *oauth2.Token

var webSSO WebSSO

func Routes(rg *gin.RouterGroup) {
	rg.GET("/", showIndex)
	rg.GET("/login", login)
	rg.GET("/logout", logout)
	rg.GET("/oauth2/code/dbwebsso", loginProcess)
	rg.GET("/info", showTokenInfo)
	rg.GET("/external", getExternalSite)
}

func logout(c *gin.Context) {
	c.SetCookie("JSESSIONID", "", 0, "/", "localhost", false, false)
	c.HTML(http.StatusOK, "logout.html", gin.H{})
}

func login(c *gin.Context) {
	url := oauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func showIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
func loginProcess(c *gin.Context) {
	if c.Bind(&webSSO) == nil {

		if webSSO.State != oauthStateString {
			fmt.Errorf("invalid oauth State")
		}
		var err error
		token, err = oauthConfig.Exchange(ctx, webSSO.Code)
		if err != nil {
			fmt.Errorf("code exchange failed: %s", err.Error())
		}

		response := fmt.Sprintf("<html><body>Login Success and Retriving token is successful<br/></body></html>")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(response))
	}
}
func showTokenInfo(c *gin.Context) {

	response := fmt.Sprintf("<html><body>accesstoken : %s<br/>"+
		"refreshtoken : %s<br/>"+
		"tokentype: %s<br/>"+
		"tokenexpiry : %s<br/></body></html>", token.AccessToken, token.RefreshToken, token.TokenType, token.Expiry)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(response))
}
func getExternalSite(c *gin.Context) {
	if webSSO.State != oauthStateString {
		fmt.Errorf("invalid oauth State")
	}

	client := oauthConfig.Client(ctx, token)
	response, err := client.Get("https://gateway.hub.db.de/bizhub-api-secured-with-jwt")

	if err != nil {
		fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Errorf("failed reading response body: %s", err.Error())
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", contents)

}