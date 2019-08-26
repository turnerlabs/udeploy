package app

import "fmt"

const (
	Pending  = "pending"
	Running  = "running"
	Stopped  = "stopped"
	Deployed = "deployed"
	Error    = "error"

	ErrorTypeAction = "action"

	ChangeTypeStatus  = "STATUS"
	ChangeTypeVersion = "VERSION"
	ChangeTypeError   = "ERROR"
)

// NewState ...
func NewState() State {
	return State{
		Is:    Stopped,
		Error: nil,
	}
}

// State ...
type State struct {
	Is string `json:"-" bson:"-"`

	Version string `json:"-" bson:"-"`
	Error   error  `json:"-" bson:"-"`
}

// IsPending ...
func (s *State) IsPending() bool {
	return s.Is == Pending
}

// IsRunning ...
func (s *State) IsRunning() bool {
	return s.Is == Running
}

// IsStopped ...
func (s *State) IsStopped() bool {
	return s.Is == Stopped
}

// SetPending ...
func (s *State) SetPending() {
	s.Is = Pending
}

// SetStopped ...
func (s *State) SetStopped() {
	s.Is = Stopped
}

// SetRunning ...
func (s *State) SetRunning() {
	s.Is = Running
}

// SetError ...
func (s *State) SetError(err error) {
	s.Error = err
}

// Change ...
type Change struct {
	Before string
	After  string
}

func (c Change) String() string {
	return fmt.Sprintf("%s => %s", c.Before, c.After)
}

// ChangedFrom ...
func (s State) ChangedFrom(prev State) (bool, map[string]Change) {
	changes := map[string]Change{}

	if s.Is != prev.Is {
		changes[ChangeTypeStatus] = Change{
			Before: prev.Is,
			After:  s.Is,
		}
	}

	if s.Version != prev.Version {
		changes[ChangeTypeVersion] = Change{
			Before: prev.Version,
			After:  s.Version,
		}
	}

	if s.Error != nil && prev.Error == nil {
		changes[ChangeTypeError] = Change{
			Before: fmt.Sprintf("%s", prev.Error),
			After:  fmt.Sprintf("%s", s.Error),
		}
	}

	return len(changes) > 0, changes
}

// StatusError ...
type StatusError struct {
	Type  string
	Value string
}

func (s StatusError) Error() string {
	return s.Value
}
