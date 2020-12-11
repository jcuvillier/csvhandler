package csvhandler

import (
	"encoding/csv"
	"sync"
)

// Reader reads records from a CSV-encoded file.
//
// It internally wraps a `encoding/csv.Reader` and uses it to read the data.
// It also holds a map keeping the column names with their indexes.
// This Reader is thread safe.
type Reader struct {
	reader  *csv.Reader
	columns map[string]int
	mutex   *sync.Mutex
}

// NewReader creates a new Reader from the given `encoding/csv.Reader`.
//
// If header is empty NewReader will read the first record and extract column names.
// As it wraps `encoding/csv.Reader`, any errors returned by the `Read` function can be returned here (including io.EOF is the reader is empty).
//
// If a duplicate is detected among column names, ErrDuplicateKey is returned.
func NewReader(r *csv.Reader, header ...string) (*Reader, error) {
	if len(header) == 0 {
		// Read headers to save column keys
		var err error
		header, err = r.Read()
		if err != nil {
			return nil, err
		}
	}

	// Iterate over fields header in order to save the column names with the index
	columns := make(map[string]int)
	for i, v := range header {
		// Check if key already exist to stop with error
		_, ok := columns[v]
		if ok {
			return nil, ErrDuplicateKey{key: v}
		}

		columns[v] = i
	}
	return &Reader{
		reader:  r,
		columns: columns,
		mutex:   &sync.Mutex{},
	}, nil
}

// Read reads one record (a slice of fields) from handler.
//
// If the record has an unexpected number of fields, Read returns the record along with the error csv.ErrFieldCount.
// If there is no data left to be read, Read returns nil, io.EOF.
func (r *Reader) Read() (*Record, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	record, err := r.reader.Read()
	return &Record{
		values:  record,
		columns: r.columns,
	}, err
}
