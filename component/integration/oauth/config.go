package oauth

import (
	"encoding/gob"
	"strings"

	"github.com/turnerlabs/udeploy/component/cfg"
	"golang.org/x/oauth2"
)

// Config ...
var Config *oauth2.Config

func init() {

	gob.Register(&oauth2.Token{})
	gob.Register(&State{})

	endpoint := oauth2.Endpoint{
		AuthURL:  cfg.Get["OAUTH_AUTH_URL"],
		TokenURL: cfg.Get["OAUTH_TOKEN_URL"],
	}

	scopes := []string{}

	s, found := cfg.Get["OAUTH_SCOPES"]
	if found {
		scopes = strings.Split(s, ",")
	}

	Config = &oauth2.Config{
		ClientID:     cfg.Get["OAUTH_CLIENT_ID"],
		ClientSecret: cfg.Get["OAUTH_CLIENT_SECRET"],
		RedirectURL:  cfg.Get["OAUTH_REDIRECT_URL"],
		Scopes:       scopes,
		Endpoint:     endpoint,
	}
}
