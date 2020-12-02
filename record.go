package csvhandler

// Record holds the fields for a given entry.
// It offers utility functions to access field based on the column name.
type Record struct {
	headers map[string]int
	values  []string
}

// Get returns as a string the field corresponding to the given key.
// If the key is missing, ErrUnknownKey is returned.
// If the corresponding index is out of bounds the current record, ErrOutOfBounds is returned.
func (r Record) Get(key string) (string, error) {
	i, ok := r.headers[key]
	if !ok {
		return "", ErrUnknownKey{key}
	}
	if i < 0 || i > len(r.values) {
		return "", ErrOutOfBounds{key: key, index: i}
	}
	return r.values[i], nil
}
