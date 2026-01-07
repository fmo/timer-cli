package services

import "errors"

func RemoveCSVHeader(data [][]string) ([][]string, error) {
	if len(data) < 1 {
		return nil, errors.New("csv has no header")
	}
	return data[1:], nil
}
