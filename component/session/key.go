package session

// ContextKey ...
type ContextKey string

func (u ContextKey) String() string {
	return "context key " + string(u)
}
