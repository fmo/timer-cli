package services

import (
	"encoding/csv"
	"os"
)

func GetDataFromCSV(file *os.File) ([][]string, error) {
	r := csv.NewReader(file)
	return r.ReadAll()
}
