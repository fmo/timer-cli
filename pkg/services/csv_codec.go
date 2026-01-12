package services

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"time"
)

var ErrNoDataInCSV = errors.New("no data in csv")

type TaskStorer interface {
	Save(Task) error
	Update(Task) error
}

type CSVCodec struct {
	File   *os.File
	Writer *csv.Writer
	Reader *csv.Reader
}

func NewCSVCodec(f *os.File) (*CSVCodec, error) {
	writer := csv.NewWriter(f)
	reader := csv.NewReader(f)

	codec := &CSVCodec{
		File:   f,
		Writer: writer,
		Reader: reader,
	}

	if data, err := codec.Load(); data == nil {
		if err != nil {
			log.Printf("cant load csv: %v", err)
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

func (c *CSVCodec) Update(task Task) error {
	return nil
}

func (c *CSVCodec) Load() ([][]string, error) {
	data, err := c.Reader.ReadAll()
	if err != nil {
		return nil, err
	}
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
