package csvhandler

import (
	"encoding/csv"
	"io"
	"sync"
)

// CSVHandler reads records from a CSV-encoded file.
//
// It internally wraps a `encoding/csv.Reader` and uses it to read the data.
// It also holds a map keeping the column names with their indexes.
// This handler is thread safe.
type CSVHandler struct {
	reader  *csv.Reader
	headers map[string]int
	mutex   *sync.Mutex
}

// New creates a new CSVHandler from the given reader.
// It also read the first record as column names and store the indexes.
// As it wraps `encoding/csv.Reader`, any errors returned by the `Read`function can be returned here (including io.EOF is the reader is empty).
// If a duplicate is detected among column names, ErrDuplicateKey is returned.
func New(r io.Reader) (CSVHandler, error) {
	csvReader := csv.NewReader(r)

	// Read headers to save column keys
	record, err := csvReader.Read()
	if err != nil {
		return CSVHandler{}, err
	}

	// Iterate over record in header in order to save the colum keys with the index
	headers := make(map[string]int)
	for i, v := range record {
		// Check if key already exist to stop with error
		_, ok := headers[v]
		if ok {
			return CSVHandler{}, ErrDuplicateKey{key: v}
		}

		headers[v] = i
	}
	return CSVHandler{
		reader:  csvReader,
		headers: headers,
		mutex:   &sync.Mutex{},
	}, nil
}

// Read reads one record (a slice of fields) from handler.
// If the record has an unexpected number of fields, Read returns the record along with the error csv.ErrFieldCount.
// If there is no data left to be read, Read returns nil, io.EOF.
func (h CSVHandler) Read() (*Record, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	r, err := h.reader.Read()
	return &Record{
		values:  r,
		headers: h.headers,
	}, err
}
