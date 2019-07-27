package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	echosession "github.com/turnerlabs/udeploy/component/session"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/integration/oauth"
	"golang.org/x/oauth2"
)

const (
	// AuthSessionName defines session for the user.
	AuthSessionName = "oauth"

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
	v, ok := store.Get(AuthSessionName)
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

		store.Set(AuthSessionName, newToken)
		if err := store.Save(); err != nil {
			return "", err
		}

		token = newToken
	}

	parsedToken, _ := jwt.Parse(token.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Currently this function causes an error that is ignored since the public key is not
		// defined. Parsing the token does not required the JWT signature verification. At
		// some point it may be worth verifying the signer.
		//
		// https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-signing-key-rollover
		return jwt.ParseRSAPublicKeyFromPEM([]byte{})
	})

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse JWT claims")
	}

	return claims["upn"].(string), nil
}
