package elves

import "errors"

var (
	ErrInvalidCapacity = errors.New("can not set up a negative capacity")
)
