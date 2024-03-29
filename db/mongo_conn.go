package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func init() {
	client, err := GetMongoConn()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	MongoClient = client
}

type Mongo struct {
}

func GetMongoConn() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(
		"mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetDefaultDB() (*mongo.Database, error) {
	db, err := GetDB("hackathons_db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetDB(dbName string) (*mongo.Database, error) {
	client, err := GetMongoConn()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	db := client.Database(dbName)
	return db, nil
}

func GetCollection(colName string) (*mongo.Collection, error) {
	db, err := GetDefaultDB()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	col := db.Collection(colName)
	err = CreateIndexes()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	return col, nil
}

func (m Mongo) GetAccountCollection() (*mongo.Collection, error) {
	col, err := GetCollection("accounts")
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	err = CreateIndexes()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	return col, nil
}

func (m Mongo) GetParticipantCollection() (*mongo.Collection, error) {
	col, err := GetCollection("participants")

	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	err = CreateParticipantColIndexes()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	return col, err
}

func (m Mongo) GetTokenCollection() (*mongo.Collection, error) {
	col, err := GetCollection("tokens")
	return col, err
}

func (m Mongo) GetScoreCollection() (*mongo.Collection, error) {
	col, err := GetCollection("scores")
	return col, err
}

func CreateParticipantColIndexes() error {
	db, err := GetDefaultDB()
	if err != nil {
		return err
	}
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"participant_id", -1}},
		Options: options.Index().SetUnique(true),
	}
	indexName, err := db.Collection("participants").Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	fmt.Printf("\n%s\n", indexName)
	return nil
}

func CreateIndexes() error {
	db, err := GetDefaultDB()
	if err != nil {
		return err
	}
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", -1}},
		Options: options.Index().SetUnique(true),
	}
	indexName, err := db.Collection("accounts").Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{"token", -1}, {"token_type", 1}, {"token_type_value", 1}},
		Options: options.Index().SetUnique(true),
	}
	indexName, err = db.Collection("tokens").Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	fmt.Printf("%s", indexName)
	return nil
}
