package line

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var channelID, channelSecret string
var baseUrl = "https://api.line.me/v1"
var oauthUrl = fmt.Sprint(baseUrl, "/oauth")
var LINE_EP = oauth2.Endpoint{
	AuthURL:  "https://access.line.me/dialog/oauth/weblogin",
	TokenURL: fmt.Sprint(oauthUrl, "/accessToken"),
}

type Auth struct {
	config *oauth2.Config
	state  string
}

func (a *Auth) GetURL() string {
	return a.config.AuthCodeURL(a.state, nil)
}

type API struct {
	*oauth2.Config
	*oauth2.Token
}

func (a *Auth) NewAPI(ctx context.Context, code string) *API {
	t, err := a.config.Exchange(ctx, code)
	if err != nil {
		panic(err)
	}
	return &API{Config: a.config, Token: t}
}

func (a *API) Logout() error {
	req, err := http.NewRequest("DELETE", fmt.Sprint(oauthUrl, "/logout"), nil)
	if err != nil {
		return err
	}
	req.Header["X-Line-ChannelToken"] = []string{a.Token.AccessToken}
	c := a.Config.Client(oauth2.NoContext, a.Token)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b := new(bytes.Buffer)
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	if b.String() != `{"result":"OK"}` {
		return fmt.Errorf("unexpected json %s", b.String())
	}
	return nil
}

// Set channel ID
func SetID(id string) {
	channelID = id
}

// Set channel secret.
func SetSecret(secret string) {
	channelSecret = secret
}

// GetAuth will return *Auth that is needed to get authorize URL and also token.
// callback is the callback URL use by LINE server to webserver. state is used
// as CSRF protection. c can be nil or a custom preset oauth2.Config.
func GetAuth(callback, state string, c *oauth2.Config) *Auth {
	var conf *oauth2.Config
	if c == nil {
		conf = &oauth2.Config{
			ClientID:     channelID,
			ClientSecret: channelSecret,
			Endpoint:     LINE_EP,
			RedirectURL:  callback,
		}
	} else {
		conf = c
		if conf.ClientID == "" {
			conf.ClientID = channelID
		}
		if conf.ClientSecret == "" {
			conf.ClientSecret = channelSecret
		}
		if conf.Endpoint.AuthURL == "" {
			conf.Endpoint.AuthURL = LINE_EP.AuthURL
		}
		if conf.Endpoint.TokenURL == "" {
			conf.Endpoint.TokenURL = LINE_EP.TokenURL
		}
		conf.RedirectURL = callback
	}
	if state == "" {
		state = strconv.FormatInt(time.Now().Unix(), 10)
	}
	return &Auth{config: conf, state: state}
}
