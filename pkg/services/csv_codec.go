package services

import (
	"encoding/csv"
	"os"

	"github.com/fmo/timer-cli/pkg/logger"
)

const taskFile = "tasks.csv"

type Persister interface {
	Save([]string) error
	Update([]string) error
	ResetData() error
	CreateHeader() error
	LoadData() ([][]string, error)
}

type CSVCodec struct {
	file   *os.File
	logger logger.Logger
	writer *csv.Writer
	reader *csv.Reader
	header []string
}

func NewCSVCodec(logger logger.Logger) (*CSVCodec, error) {
	f, err := os.OpenFile(taskFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &CSVCodec{
		file:   f,
		logger: logger,
		writer: csv.NewWriter(f),
		reader: csv.NewReader(f),
		header: []string{"start", "end", "status"}}, nil
}

// CreateHeader creates header if file has no data. If there is data, this does not
// check if its a actual header or not.
func (c *CSVCodec) CreateHeader() error {
	data, err := c.LoadData()
	if err != nil {
		return err
	}

	if len(data) == 0 {
		if err := c.writer.Write(c.header); err != nil {
			return err
		}
		c.writer.Flush()
	}

	return nil
}

// Save new row to the csv
func (c *CSVCodec) Save(row []string) error {
	if err := c.writer.Write(row); err != nil {
		return err
	}

	c.writer.Flush()
	return nil
}

// Update rewrites whole tasks list with the new record
func (c *CSVCodec) Update(rowToUpdate []string) error {
	data, err := c.LoadData()
	if err != nil {
		return err
	}

	// remove header
	data = data[1:]

	if err := c.ResetData(); err != nil {
		return err
	}

	if err := c.CreateHeader(); err != nil {
		return err
	}

	for _, row := range data {
		if row[0] == rowToUpdate[0] {
			c.writer.Write(rowToUpdate)
		} else {
			c.writer.Write(row)
		}
		c.writer.Flush()
	}

	return nil
}

// ResetData allows reset data
func (c *CSVCodec) ResetData() error {
	if err := c.file.Truncate(0); err != nil {
		return err
	}

	_, err := c.file.Seek(0, 0)
	if err != nil {
		return err
	}

	return nil
}

// LoadData loads whole data including header
func (c *CSVCodec) LoadData() ([][]string, error) {
	// move cursor to the top
	c.file.Seek(0, 0)

	data, err := c.reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}
