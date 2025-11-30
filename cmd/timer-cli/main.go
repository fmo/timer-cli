package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	taskFile = "tasks.csv"
	layout   = "02-01-2006 15:04:05"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("need at least one argument"))
	}

	var start, end time.Time
	var startString, endString string

	f, err := getCsv()
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, err := read(f)
	if err != nil {
		log.Fatal(err)
	}
    if data != nil {
    	start, err = time.Parse(layout, data[1][0])
	    end, err = time.Parse(layout, data[1][1])

	    startString = start.Format(layout)
	    endString = end.Format(layout)
    }

	switch os.Args[1] {
	case "start":
		start = time.Now()
		startString = start.Format(layout)
		endString = ""
	case "end":
		end = time.Now()
		endString = end.Format(layout)
	default:
		duration := end.Sub(start)
		fmt.Println(duration)
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

	if err := write(f, startString, endString); err != nil {
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
