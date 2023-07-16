package node2

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataNode struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Ip              string             `bson:"ip" validate:"required"`
	Generation      int                `bson:"generation" validate:"required"`
	NodeStatus      string             `bson:"node_status"`
	Created         string             `bson:"created"`
	LastUpdate      string             `bson:"last_update"`
	LastOfflineTime string             `bson:"last_offline_time"`
}
