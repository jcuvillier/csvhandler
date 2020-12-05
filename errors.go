package csvhandler

import "fmt"

// ErrDuplicateKey means a duplicate key is detected within header
type ErrDuplicateKey struct {
	key string
}

func (e ErrDuplicateKey) Error() string {
	return fmt.Sprintf("key '%s' already exists", e.key)
}

// ErrUnknownKey means a requested key does not exist within header
type ErrUnknownKey struct {
	key string
}

func (e ErrUnknownKey) Error() string {
	return fmt.Sprintf("key '%s' does not exist", e.key)
}

// ErrOutOfBounds means a requested key corresponds to an index that is out of bounds the current record
type ErrOutOfBounds struct {
	key   string
	index int
}

func (e ErrOutOfBounds) Error() string {
	return fmt.Sprintf("key '%s' with index %d is out of bounds", e.key, e.index)
}

// ErrWrongType means the field with the requested key is not the expected type
type ErrWrongType struct {
	key string
	err error
}

func (e ErrWrongType) Error() string {
	return fmt.Sprintf("field with key '%s' is not the expected type, %v", e.key, e.err)
}
