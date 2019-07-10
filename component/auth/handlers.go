package auth

import (
	"errors"

	"github.com/turnerlabs/udeploy/component/cfg"
	"github.com/turnerlabs/udeploy/component/integration/oauth"

	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	echosession "github.com/turnerlabs/udeploy/component/session"
)

// Logout ...
func Logout(c echo.Context) error {
	store := echosession.FromContext(c)
	store.Delete(authSessionName)

	if err := store.Save(); err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, cfg.Get["OAUTH_SIGN_OUT_URL"])
}

// Login ...
func Login(c echo.Context) error {
	state, err := json.Marshal(oauth.UpdateState(c.QueryParam(userURLParam)))
	if err != nil {
		return err
	}

	url := oauth.Config.AuthCodeURL(string(state))

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// Response ...
func Response(c echo.Context) error {

	returnedState := oauth.State{}
	if err := json.Unmarshal([]byte(c.QueryParam("state")), &returnedState); err != nil {
		return err
	}

	if returnedState.Invalid() {
		return errors.New("invalid state")
	}

	token, err := oauth.Config.Exchange(context.Background(), c.QueryParam("code"))
	if err != nil {
		return err
	}

	store := echosession.FromContext(c)
	store.Set(authSessionName, token)

	if err := store.Save(); err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, returnedState.UserRequestedPath)
}
