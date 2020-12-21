package csvhandler

import (
	"encoding/csv"
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

const defaultEmptyValue = ""

// A Writer writes records using CSV encoding.
//
// It internally uses a `encoding/csv.Writer` to write the records.
// Therefore
type Writer struct {
	writer     *csv.Writer
	header     []string
	defaults   map[string]string
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
		defaults:   make(map[string]string),
		mutex:      &sync.Mutex{},
		EmptyValue: defaultEmptyValue,
	}, nil
}

// SetDefault sets the default value to be used if there is no value for this key in the record.
//
// If the defined value is nil, default value is used.
func (w *Writer) SetDefault(key string, value interface{}) {
	w.defaults[key] = fmt.Sprintf("%v", value)
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
			return errors.Wrap(err, "cannot write header line")
		}
	}
	w.writer.Flush()
	return errors.Wrap(w.writer.Error(), "cannot write header line")
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
		if value, hasValue := r.fields[h]; hasValue {
			record = append(record, value)
		} else if defValue, hasDefault := w.defaults[h]; hasDefault {
			record = append(record, fmt.Sprintf("%v", defValue))
		} else {
			record = append(record, fmt.Sprintf("%v", w.EmptyValue))
		}
	}

	if err := w.writer.Write(record); err != nil {
		return errors.Wrap(err, "cannot write record")
	}
	w.writer.Flush()
	return errors.Wrap(w.writer.Error(), "cannot write record")
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
