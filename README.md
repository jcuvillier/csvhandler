# CSVHandler [![PkgGoDev](https://pkg.go.dev/badge/github.com/jcuvillier/csvhandler)](https://pkg.go.dev/github.com/jcuvillier/csvhandler) [![CI](https://github.com/jcuvillier/csvhandler/workflows/Go/badge.svg)](https://github.com/jcuvillier/csvhandler/actions?query=workflow%3AGo)

CSVHandler is a utility package on top of `encoding/csv` to ease read by allowing direct value access with column name. 

```
go get github.com/jcuvillier/csvhandler
```

Working with csv can be tedious, especially when dealing with large number of columns. This package keep correspondence between column name and index and provides utility functions to access data in records.

```golang
    handler, err := csvhandler.New(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	for {
		// Read handler to get a record
		record, err := handler.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Read first_name, last_name and age column from record
		firstName, err := record.Get("first_name")
		if err != nil {
			log.Fatal(err)
		}
		lastName, err := record.Get("last_name")
		if err != nil {
			log.Fatal(err)
		}
		age, err := record.GetInt("age")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s is %d\n", firstName, lastName, age)
	}
```
with
```
id,first_name,last_name,age
cb846443-97b2-4b1a-ad62-682c99b70604,Holly,Franklin,27
618d824c-0052-4298-b049-d35b44e53e03,Giacobo,Tolumello,18
1768891f-c658-44b4-9912-8e0a2f51cf4f,Aubrie,Bellie,32
06dca1a4-8686-4126-9bf9-2071b2881810,Kristoforo,Lifsey,59
61d4e115-d6d6-4c9e-97a2-63086fdfa03e,Jasmine,Rayhill,35
```
prints
```
Holly Franklin is 27
Giacobo Tolumello is 18
Aubrie Bellie is 32
Kristoforo Lifsey is 59
Jasmine Rayhill is 35
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
r.Println("first_name", "last_name")
```
prints to standard output
```
first_name='Holly' last_name='Franklin'
```

