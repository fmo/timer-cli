package services_test

import (
	"log"
	"testing"

	"github.com/fmo/timer-cli/pkg/logger"
	"github.com/fmo/timer-cli/pkg/services"
)

func newTestCodec() services.Persister {
	logger, err := logger.New()
	if err != nil {
		log.Fatal("cant initiate test due to logger setup")
	}

	codec, err := services.NewCSVCodec(logger)
	if err != nil {
		log.Fatal("cant initiate csv codec")
	}

	codec.ResetData()

	return codec
}

func TestSave(t *testing.T) {
	codec := newTestCodec()

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
