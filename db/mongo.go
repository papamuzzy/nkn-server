package db

import (
	"context"
	"nkn-server/config"
	"nkn-server/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DataBase *mongo.Database
var NodeCollection *mongo.Collection
var Node2Collection *mongo.Collection

func Start() {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MongoUri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	Client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.MyLog.Println(err)
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := Client.Database(config.MongoBase).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		log.MyLog.Println(err)
	}
	log.MyLog.Println("Pinged your deployment. You successfully connected to MongoDB!")

	DataBase = Client.Database(config.MongoBase)
	NodeCollection = DataBase.Collection(config.MongoCollection)
	Node2Collection = DataBase.Collection(config.MongoCollection)
}

func Stop() {
	if err := Client.Disconnect(context.TODO()); err != nil {
		log.MyLog.Println(err)
	}
}
