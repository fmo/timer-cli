package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

var ErrNoDataInCSV = errors.New("no data in csv")

type TaskStorer interface {
	Save(Task) error
	Update(*Task) error
}

type CSVCodec struct {
	File   *os.File
	Logger *log.Logger
	Writer *csv.Writer
	Reader *csv.Reader
}

func NewCSVCodec(f *os.File, logger *log.Logger) (*CSVCodec, error) {
	writer := csv.NewWriter(f)
	reader := csv.NewReader(f)

	codec := &CSVCodec{
		File:   f,
		Logger: logger,
		Writer: writer,
		Reader: reader,
	}

	// only for adding header to the csv if its an empty csv
	if data, err := codec.Load(); data == nil {
		if err != nil {
			return nil, fmt.Errorf("cant load csv: %w", err)
		}
		writer.Write([]string{"start", "end", "status"})
		writer.Flush()
	}

	return codec, nil
}

func (c *CSVCodec) Save(task Task) error {
	start := task.Start.Format(time.RFC3339)
	var end string
	if task.End.IsZero() {
		end = ""
	} else {
		end = task.End.Format(time.RFC3339)
	}

	if err := c.Writer.Write([]string{start, end, string(task.Status)}); err != nil {
		return err
	}

	c.Writer.Flush()
	return nil
}

func (c *CSVCodec) Update(task *Task) error {
	data, err := c.Load()
	if err != nil {
		c.Logger.Printf("Cant load data: %v", err)
		return err
	}

	for rowKey, row := range data {
		c.Logger.Printf("row %d: %s, task: %v", rowKey, row, task)
		start, err := time.Parse(time.RFC3339, row[0])
		if err != nil {
			return fmt.Errorf("cant convert start time to time: %w", err)
		}
		c.Logger.Printf("start in time obj: %v", start)
		if task.Start.Equal(start) {
			c.Logger.Printf("found task to update")
		}
	}

	return nil
}

func (c *CSVCodec) Load() ([][]string, error) {
	c.File.Seek(0, 0)
	data, err := c.Reader.ReadAll()
	if err != nil {
		return nil, err
	}
	c.Logger.Println("data:", data)
	if len(data) == 0 {
		return nil, ErrNoDataInCSV
	}
	if data[0][0] != "start" || data[0][1] != "end" || data[0][2] != "status" {
		return nil, errors.New("dont have headers")
	}
	return data[1:], nil
}

func (c *CSVCodec) ResetData() error {
	if err := c.File.Truncate(0); err != nil {
		return err
	}

	_, err := c.File.Seek(0, 0)
	if err != nil {
		return err
	}

	if err := c.Writer.Write([]string{"start", "end", "status"}); err != nil {
		return err
	}

	c.Writer.Flush()
	return nil
}
