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
	reader *csv.Reader
	header []string
	mutex  *sync.Mutex
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

	// Check for duplicates
	set := make(map[string]struct{})
	for _, h := range header {
		if _, duplicate := set[h]; duplicate {
			return nil, ErrDuplicateKey{key: h}
		}
		set[h] = struct{}{}
	}

	return &Reader{
		reader: r,
		header: header,
		mutex:  &sync.Mutex{},
	}, nil
}

// Read reads one record (a slice of fields) from handler.
//
// If the record has an unexpected number of fields, Read returns the record along with the error csv.ErrFieldCount.
// If there is no data left to be read, Read returns nil, io.EOF.
func (r *Reader) Read() (*Record, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.reader.FieldsPerRecord = len(r.header)
	record, err := r.reader.Read()

	fields := make(map[string]string)
	for i, v := range record {
		// At this point, we are sure `record` and `r.header` ahve the same size
		fields[r.header[i]] = v
	}

	return &Record{
		fields: fields,
	}, err
}
