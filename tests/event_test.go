package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/DavidBuzatu-Marian/go_mongo"
)

type Config struct {
	MongoURI string
}

var config Config

func ReadConfig() {
	file, _ := os.Open("../config/default.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error reading config file!")
		os.Exit(1)
	}
}

func TestCollectEvents(t *testing.T) {
	ReadConfig()
	go_mongo.ConnectToMongo(config.MongoURI)
	client := go_mongo.ConnectToMongo(config.MongoURI)
	events := go_mongo.CollectEvents(client)
	for _, val := range events {
		t.Log(val)
	}
}
