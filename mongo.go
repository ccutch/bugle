package bugle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBClient struct {
	*mongo.Client
}

func Mongo(url string) (*DBClient, error) {
	options := options.Client().ApplyURI(url)
	client, err := mongo.NewClient(options)
	if err != nil {
		return nil, err
	}

	db := DBClient{client}
	return &db, db.Connect(context.Background())
}

func (db *DBClient) mongoDB() *mongo.Database { return db.Database("default") }
func (db *DBClient) lists() *mongo.Collection { return db.mongoDB().Collection("lists") }
func (db *DBClient) subs() *mongo.Collection  { return db.mongoDB().Collection("subs") }

func (db *DBClient) NewList(name string) (string, error) {
	res, err := db.lists().InsertOne(context.TODO(), bson.M{"name": name})
	return fmt.Sprint(res.InsertedID), err
}

func (db *DBClient) NewSubscription(listName, name, address string) (string, error) {
	res, err := db.subs().InsertOne(context.TODO(), bson.M{
		"name":     name,
		"address":  address,
		"listName": listName,
		"added":    time.Now().Format("01-01-1970"),
	})
	return fmt.Sprint(res.InsertedID), err
}

func (db *DBClient) GetSubscriptions(listName string) ([]string, error) {
	cursor, err := db.subs().Find(context.TODO(), bson.M{"listName": listName})
	if err != nil {
		return []string{}, err
	}

	var documents []bson.M
	if err = cursor.All(context.TODO(), &documents); err != nil {
		return []string{}, err
	}

	var subscriptions []string
	var buff bytes.Buffer
	enc := json.NewEncoder(&buff)
	for _, doc := range documents {
		if err = enc.Encode(&doc); err != nil {
			break
		}
		subscriptions = append(subscriptions, buff.String())
		buff.Reset()
	}

	return subscriptions, err
}

func (db *DBClient) DeleteSubscription(listName, address string) error {
	_, err := db.subs().DeleteOne(context.TODO(), bson.M{
		"listName": listName,
		"address":  address,
	})
	return err
}
