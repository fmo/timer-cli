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
    var complete = false

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
        loc := time.FixedZone("UTC+3", 3*3600)
    	start, err = time.ParseInLocation(layout, data[1][0], loc)
	    end, err = time.ParseInLocation(layout, data[1][1], loc)
        if data[1][2] == "true" {
            complete = true
        } else {
            complete = false
	    }
        startString = start.Format(layout)
	    endString = end.Format(layout)
    }

	switch os.Args[1] {
	case "start":
		start = time.Now()
		startString = start.Format(layout)
		endString = ""
	case "end":
        complete = true
		end = time.Now()
		endString = end.Format(layout)
	default:
        countTime(start, end, complete)
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

	if err := write(f, startString, endString, complete); err != nil {
		log.Fatal(err)
	}
}

func countTime(start, end time.Time, complete bool) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    var d time.Duration
    for range ticker.C {
        if complete {
            d = end.Sub(start)
        } else {
            now := time.Now()
            d = now.Sub(start)
        }
        d = d.Truncate(time.Second)
        fmt.Print("\033[H\033[2J")
        fmt.Println(d.String())
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

	record := []string{"start", "end", "complete"}

	if err := w.Write(record); err != nil {
		return err
	}

	w.Flush()
	return nil
}

func write(f *os.File, start, end string, complete bool) error {
	w := csv.NewWriter(f)

    completeString := ""

    if complete {
        completeString = "true"
    } else {
        completeString = "false"
    }

	record := []string{start, end, completeString}

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
