package services_test

import (
	"testing"

	"github.com/fmo/timer-cli/pkg/logger"
	"github.com/fmo/timer-cli/pkg/services"
)

func TestCreateHeader(t *testing.T) {
	logger, err := logger.New()
	if err != nil {
		t.Error("unexpected error while initiating logger")
	}

	codec, err := services.NewCSVCodec(logger)
	if err != nil {
		t.Error("cant initiate csv codec")
	}

	if err := codec.CreateHeader(); err != nil {
		t.Error("creating header unexpectedly failed")
	}

	data, err := codec.LoadData()
	if err != nil {
		t.Error("loading data failed")
	}

	if len(data) != 1 {
		t.Error("should have only header record")
	}

	if data[0][0] != "start" || data[0][1] != "end" || data[0][2] != "status" {
		t.Error("header has not have expected data")
	}
}
