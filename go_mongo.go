package go_mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx context.Context

func ConnectToMongo(mongoUri string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func CollectBirthdays(client *mongo.Client) []bson.D {
	collection := client.Database("myFirstDatabase").Collection("people")
	query, err := getBirthdaysFromCollection(collection)
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close(ctx)
	people := getBsonDFromQuery(query)
	return people
}

func CollectEvents(client *mongo.Client) []bson.D {
	collection := client.Database("myFirstDatabase").Collection("events")
	query, err := getEventsFromCollection(collection)
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close(ctx)
	events := getBsonDFromQuery(query)
	return events
}

// Get all events that were added the current day
// Query this endpoint every day, maybe every hour?
func getEventsFromCollection(collection *mongo.Collection) (*mongo.Cursor, error) {
	return collection.Find(ctx, bson.M{"date_added": bson.M{"$gte": time.Now()}})
}

func getBsonDFromQuery(query *mongo.Cursor) []bson.D {
	var result []bson.D
	for query.Next(ctx) {
		var curr bson.D
		err := query.Decode(&curr)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, curr)
	}
	if err := query.Err(); err != nil {
		log.Fatal(err)
	}
	return result
}

func getBirthdaysFromCollection(collection *mongo.Collection) (*mongo.Cursor, error) {
	return collection.Aggregate(ctx, bson.A{
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
}
