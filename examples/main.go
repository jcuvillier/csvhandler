package main

import (
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
	handler, err := csvhandler.New(f)
	if err != nil {
		log.Fatal(err)
	}
	for {
		record, err := handler.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s is %s\n", record.Get("first_name"), record.Get("last_name"), record.Get("age"))
	}
}
