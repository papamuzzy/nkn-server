package node

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataNode struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Ip              string             `bson:"ip" validate:"required"`
	Generation      int                `bson:"generation" validate:"required"`
	Height          int                `bson:"height"`
	Version         string             `bson:"version"`
	WorkTime        string             `bson:"work_time"`
	MinedEver       int                `bson:"mined_ever"`
	MinedToday      int                `bson:"mined_today"`
	NodeStatus      string             `bson:"node_status"`
	LastBlockNumber int                `bson:"last_block_number"`
	LastUpdate      string             `bson:"last_update"`
	LastOfflineTime string             `bson:"last_offline_time"`
}
