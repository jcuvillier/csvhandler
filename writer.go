package csvhandler

import (
	"encoding/csv"
	"fmt"
	"sync"
)

const defaultEmptyValue = ""

// A Writer writes records using CSV encoding.
//
// It internally uses a `encoding/csv.Writer` to write the records.
type Writer struct {
	writer     *csv.Writer
	header     []string
	defaults   map[string]field
	formatters map[string]Formatter
	mutex      *sync.Mutex
	EmptyValue string
}

// NewWriter creates a new Writer from the given `encoding/csv.Wrtiter` and header.
//
// If a duplicate is detected among column names, ErrDuplicateKey is returned.
func NewWriter(w *csv.Writer, header ...string) (*Writer, error) {
	// Check for duplicates in header
	set := make(map[string]struct{})
	for _, h := range header {
		if _, duplicate := set[h]; duplicate {
			return nil, ErrDuplicateKey{key: h}
		}
		set[h] = struct{}{}
	}

	return &Writer{
		writer:     w,
		header:     header,
		defaults:   make(map[string]field),
		formatters: make(map[string]Formatter),
		mutex:      &sync.Mutex{},
		EmptyValue: defaultEmptyValue,
	}, nil
}

// SetDefault sets the default value to be used if there is no value for this key in the record.
//
// If the defined value is nil, default value is used.
func (w *Writer) SetDefault(key string, value interface{}, formatter ...Formatter) {
	var f Formatter
	if len(formatter) == 1 {
		f = formatter[0]
	} else if len(formatter) > 1 {
		f = chainFormatter(formatter...)
	}
	w.defaults[key] = field{
		value:     value,
		formatter: f,
	}
}

// SetFormatter sets the formatter to be used
func (w *Writer) SetFormatter(key string, formatter ...Formatter) {
	if len(formatter) == 0 {
		return
	} else if len(formatter) == 1 {
		w.formatters[key] = formatter[0]
	} else if len(formatter) > 1 {
		w.formatters[key] = chainFormatter(formatter...)
	}
}

// WriteHeader writes the header line of the CSV.
//
// Field delimiter used is the one specified in the `encoding/csv.Writer` given when creating this Writer.
// Header keys are written in the same order as specified in `NewWriter` function.
func (w *Writer) WriteHeader() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if len(w.header) != 0 {
		if err := w.writer.Write(w.header); err != nil {
			return fmt.Errorf("cannot write header line: %s", err)
		}
	}
	w.writer.Flush()
	if err := w.writer.Error(); err != nil {
		return fmt.Errorf("cannot write header line: %s", err)
	}
	return nil
}

// Write writes the given record as a new line.
//
// Field delimiter used is the one specified in the `encoding/csv.Writer` given when creating this Writer.
// Fields are written in the header order specified in `NewWriter` function.
// If field is not specified in the record, a specified default value (see function SetDefault())
// 	can be used, otherwise EmptyValue is used.
// Fields with key not in header will be ignored.
func (w *Writer) Write(r *Record) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	record := make([]string, 0, len(w.header))

	for _, h := range w.header {
		value, err := w.getFormattedValue(r, h)
		if err != nil {
			return err
		}
		record = append(record, value)
	}

	if err := w.writer.Write(record); err != nil {
		return fmt.Errorf("cannot write record: %s", err)
	}
	w.writer.Flush()
	if err := w.writer.Error(); err != nil {
		return fmt.Errorf("cannot write record: %s", err)
	}
	return nil
}

// getFormattedValue returned the formatted value of the given record and column.
//
// Value used is from:
// 1. corresponding field of the record
// 2. default value defined for the column
// 3. Writer's EmptyValue
//
// Formatter used is from:
// 1. associated formatter to the field or defaultValue depending on the value used
// 2. defaultFormatter if both are missing
// 3. formatter defined for column is chained if specified
func (w *Writer) getFormattedValue(record *Record, column string) (string, error) {
	var f Formatter
	var v interface{}
	// Use EmptyValue if record has no field and no defaultValue is set
	v = w.EmptyValue
	if field, hasField := record.fields[column]; hasField {
		v = field.value
		f = field.formatter
	} else if defValue, hasDefault := w.defaults[column]; hasDefault {
		v = defValue.value
		f = defValue.formatter
	}

	if f == nil {
		// No formatter defined at all, fallback to defaultFormatter
		f = defaultFormatter
	}

	// Finally, check for column formatter, if present chain with field formatter
	if formatter, hasFormatter := w.formatters[column]; hasFormatter && formatter != nil {
		f = chainFormatter(f, formatter)
	}

	return f(v)
}

// WriteAll writes all the given records using the Write function.
func (w *Writer) WriteAll(r []*Record) error {
	for _, record := range r {
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.writer.Flush()
	return w.writer.Error()
}
