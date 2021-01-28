package csvhandler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultFormatter(t *testing.T) {
	testcases := map[string]struct {
		value    interface{}
		expected string
	}{
		"string": {
			value:    "foo",
			expected: "foo",
		},
		"bool": {
			value:    false,
			expected: "false",
		},
		"int": {
			value:    10,
			expected: "10",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			res, err := defaultFormatter(tc.value)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestStringFormatter(t *testing.T) {
	testcases := map[string]struct {
		value    interface{}
		format   string
		expected string
	}{
		"regular": {
			value:    "World",
			format:   "Hello %s !",
			expected: "Hello World !",
		},
	}
	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			res, err := StringFormatter(tc.format)(tc.value)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestTimeFormatter(t *testing.T) {
	value := time.Date(1998, time.July, 12, 22, 30, 0, 0, time.Local)
	testcases := map[string]struct {
		value    interface{}
		layout   string
		expected string
		err      bool
	}{
		"time": {
			value:    value,
			layout:   time.ANSIC,
			expected: "Sun Jul 12 22:30:00 1998",
		},
		"time pointer": {
			value:    &value,
			layout:   time.ANSIC,
			expected: "Sun Jul 12 22:30:00 1998",
		},
		"not time": {
			value: false,
			err:   true,
		},
	}
	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			res, err := TimeFormatter(tc.layout)(tc.value)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, res)
			}
		})
	}
}

func TestChainFormatter(t *testing.T) {
	testcases := map[string]struct {
		formatters []Formatter
		value      interface{}
		expected   string
		err        bool
	}{
		"regular": {
			formatters: []Formatter{
				TimeFormatter("_2 Jan"),
				StringFormatter("Santa is coming on %s"),
			},
			value:    time.Date(2021, time.December, 25, 8, 0, 0, 0, time.Local),
			expected: "Santa is coming on 25 Dec",
		},
		"with error": {
			formatters: []Formatter{
				TimeFormatter("_2 Jan"),
				StringFormatter("Santa is coming on %s"),
			},
			value: nil,
			err:   true,
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			res, err := chainFormatter(tc.formatters...)(tc.value)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, res)
			}
		})
	}
}
