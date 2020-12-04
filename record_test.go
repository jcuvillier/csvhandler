package csvhandler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

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
	values: []string{"John", "Smith", "25", "true", "2018-11-05 12:55:10", "15.65", "12m"},
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
