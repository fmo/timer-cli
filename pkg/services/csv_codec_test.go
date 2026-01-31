package services_test

import (
	"log"
	"testing"

	"github.com/fmo/timer-cli/pkg/logger"
	"github.com/fmo/timer-cli/pkg/services"
)

func newTestCodec(t *testing.T) services.Persister {
	logger, err := logger.New()
	if err != nil {
		log.Fatal("cant initiate test due to logger setup")
	}

	codec, err := services.NewCSVCodec(logger)
	if err != nil {
		log.Fatal("cant initiate csv codec")
	}

	t.Cleanup(func() {
		_ = codec.ResetData()
	})

	return codec
}

func TestCreateHeader(t *testing.T) {
	codec := newTestCodec(t)

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

func TestSave(t *testing.T) {
	codec := newTestCodec(t)

	if err := codec.Save([]string{"saves", "any", "data"}); err != nil {
		t.Error("cant save")
	}

	data, err := codec.LoadData()
	if err != nil {
		t.Error("cant load data")
	}

	if len(data) != 1 {
		t.Errorf("expected %d row, got %d row", 1, len(data))
	}
}
