package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

const (
    taskFile = "tasks.csv"
    layout = "02-01-2006 15:04:05"
)

func main() {
	var start, end time.Time

	f, err := getCsv()
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, err := read(f)
	if err != nil {
		log.Fatal(err)
	}
    start, err = time.Parse(layout, data[1][0])
    end, err = time.Parse(layout, data[1][1])

	switch os.Args[1] {
	case "start":
		start = time.Now()
    case "end":
        end = time.Now()
	default:
		fmt.Printf("%v %v", data[1][0], data[1][1])
	}

	if err := f.Truncate(0); err != nil {
		log.Fatal(err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	if err := writeHeader(f); err != nil {
		log.Fatal(err)
	}

	if err := write(f, start.Format(layout), end.Format(layout)); err != nil {
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

func write(f *os.File, start, end string) error {
	w := csv.NewWriter(f)

	record := []string{start, end}

	if err := w.Write(record); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func read(f *os.File) ([][]string, error) {
	r := csv.NewReader(f)

	return r.ReadAll()
}
