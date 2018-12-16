package mimi

const (
	ErrID = iota
)

type Error struct {
	ID      int
	Message string
}

var (
	ErrAA = &Error{
		ID:      ErrID,
		Message: "jagajaga",
	}
)
