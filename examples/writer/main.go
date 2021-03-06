package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jcuvillier/csvhandler"
)

const (
	firstName = "first_name"
	lastName  = "last_name"
	age       = "age"
)

func toUpper(value interface{}) (string, error) {
	return strings.ToUpper(fmt.Sprintf("%v", value)), nil
}

func main() {
	// Create a `encoding/csv.Writer` to be used
	csvWriter := csv.NewWriter(os.Stdout)
	// This Writer is used to customize how the csv is formated (field delimiter and line terminator)
	csvWriter.Comma = ';'

	// Create the Writer using the previous `encoding/csv.Writer` and the given header
	writer, err := csvhandler.NewWriter(csvWriter, firstName, lastName, age)
	if err != nil {
		log.Fatal(err)
	}
	// ToUpper formatter for _lastname_ column
	writer.SetFormatter(lastName, toUpper)

	// Write header line
	if err := writer.WriteHeader(); err != nil {
		log.Fatal(err)
	}

	// Create a Record then writes it
	r := csvhandler.NewRecord()
	r.Set(lastName, "Smith")
	r.Set(firstName, "John")
	r.Set(age, 25)
	if err := writer.Write(r); err != nil { // Writes John;SMITH;25
		log.Fatal(err)
	}

	// Default value can be specified for a given key
	writer.SetDefault(age, 20)

	r = csvhandler.NewRecord()
	r.Set(lastName, "Smith")
	r.Set(firstName, "Laura")
	if err := writer.Write(r); err != nil { // Writes Laura;SMITH;20
		log.Fatal(err)
	}
}
