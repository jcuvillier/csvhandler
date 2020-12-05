package csvhandler

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
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

// GetBool returns as a boolean the field corresponding to the given key.
// If the key is missing, ErrUnknownKey is returned.
// If the corresponding index is out of bounds the current record, ErrOutOfBounds is returned.
// If the field is not the expected type, ErrWrongType is returned.
func (r Record) GetBool(key string) (bool, error) {
	v, err := r.Get(key)
	if err != nil {
		return false, err
	}
	if v != "true" && v != "false" {
		return false, ErrWrongType{key: key, err: fmt.Errorf("'%s' is not a boolean", v)}
	}
	return v == "true", nil
}

// GetInt returns as an integer the field corresponding to the given key.
// If the key is missing, ErrUnknownKey is returned.
// If the corresponding index is out of bounds the current record, ErrOutOfBounds is returned.
// If the field is not the expected type, ErrWrongType is returned.
func (r Record) GetInt(key string) (int, error) {
	i, err := r.GetInt64(key)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// GetInt64 returns as an integer64 the field corresponding to the given key.
// If the key is missing, ErrUnknownKey is returned.
// If the corresponding index is out of bounds the current record, ErrOutOfBounds is returned.
// If the field is not the expected type, ErrWrongType is returned.
func (r Record) GetInt64(key string) (int64, error) {
	v, err := r.Get(key)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, ErrWrongType{key: key, err: err}
	}

	return i, nil
}

// GetFloat64 returns as an float the field corresponding to the given key.
// If the key is missing, ErrUnknownKey is returned.
// If the corresponding index is out of bounds the current record, ErrOutOfBounds is returned.
// If the field is not the expected type, ErrWrongType is returned.
func (r Record) GetFloat64(key string) (float64, error) {
	v, err := r.Get(key)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, ErrWrongType{key: key, err: err}
	}
	return f, nil
}

// GetTime returns as a time.Time the field corresponding to the given key.
// If the key is missing, ErrUnknownKey is returned.
// If the corresponding index is out of bounds the current record, ErrOutOfBounds is returned.
// If the field cannot  be parsed as a time using the given layout, ErrWrongType is returned.
func (r Record) GetTime(layout, key string) (time.Time, error) {
	v, err := r.Get(key)
	if err != nil {
		return time.Unix(0, 0), err
	}

	t, err := time.Parse(layout, v)
	if err != nil {
		return time.Unix(0, 0), ErrWrongType{key: key, err: err}
	}
	return t, nil
}

// GetDuration returns as a time.Duration the field corresponding to the given key.
// If the key is missing, ErrUnknownKey is returned.
// If the corresponding index is out of bounds the current record, ErrOutOfBounds is returned.
// If the field cannot be parsed as a duration, ErrWrongType is returned.
func (r Record) GetDuration(key string) (time.Duration, error) {
	v, err := r.Get(key)
	if err != nil {
		return 0, err
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, ErrWrongType{key: key, err: err}
	}
	return d, nil
}
