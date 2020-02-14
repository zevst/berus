package berus

import "fmt"

func newError(text string) error {
	return &errorString{fmt.Sprintf("berus: %s", text)}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
