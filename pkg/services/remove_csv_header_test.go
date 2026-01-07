package services

import (
	"testing"
	"time"
)

func TestRemoveCsvHeader(t *testing.T) {
	start := time.Now()
	end := start.Add(30 * time.Minute)

	data := [][]string{
		{"start", "end", "status"},
		{start.Format(time.RFC3339), end.Format(time.RFC3339), "done"},
	}

	tests := []struct {
		name string
		data [][]string
		want [][]string
		err  bool
	}{
		{
			name: "valid data",
			data: data,
			want: data[1:],
			err:  false,
		},
		{
			name: "invalid data",
			data: [][]string{},
			want: [][]string{},
			err:  true,
		},
		{
			name: "empty valid data",
			data: data[0:1],
			want: [][]string{},
			err:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RemoveCSVHeader(tt.data)
			if tt.err && err == nil {
				t.Error("error expected")
			}
			if len(tt.want) != len(got) {
				t.Errorf("want: %d, got: %d", len(tt.want), len(got))
			}
		})
	}
}
