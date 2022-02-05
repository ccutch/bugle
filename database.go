package bugle

import (
	"bytes"
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

type DBClient struct {
	ctx    context.Context
	client *mongo.Client
}

func NewDBClient(ctx context.Context) *DBClient {
	var db DBClient
	var err error
	var cancel context.CancelFunc

	options := options.Client().ApplyURI(os.Getenv("DB_URI"))
	db.ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	db.client, err = mongo.NewClient(options)
	if err != nil {
		log.Fatal(err)
	}

	err = db.client.Connect(db.ctx)
	if err != nil {
		log.Fatal(err)
	}

	return &db
}

func (db *DBClient) mongoDB() *mongo.Database { return db.client.Database("default") }
func (db *DBClient) lists() *mongo.Collection { return db.mongoDB().Collection("lists") }
func (db *DBClient) subs() *mongo.Collection  { return db.mongoDB().Collection("subs") }

func (db *DBClient) NewList(name string) string {
	res, err := db.lists().InsertOne(db.ctx, bson.M{"name": name})
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprint(res.InsertedID)
}

func (db *DBClient) NewSubscription(listName, name, address string) string {
	res, err := db.subs().InsertOne(db.ctx, bson.M{
		"name":     name,
		"address":  address,
		"listName": listName,
		"added":    time.Now().Format("01-01-1970"),
	})
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprint(res.InsertedID)
}

func (db *DBClient) GetSubscriptions(listName string) []string {
	cursor, err := db.subs().Find(db.ctx, bson.M{"listName": listName})
	if err != nil {
		log.Fatal(err)
	}

	var documents []bson.M
	err = cursor.All(db.ctx, &documents)
	if err != nil {
		log.Fatal(err)
	}

	var subscriptions []string
	var buff bytes.Buffer
	for _, doc := range documents {
		json.NewEncoder(&buff).Encode(&doc)
		subscriptions = append(subscriptions, buff.String())
		buff.Reset()
	}

	return subscriptions
}

func (db *DBClient) DeleteSubscription(listName, address string) {
	_, err := db.subs().DeleteOne(db.ctx, bson.M{"listName": listName, "address": address})
	if err != nil {
		log.Fatal(err)
	}
}
