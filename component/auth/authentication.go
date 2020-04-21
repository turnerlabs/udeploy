package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	echosession "github.com/turnerlabs/udeploy/component/session"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/integration/oauth"
	"golang.org/x/oauth2"
)

const (
	// AuthTokenName defines session for the user.
	AuthTokenName = "oauth"

	// IDTokenName defines id_token for the user.
	IDTokenName = "id_token"

	invalidSessionErr = "invalid session"

	loginBaseURL = "/oauth2/login"

	// UserURLParam is the url to redirect the user back to after authentication is complete.
	UserURLParam = "user_url"

	// UserIDParam is stored in the request context for auditing, logging, notifications, etc...
	UserIDParam = "user_id"
)

// UnAuthRedirect ...
func UnAuthRedirect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		user, err := ensureUserToken(c)
		if err != nil {
			loginURL := fmt.Sprintf("%s?%s=%s", loginBaseURL, UserURLParam, c.Request().URL.Path)

			return c.Redirect(http.StatusTemporaryRedirect, loginURL)
		}

		c.Set(UserIDParam, user)

		return next(c)
	}
}

// UnAuthError ...
func UnAuthError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		user, err := ensureUserToken(c)
		if err != nil {
			return err
		}

		c.Set(UserIDParam, user)

		return next(c)
	}
}

func ensureUserToken(c echo.Context) (string, error) {
	store := echosession.FromContext(c)
	v, ok := store.Get(AuthTokenName)
	if !ok {
		return "", errors.New(invalidSessionErr)
	}

	uid, ok := store.Get(UserIDParam)
	if !ok {
		return "", errors.New(invalidSessionErr)
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		log.Println("session token is not a map")
		return "", errors.New(invalidSessionErr)
	}

	b, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return "", errors.New(invalidSessionErr)
	}

	token := &oauth2.Token{}
	if err := json.Unmarshal(b, token); err != nil {
		log.Println(err)
		return "", errors.New(invalidSessionErr)
	}

	if !token.Valid() {
		s := oauth.Config.TokenSource(context.Background(), token)

		newToken, err := s.Token()
		if err != nil {
			return "", err
		}

		store.Set(AuthTokenName, newToken)
		if err := store.Save(); err != nil {
			return "", err
		}

		token = newToken
	}

	return uid.(string), nil
}
