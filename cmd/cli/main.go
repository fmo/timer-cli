package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

const taskFile = "tasks.csv"

func main() {
	// var start *time.Time

	f, err := getCsv()
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, err := read(f)
	if err != nil {
		log.Fatal(err)
	}

    fmt.Println(data[0])


	if err := writeHeader(f); err != nil {
		log.Fatal(err)
	}

}

func getCsv() (*os.File, error) {
	var f *os.File
	var err error

	f, err = os.OpenFile(taskFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	log.Printf("file %s exists returning", taskFile)

	return f, nil
}

func writeHeader(f *os.File) error {
    w := csv.NewWriter(f)
    
    record := []string{"start", "end"}

    if err := w.Write(record); err != nil {
        return err
    }

    w.Flush()
    return nil
}

func write(f *os.File) error {
	w := csv.NewWriter(f)

	records := [][]string{
		{"start", "end"},
		{"asd", "asd"},
	}

	if err := w.WriteAll(records); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func read(f *os.File) ([][]string, error) {
	r := csv.NewReader(f)

	return r.ReadAll()
}
