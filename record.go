package csvhandler

import (
	"fmt"
	"io"
	"os"
)

// Record holds the fields for a given entry.
// It offers utility functions to access field based on the column name
type Record struct {
	headers map[string]int
	values  []string
}

// Fprintln prints into the given writer each given column with a 'key=value' format.
// For instance, Fprintln(w, "first_name", "last_name") writes "first_name='John' last_name='Smith'"
// Expected errors are the same Get() may return
func (r Record) Fprintln(w io.Writer, columns ...string) error {
	for i, c := range columns {
		v, err := r.Get(c)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "%s='%s'", c, v); err != nil {
			return err
		}
		// Add space separator if not the last element
		if i != len(columns)-1 {
			fmt.Fprintf(w, " ")
		}
	}
	_, err := fmt.Fprint(w, "\n")
	return err
}

// Println prints into the standard output each given column with a 'key=value' format.
// For instance, Fprintln(w, "first_name", "last_name") writes "first_name='John' last_name='Smith'"
// Expected errors are the same Get() may return
func (r Record) Println(columns ...string) error {
	return r.Fprintln(os.Stdout, columns...)
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
