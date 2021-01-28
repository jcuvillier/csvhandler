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

// errFormatter is a formatter to trigger errors
var errFormatter = func() Formatter {
	return func(v interface{}) (string, error) {
		return "", fmt.Errorf("formatter error")
	}
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
		values    map[string]field
		defaults  map[string]field
		expected  string
		errWriter bool
		comma     rune
	}{
		"regular": {
			values: map[string]field{
				"first_name": {value: "John"},
				"last_name":  {value: "Smith"},
				"age":        {value: 20},
			},
			expected: "John;Smith;20\n",
			comma:    ';',
		},
		"with default and empty": {
			values: map[string]field{
				"first_name": {value: "John"},
			},
			defaults: map[string]field{
				"age": {value: 30},
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
		"formatter error": {
			values: map[string]field{
				"first_name": {value: "John", formatter: errFormatter()},
			},
			errWriter: true,
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
			for k, f := range tc.defaults {
				w.SetDefault(k, f.value, f.formatter)
			}
			// Create record and sets record values
			r := NewRecord()
			for k, f := range tc.values {
				r.Set(k, f.value, f.formatter)
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

func TestGetFormattedValue(t *testing.T) {
	tstFormatter := func(v interface{}) (string, error) {
		return fmt.Sprintf("this is a test, %v", v), nil
	}

	testcases := map[string]struct {
		value      interface{}
		formatter  Formatter
		isDefault  bool
		wFormatter Formatter
		expected   string
		isErr      bool
	}{
		"record field": {
			value:     "record field",
			formatter: defaultFormatter,
			expected:  "record field",
		},
		"formatter error": {
			value:     "record field",
			formatter: errFormatter(),
			isErr:     true,
		},
		"default field": {
			value:     "default field",
			formatter: defaultFormatter,
			isDefault: true,
			expected:  "default field",
		},
		"writer formatter": {
			value:      "toto",
			formatter:  StringFormatter("%v !!"),
			wFormatter: tstFormatter,
			expected:   "this is a test, toto !!",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			w, err := NewWriter(csv.NewWriter(nil))
			require.NoError(t, err)
			r := NewRecord()
			w.SetFormatter("foo", tc.wFormatter)

			if !tc.isDefault {
				r.Set("foo", tc.value, tc.formatter)
			} else {
				w.SetDefault("foo", tc.value, tc.formatter)
			}

			v, err := w.getFormattedValue(r, "foo")
			if tc.isErr {
				require.Error(t, err)
			} else {
				assert.Equal(t, tc.expected, v)
			}
		})
	}
}

func TestSetDefault(t *testing.T) {
	testcases := map[string]struct {
		value      interface{}
		formatters []Formatter
		expected   string
	}{
		"no formatter": {
			value: "value",
		},
		"single formatter": {
			value:      "value",
			formatters: []Formatter{defaultFormatter},
			expected:   "value",
		},
		"chain formatters": {
			value: "value",
			formatters: []Formatter{
				StringFormatter("prefix %v"),
				StringFormatter("%v suffix"),
			},
			expected: "prefix value suffix",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			w := &Writer{
				defaults: make(map[string]field),
			}
			w.SetDefault("key", tc.value, tc.formatters...)

			f, ok := w.defaults["key"]
			require.True(t, ok)
			assert.Equal(t, "value", f.value)
			if f.formatter != nil {
				s, err := f.formatter(f.value)
				require.NoError(t, err)
				assert.Equal(t, tc.expected, s)
			}
		})
	}
}

func TestSetFormatter(t *testing.T) {
	testcases := map[string]struct {
		value      interface{}
		formatters []Formatter
		expected   string
		isNil      bool
	}{
		"no formatter": {
			value: "value",
			isNil: true,
		},
		"single formatter": {
			value:      "value",
			formatters: []Formatter{defaultFormatter},
			expected:   "value",
		},
		"chain formatters": {
			value: "value",
			formatters: []Formatter{
				StringFormatter("prefix %v"),
				StringFormatter("%v suffix"),
			},
			expected: "prefix value suffix",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			w := &Writer{
				formatters: make(map[string]Formatter),
			}
			w.SetFormatter("key", tc.formatters...)

			f, ok := w.formatters["key"]
			if tc.isNil {
				require.False(t, ok)
			} else {
				require.True(t, ok)
				v, err := f(tc.value)
				require.NoError(t, err)
				assert.Equal(t, tc.expected, v)
			}
		})
	}
}
