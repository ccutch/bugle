package bugle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBClient struct {
	ctx    context.Context
	err    error
	client *mongo.Client
}

func NewDBClient(ctx context.Context) *DBClient {
	var db DBClient
	var cancel context.CancelFunc

	options := options.Client().ApplyURI(os.Getenv("DB_URI"))
	db.ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	db.client, db.err = mongo.NewClient(options)
	if db.err != nil {
		return &db
	}

	db.err = db.client.Connect(db.ctx)
	return &db
}

func (db *DBClient) mongoDB() *mongo.Database { return db.client.Database("default") }
func (db *DBClient) lists() *mongo.Collection { return db.mongoDB().Collection("lists") }
func (db *DBClient) subs() *mongo.Collection  { return db.mongoDB().Collection("subs") }

func (db *DBClient) NewList(name string) string {
	if db.err != nil {
		return ""
	}

	var res *mongo.InsertOneResult
	res, db.err = db.lists().InsertOne(db.ctx, bson.M{"name": name})
	return fmt.Sprint(res.InsertedID)
}

func (db *DBClient) NewSubscription(listName, name, address string) string {
	if db.err != nil {
		return ""
	}

	var res *mongo.InsertOneResult
	res, db.err = db.subs().InsertOne(db.ctx, bson.M{
		"name":     name,
		"address":  address,
		"listName": listName,
		"added":    time.Now().Format("01-01-1970"),
	})

	return fmt.Sprint(res.InsertedID)
}

func (db *DBClient) GetSubscriptions(listName string) []string {
	if db.err != nil {
		return []string{}
	}

	var cursor *mongo.Cursor
	cursor, db.err = db.subs().Find(db.ctx, bson.M{"listName": listName})
	if db.err != nil {
		return []string{}
	}

	var documents []bson.M
	db.err = cursor.All(db.ctx, &documents)
	if db.err != nil {
		return []string{}
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
	if db.err != nil {
		return
	}

	_, db.err = db.subs().DeleteOne(db.ctx, bson.M{"listName": listName, "address": address})
}
