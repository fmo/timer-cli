package services

func RemoveCSVHeader(data [][]string) [][]string {
	if data == nil {
		return data
	}

	if data[0][0] == "start" {
		return data[1:]
	}

	return data
}
