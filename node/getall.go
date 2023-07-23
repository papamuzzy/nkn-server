package node

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nkn-server/db"
	"nkn-server/log"
)

func GetAll(limit int64, offset int64) ([]DataNode, int) {
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"generation", 1}}).SetSkip(offset).SetLimit(limit)
	cursor, err := db.NodeCollection.Find(context.TODO(), filter, opts)
	var results []DataNode
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.MyLog.Println(err)
	}
	log.MyLog.Println(results)

	optsCount := options.Count().SetHint("_id_")
	total, err := db.NodeCollection.CountDocuments(context.TODO(), bson.D{}, optsCount)
	if err != nil {
		log.MyLog.Println(err)
	}

	return results, int(total)
}
