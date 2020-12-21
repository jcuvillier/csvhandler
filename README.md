# CSVHandler

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jcuvillier/csvhandler)](https://pkg.go.dev/github.com/jcuvillier/csvhandler) [![CI](https://github.com/jcuvillier/csvhandler/workflows/Go/badge.svg)](https://github.com/jcuvillier/csvhandler/actions?query=workflow%3AGo) [![codecov](https://codecov.io/gh/jcuvillier/csvhandler/branch/master/graph/badge.svg?token=EUSKNU9LOP)](https://codecov.io/gh/jcuvillier/csvhandler)

Golang built-in `encoding/csv` package can be tedious to use as it handles records as low level types `[]string`, forcing users to access field using indexes.  
This package aims to ease this use by allowing direct access with the column name, with an API close to `encoding/csv`.  

For a usage closer to `encoding/json` with `marshal/unmarshal` functions to/from a struct, consider using other package such as [github.com/jszwec/csvutil](https://github.com/jszwec/csvutil) or [github.com/gocarina/gocsv](https://github.com/gocarina/gocsv).

## Installation

```
go get github.com/jcuvillier/csvhandler
```

## Reader

```golang
csvInput := bytes.NewBuffer([]byte(`
first_name,last_name,age
Holly,Franklin,27
Giacobo,Tolumello,18`,
))

reader, _ := csvhandler.NewReader(csv.NewReader(csvInput))
record, _ := reader.Read()
record.Get("first_name") // returns Holly
record.GetInt("age")     // return 27
```

## Writer

```golang
// Create a writer to stdout with header "first_name,last_name,age"
writer, _ := csvhandler.NewWriter( csv.NewWriter(os.Stdout), "first_name", "last_name","age")

// Write header line
writer.WriteHeader() // Writes first_name,last_name,age

// Create a record to be written
record := NewRecord()
record.Set("first_name", "Holly")
record.Set("last_name", "Franklin")
record.Set("age", 27)
writer.Write(record) // Writes Holly,Franklin,27
```

If a field is not specified, `Writer.EmptyValue` is used. A default value can also be provided with `Writer.SetDefault` function.

```golang
writer, _ := csvhandler.NewWriter( csv.NewWriter(os.Stdout), "first_name", "last_name","age")
writer.SetDefault("age", 18)

record := NewRecord()
record.Set("first_name", "Holly")
writer.Write(record) // Writes Holly,,18
```

## Record

### Get fields

Record holds the fields for a given entry.  
It offers utility functions to `Get` field based on the column name to a given type, for instance:
```golang
GetInt(key string) (int, error)
```
Supported types are `string`, `bool`, `int`, `int64`, `float64`, `time.Time` and `timne.Duration`.

### Print fields

You can also print as key/value pairs a record by giving the column name.

```golang
r.Println("first_name", "last_name") // prints to stdout "first_name='Holly' last_name='Franklin'"
```
