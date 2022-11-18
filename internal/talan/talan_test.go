package talan

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"log"
	"nft/config"
	"testing"
)

func initConfig(down fx.Shutdowner) {
	config.InitConfigs(down, "../../")
}

func TestGenerateAddress(t *testing.T) {
	err := fxtest.New(t, fx.Invoke(initConfig)).Start(context.Background())
	if err != nil {
		return
	}
	repository := NewTalanRepository(TalanRepositoryParams{})
	_, err = repository.GenerateAddress(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestGetBalance(t *testing.T) {
	err := fxtest.New(t, fx.Invoke(initConfig)).Start(context.Background())
	if err != nil {
		return
	}
	repository := NewTalanRepository(TalanRepositoryParams{})
	_, err = repository.GetBalance(context.Background(), "TfqaEUs2GBdJoPFFYPE9eHn8jR9qgY5ZsR")
	if err != nil {
		t.Error(err)
	}
}

func TestGetTransactions(t *testing.T) {
	err := fxtest.New(t, fx.Invoke(initConfig)).Start(context.Background())
	if err != nil {
		return
	}
	repository := NewTalanRepository(TalanRepositoryParams{})
	txs, err := repository.GetTransactions(context.Background(), "TcW3npnMSKNwjUbmKHaT76njayWbKBytK1")
	if err != nil {
		t.Error(err)
	}
	log.Println(txs)
}
