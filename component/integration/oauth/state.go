package oauth

import "github.com/google/uuid"

var instanceID = uuid.New().String()

// State ...
type State struct {
	ID                string
	UserRequestedPath string
}

// Invalid ...
func (s State) Invalid() bool {
	return s.ID != instanceID
}

// UpdateState ...
func UpdateState(path string) State {
	return State{
		ID:                instanceID,
		UserRequestedPath: path,
	}
}
