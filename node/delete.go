package node

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"nkn-server/db"
	"nkn-server/log"
)

func Delete(generation int) {
	filter := bson.D{{"generation", generation}}
	_, err := db.NodeCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.MyLog.Println(err)
	}
}
