package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jcuvillier/csvhandler"
)

func main() {
	f, err := os.Open("./test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader, err := csvhandler.NewReader(csv.NewReader(f))
	if err != nil {
		log.Fatal(err)
	}
	for {
		// Read handler to get a record
		record, err := reader.Read()
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
}
