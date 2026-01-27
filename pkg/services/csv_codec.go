package services

import (
	"encoding/csv"
	"os"

	"github.com/fmo/timer-cli/pkg/logger"
)

type Persister interface {
	Save([]string) error
	Update([]string) error
	LoadData() ([][]string, error)
	ResetData() error
	CreateHeader() error
}

type CSVCodec struct {
	File   *os.File
	Logger logger.Logger
	Writer *csv.Writer
	Reader *csv.Reader
}

func NewCSVCodec(f *os.File, logger logger.Logger) *CSVCodec {
	writer := csv.NewWriter(f)
	reader := csv.NewReader(f)

	return &CSVCodec{
		File:   f,
		Logger: logger,
		Writer: writer,
		Reader: reader,
	}
}

func (c *CSVCodec) saveHeader() error {
	if err := c.Writer.Write([]string{"start", "end", "status"}); err != nil {
		return err
	}
	c.Writer.Flush()
	return nil
}

func (c *CSVCodec) CreateHeader() error {
	data, err := c.Load()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		if err := c.saveHeader(); err != nil {
			return err
		}
	}
	return nil
}

func (c *CSVCodec) Save(row []string) error {
	if err := c.Writer.Write(row); err != nil {
		return err
	}

	c.Writer.Flush()
	return nil
}

func (c *CSVCodec) Update(rowToUpdate []string) error {
	data, err := c.LoadData()
	if err != nil {
		return err
	}

	if err := c.ResetData(); err != nil {
		return err
	}

	if err := c.saveHeader(); err != nil {
		return err
	}

	for _, row := range data {
		if row[0] == rowToUpdate[0] {
			c.Writer.Write(rowToUpdate)
		} else {
			c.Writer.Write(row)
		}
		c.Writer.Flush()
	}

	return nil
}

// Load whole csv file
func (c *CSVCodec) Load() ([][]string, error) {
	c.File.Seek(0, 0)

	return c.Reader.ReadAll()
}

// LoadData assumes first line is already header and loads without it
func (c *CSVCodec) LoadData() ([][]string, error) {
	c.File.Seek(0, 0)
	data, err := c.Reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return data[1:], nil
}

// ResetData allows reset data, NewCodec creates header if it does not there
func (c *CSVCodec) ResetData() error {
	if err := c.File.Truncate(0); err != nil {
		return err
	}

	_, err := c.File.Seek(0, 0)
	if err != nil {
		return err
	}

	return nil
}
