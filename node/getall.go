package node

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nkn-server/db"
	"nkn-server/log"
)

func GetAll() []DataNode {
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"generation", 1}})
	cursor, err := db.NodeCollection.Find(context.TODO(), filter, opts)
	var results []DataNode
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.MyLog.Println(err)
	}

	return results
}
