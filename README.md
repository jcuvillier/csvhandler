# csvhandler [![PkgGoDev](https://pkg.go.dev/badge/github.com/jcuvillier/csvhandler)](https://pkg.go.dev/github.com/jcuvillier/csvhandler)

CSVHandler is a utility package on top of `encoding/csv` to ease read by allowing direct value access with column name.

```golang
// Reads first record and prints column with key 'username'
handler, _ := csvhandler.New(os.Stdin)
record, _ := handler.Read()
fmt.Println(record.Get("username"))
```