package csvhandler

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tstWriter struct {
	w   io.Writer
	err error
}

func (t tstWriter) Write(p []byte) (n int, err error) {
	if t.err != nil {
		return 0, t.err
	}
	return t.w.Write(p)
}
func TestNewWriter(t *testing.T) {
	testcases := map[string]struct {
		header  []string
		err     bool
		errType interface{}
	}{
		"regular": {
			header: []string{"first_name", "last_name", "age"},
		},
		"error duplicate header": {
			header:  []string{"first_name", "first_name", "age"},
			err:     true,
			errType: &ErrDuplicateKey{},
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			w, err := NewWriter(csv.NewWriter(nil), tc.header...)
			if tc.err {
				require.Error(t, err)
				assert.True(t, errors.As(err, tc.errType))
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.header, w.header)
				assert.NotNil(t, w.writer)
			}
		})
	}
}

func TestWriteHeader(t *testing.T) {
	testcases := map[string]struct {
		header    []string
		expected  string
		errWriter bool
		comma     rune
	}{
		"regular": {
			header:   []string{"first_name", "last_name", "age"},
			expected: "first_name,last_name,age\n",
			comma:    ',',
		},
		"empty header": {},
		"write error": {
			header:    []string{"first_name", "last_name", "age"},
			errWriter: true,
		},
		"csv writer error": { // Forcing a csv write error by specifying an invalid comma rune
			header:    []string{"first_name", "last_name", "age"},
			errWriter: true,
			comma:     'a',
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			// Create tst writer
			var b bytes.Buffer
			tw := tstWriter{
				w: &b,
			}
			if tc.errWriter {
				tw.err = fmt.Errorf("write error")
			}
			cw := csv.NewWriter(tw)
			cw.Comma = tc.comma

			// Create writer
			w, err := NewWriter(cw, tc.header...)
			require.NoError(t, err)

			// Write header and check
			err = w.WriteHeader()
			if tc.errWriter {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, b.String())
			}
		})
	}
}

func TestWrite(t *testing.T) {
	testcases := map[string]struct {
		values    map[string]interface{}
		defaults  map[string]interface{}
		expected  string
		errWriter bool
		comma     rune
	}{
		"regular": {
			values: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Smith",
				"age":        20,
			},
			expected: "John;Smith;20\n",
			comma:    ';',
		},
		"with default and empty": {
			values: map[string]interface{}{
				"first_name": "John",
			},
			defaults: map[string]interface{}{
				"age": 30,
			},
			expected: "John;;30\n",
			comma:    ';',
		},
		"write error": {
			errWriter: true,
		},
		"csv writer error": { // Forcing a csv write error by specifying an invalid comma rune
			errWriter: true,
			comma:     'a',
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			// Create tst writer
			var b bytes.Buffer
			tw := tstWriter{
				w: &b,
			}
			if tc.errWriter {
				tw.err = fmt.Errorf("write error")
			}
			cw := csv.NewWriter(tw)
			cw.Comma = tc.comma

			// Create writer
			w, err := NewWriter(cw, "first_name", "last_name", "age")
			require.NoError(t, err)

			// Set default values
			for k, v := range tc.defaults {
				w.SetDefault(k, v)
			}
			// Create record and sets record values
			r := NewRecord()
			for k, v := range tc.values {
				r.Set(k, v)
			}

			// Write record and check
			err = w.Write(r)
			if tc.errWriter {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, b.String())
			}

		})
	}
}

func TestWriteAll(t *testing.T) {
	testcases := map[string]struct {
		values    []map[string]interface{}
		expected  string
		errWriter bool
		comma     rune
	}{
		"regular": {
			values: []map[string]interface{}{
				{
					"first_name": "John",
					"last_name":  "Smith",
					"age":        20,
				},
				{
					"first_name": "Holly",
					"last_name":  "Franklin",
					"age":        27,
				},
			},
			expected: "John,Smith,20\nHolly,Franklin,27\n",
			comma:    ',',
		},
		"write error": {
			values: []map[string]interface{}{
				{
					"first_name": "John",
					"last_name":  "Smith",
					"age":        20,
				},
			},
			errWriter: true,
			comma:     'a',
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			// Create tst writer
			var b bytes.Buffer
			tw := tstWriter{
				w: &b,
			}
			if tc.errWriter {
				tw.err = fmt.Errorf("write error")
			}
			cw := csv.NewWriter(tw)
			cw.Comma = tc.comma

			// Create writer
			w, err := NewWriter(cw, "first_name", "last_name", "age")
			require.NoError(t, err)

			var records []*Record
			for _, values := range tc.values {
				r := NewRecord()
				for k, v := range values {
					r.Set(k, v)
				}
				records = append(records, r)
			}

			// Write records and check
			err = w.WriteAll(records)
			if tc.errWriter {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, b.String())
			}
		})
	}
}
