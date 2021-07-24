package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoURI string
}

var config Config
var client *mongo.Client
var ctx context.Context

func ReadConfig() {
	file, _ := os.Open("./config/default.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error reading config file!")
		os.Exit(1)
	}
}

func ConnectToMongo() {
	ReadConfig()
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}

func CollectBirthdays() []bson.D {
	var people []bson.D
	collection := client.Database("myFirstDatabase").Collection("people")
	query, err := collection.Aggregate(ctx, bson.A{
		bson.M{
			"$redact": bson.M{
				"$cond": bson.A{
					bson.M{
						"$and": bson.A{
							bson.M{
								"$eq": bson.A{
									bson.M{"$month": "$birthday"},
									time.Now().Month()},
							},
							bson.M{
								"$eq": bson.A{
									bson.M{"$dayOfMonth": "$birthday"},
									time.Now().Day()},
							}}},
					"$$KEEP",
					"$$PRUNE"}}}})
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close(ctx)
	for query.Next(ctx) {
		var result bson.D
		err := query.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, result)
	}
	if err := query.Err(); err != nil {
		log.Fatal(err)
	}
	return people
}
