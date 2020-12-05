package csvhandler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tstReader struct {
	data []byte
	err  error
}

func (r tstReader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	buf := bytes.NewBuffer(r.data)
	return buf.Read(p)
}

func TestNew(t *testing.T) {
	testcases := map[string]struct {
		data       []byte
		headersLen int
		errReader  bool
		errType    interface{}
	}{
		"regular": {
			data: []byte(`first_name,last_name
			Holly,Franklin`),
			headersLen: 2,
		},
		"read error": {
			errReader: true,
		},
		"duplicate key": {
			data: []byte(`first_name,first_name
			Holly,Franklin`),
			errType: &ErrDuplicateKey{},
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			// Prepare reader
			r := tstReader{
				data: tc.data,
			}
			if tc.errReader {
				r.err = fmt.Errorf("read error")
			}

			handler, err := New(r)
			if tc.errReader || tc.errType != nil {
				require.Error(t, err)
				if tc.errType != nil {
					assert.True(t, errors.As(err, tc.errType))
				}
			} else {
				assert.Len(t, handler.headers, tc.headersLen)
			}
		})
	}
}

func TestReadProcessCSV(t *testing.T) {
	testcases := map[string]struct {
		filename string
		names    []string
	}{
		"regular": {
			filename: "regular.csv",
			names:    []string{"Holly", "Giacobo", "Aubrie", "Kristoforo", "Jasmine"},
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			f, err := os.Open(filepath.Join("tstdata", tc.filename))
			require.NoError(t, err)
			defer f.Close()
			handler, err := New(f)
			require.NoError(t, err)

			/// Read all records from file
			var records []Record
			for {
				// Read handler to get a record
				record, err := handler.Read()
				if err == io.EOF {
					break
				}
				require.NoError(t, err)
				records = append(records, *record)
			}
			// Check the number of records is OK
			assert.Len(t, records, len(tc.names))

			// Check the names are corrects
			for i := range records {
				n, err := records[i].Get("first_name")
				require.NoError(t, err)
				assert.Equal(t, tc.names[i], n)
			}
		})
	}

}
