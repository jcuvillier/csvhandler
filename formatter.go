package csvhandler

import (
	"fmt"
	"time"
)

// Formatter is the function that returns a string formatted version of the given value or an error.
type Formatter func(interface{}) (string, error)

// defaultFormatter is the formatter used when no formatter is specified by caller.
// It printfs the value with a basic `fmt.Sprintf("%v")`
func defaultFormatter(value interface{}) (string, error) {
	return fmt.Sprintf("%v", value), nil
}

// StringFormatter returns a new formatter that uses the given format.
// Format is applied using `fmt.Sprintf`
func StringFormatter(format string) Formatter {
	return func(value interface{}) (string, error) {
		return fmt.Sprintf(format, value), nil
	}
}

// TimeFormatter returns a new formatter that uses the given layout to format a time.
// Only time.Time and *time.Time are allowed as value for the returned formatter.
func TimeFormatter(layout string) Formatter {
	return func(value interface{}) (string, error) {
		var t time.Time
		switch v := value.(type) {
		case time.Time:
			t = time.Time(v)
		case *time.Time:
			t = *(*time.Time)(v)
		default:
			return "", fmt.Errorf("%v (%T) is not a time", value, value)
		}

		return t.Format(layout), nil
	}
}

func chainFormatter(formatters ...Formatter) Formatter {
	return func(value interface{}) (string, error) {
		v := value
		for i, f := range formatters {
			if i == len(formatters)-1 {
				return f(v)
			}
			val, err := f(v)
			if err != nil {
				return "", err
			}
			v = val
		}
		return "", nil // Function should return with the condition within the for loop
	}
}
