package user

import "go.mongodb.org/mongo-driver/bson/primitive"

// User ...
type User struct {
	ID    primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Email string              `json:"email"`
	Admin bool                `json:"admin"`
	Apps  map[string]AppClaim `json:"apps"`
	Roles []string            `json:"roles"`
}

// ListApps ...
func (u User) ListApps() []string {
	apps := []string{}

	for name := range u.Apps {
		apps = append(apps, name)
	}

	return apps
}

// HasPermission ...
func (u User) HasPermission(app, instance, claim string) bool {
	if a, exists := u.Apps[app]; exists {
		if claims, exists := a.Claims[instance]; exists {
			for _, c := range claims {
				if c == claim {
					return true
				}
			}
		}
	}

	return false
}

// AppClaim ...
type AppClaim struct {
	Claims map[string][]string `json:"claims,omitempty"`
}

// HasPermission ...
func (a *AppClaim) HasPermission(instance, claim string) bool {
	if a.Claims == nil {
		return false
	}

	for _, c := range a.Claims[instance] {
		if c == claim {
			return true
		}
	}

	return false
}

// SetPermission ...
func (a *AppClaim) SetPermission(instance, claim string) {
	if a.Claims == nil {
		a.Claims = map[string][]string{}
	}

	for _, c := range a.Claims[instance] {
		if c == claim {
			return
		}
	}

	a.Claims[instance] = append(a.Claims[instance], claim)
}

//InstanceClaim ...
type InstanceClaim struct {
	Claims []string
	Name   string
}
