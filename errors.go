package csvhandler

import "fmt"

// ErrDuplicateKey is error triggered when a duplicate key is detected among header
type ErrDuplicateKey struct {
	key string
}

func (e *ErrDuplicateKey) Error() string {
	return fmt.Sprintf("key '%s' already exists", e.key)
}
