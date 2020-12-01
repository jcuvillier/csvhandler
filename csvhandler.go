package csvhandler

import (
	"encoding/csv"
	"io"
	"sync"
)

type CSVHandler struct {
	reader  *csv.Reader
	headers map[string]int
	mutex   *sync.Mutex
}

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
			return CSVHandler{}, &ErrDuplicateKey{key: v}
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
