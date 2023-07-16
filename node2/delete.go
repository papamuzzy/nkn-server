package node2

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"nkn-server/db"
	"nkn-server/log"
)

func Delete(ip string) {
	filter := bson.D{{"ip", ip}}
	_, err := db.Node2Collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.MyLog.Println(err)
	}
}
