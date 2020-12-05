package csvhandler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

type errWriter struct{}

func (w errWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("write error")
}

var r = Record{
	headers: map[string]int{
		"first_name":      0,
		"last_name":       1,
		"age":             2,
		"is_active":       3,
		"registered":      4,
		"balance":         5,
		"mean_connection": 6,
		"out_of_bounds":   -1,
	},
	values: []string{"John", "Smith", "25", "true", "2018-11-05 12:55:10", "15.65", "12m10s"},
}

func TestFprintln(t *testing.T) {
	testcases := map[string]struct {
		columns   []string
		err       bool
		expected  string
		errWriter bool
	}{
		"regular": {
			columns:  []string{"first_name", "last_name"},
			err:      false,
			expected: "first_name='John' last_name='Smith'\n",
		},
		"unknown column": {
			columns: []string{"unknown"},
			err:     true,
		},
		"null writer": {
			columns:   []string{"first_name", "last_name"},
			err:       true,
			errWriter: true,
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			var buf bytes.Buffer
			var w io.Writer
			if tc.errWriter {
				w = errWriter{}
			} else {
				w = &buf
			}
			err := r.Fprintln(w, tc.columns...)
			if tc.err {
				require.Error(t, err)
			} else {
				assert.Equal(t, tc.expected, buf.String())
			}
		})
	}
}

func TestGet(t *testing.T) {
	testcases := map[string]struct {
		key      string
		expected string
		err      bool
		errType  interface{}
	}{
		"regular": {
			key:      "first_name",
			expected: "John",
		},
		"unknown key": {
			key:     "unknown",
			err:     true,
			errType: &ErrUnknownKey{},
		},
		"out of bounds": {
			key:     "out_of_bounds",
			err:     true,
			errType: &ErrOutOfBounds{},
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			val, err := r.Get(tc.key)
			if tc.err {
				require.Error(t, err)
				if tc.errType != nil {
					assert.True(t, errors.As(err, tc.errType))
				}
			} else {
				assert.Equal(t, tc.expected, val)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	testcases := map[string]struct {
		key      string
		expected bool
		err      bool
		errType  interface{}
	}{
		"regular": {
			key:      "is_active",
			expected: true,
		},
		"unknown key": {
			key:     "unknown",
			err:     true,
			errType: &ErrUnknownKey{},
		},
		"not boolean": {
			key:     "first_name",
			err:     true,
			errType: &ErrWrongType{},
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			val, err := r.GetBool(tc.key)
			if tc.err {
				require.Error(t, err)
				if tc.errType != nil {
					assert.True(t, errors.As(err, tc.errType))
				}
			} else {
				assert.Equal(t, tc.expected, val)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	testcases := map[string]struct {
		key      string
		expected int
		err      bool
		errType  interface{}
	}{
		"regular": {
			key:      "age",
			expected: 25,
		},
		"unknown key": {
			key:     "unknown",
			err:     true,
			errType: &ErrUnknownKey{},
		},
		"not int": {
			key:     "first_name",
			err:     true,
			errType: &ErrWrongType{},
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			val, err := r.GetInt(tc.key)
			if tc.err {
				require.Error(t, err)
				if tc.errType != nil {
					assert.True(t, errors.As(err, tc.errType))
				}
			} else {
				assert.Equal(t, tc.expected, val)
			}
		})
	}
}

func TestGetInt64(t *testing.T) {
	testcases := map[string]struct {
		key      string
		expected int64
		err      bool
		errType  interface{}
	}{
		"regular": {
			key:      "age",
			expected: 25,
		},
		"unknown key": {
			key:     "unknown",
			err:     true,
			errType: &ErrUnknownKey{},
		},
		"not int": {
			key:     "first_name",
			err:     true,
			errType: &ErrWrongType{},
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			val, err := r.GetInt64(tc.key)
			if tc.err {
				require.Error(t, err)
				if tc.errType != nil {
					assert.True(t, errors.As(err, tc.errType))
				}
			} else {
				assert.Equal(t, tc.expected, val)
			}
		})
	}
}

func TestGetFloat64(t *testing.T) {
	testcases := map[string]struct {
		key      string
		expected float64
		err      bool
		errType  interface{}
	}{
		"regular": {
			key:      "balance",
			expected: 15.65,
		},
		"unknown key": {
			key:     "unknown",
			err:     true,
			errType: &ErrUnknownKey{},
		},
		"not int": {
			key:     "first_name",
			err:     true,
			errType: &ErrWrongType{},
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			val, err := r.GetFloat64(tc.key)
			if tc.err {
				require.Error(t, err)
				if tc.errType != nil {
					assert.True(t, errors.As(err, tc.errType))
				}
			} else {
				assert.Equal(t, tc.expected, val)
			}
		})
	}
}

func TestGetTime(t *testing.T) {
	testcases := map[string]struct {
		key      string
		expected time.Time
		err      bool
		errType  interface{}
	}{
		"regular": {
			key:      "registered",
			expected: time.Date(2018, time.November, 5, 12, 55, 10, 0, time.UTC),
		},
		"unknown key": {
			key:     "unknown",
			err:     true,
			errType: &ErrUnknownKey{},
		},
		"not int": {
			key:     "first_name",
			err:     true,
			errType: &ErrWrongType{},
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			val, err := r.GetTime("2006-01-02 15:04:05", tc.key)
			if tc.err {
				require.Error(t, err)
				if tc.errType != nil {
					assert.True(t, errors.As(err, tc.errType))
				}
			} else {
				assert.Equal(t, tc.expected, val)
			}
		})
	}
}

func TestGetDuration(t *testing.T) {
	testcases := map[string]struct {
		key      string
		expected time.Duration
		err      bool
		errType  interface{}
	}{
		"regular": {
			key:      "mean_connection",
			expected: 12*time.Minute + 10*time.Second,
		},
		"unknown key": {
			key:     "unknown",
			err:     true,
			errType: &ErrUnknownKey{},
		},
		"not int": {
			key:     "first_name",
			err:     true,
			errType: &ErrWrongType{},
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			val, err := r.GetDuration(tc.key)
			if tc.err {
				require.Error(t, err)
				if tc.errType != nil {
					assert.True(t, errors.As(err, tc.errType))
				}
			} else {
				assert.Equal(t, tc.expected, val)
			}
		})
	}
}
