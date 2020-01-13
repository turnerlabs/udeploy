package user

import (
	"context"
)

// Inherit ...
func Inherit(ctx context.Context, u User) (User, error) {

	for _, role := range u.Roles {

		usr, err := Get(ctx, role)
		if err != nil {
			return u, err
		}

		usr, err = Inherit(ctx, usr)
		if err != nil {
			return u, err
		}

		u = override(u, usr)
	}

	return u, nil
}

func override(u, i User) User {

	for name, a := range i.Apps {
		_, found := u.Apps[name]
		if !found {
			u.Apps[name] = merge(u.Apps[name], a)
		}

		u.Apps[name] = a
	}

	return u
}

func merge(userApp AppClaim, app AppClaim) AppClaim {

	for appName, claims := range app.Claims {
		for _, c := range claims {
			userApp.SetPermission(appName, c)
		}
	}

	return userApp
}
