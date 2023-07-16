package node2

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"nkn-server/block"
	"nkn-server/db"
	"nkn-server/log"
	"nkn-server/xtime"
	"time"
)

func UpdateBase() {
	for {
		var nodes []DataNode

		cursor, err := db.Node2Collection.Find(context.TODO(), bson.D{})
		if err != nil {
			log.UpdateLog.Println(err)
		}

		if err = cursor.All(context.TODO(), &nodes); err != nil {
			log.UpdateLog.Println(err)
		}
		for _, res := range nodes {
			filter := bson.D{{"ip", res.Ip}}
			if CheckConnection(res.Ip) {
				update := bson.D{
					{"$set",
						bson.D{
							{"node_status", "ACTIVE"},
							{"last_update", xtime.ToStr(time.Now())},
						},
					},
				}

				_, err := db.Node2Collection.UpdateOne(context.TODO(), filter, update)
				if err != nil {
					log.UpdateLog.Println(err)
				}
			} else {
				if res.NodeStatus != "OFFLINE" {
					update := bson.D{
						{"$set",
							bson.D{
								{"node_status", "OFFLINE"},
								{"last_update", xtime.ToStr(time.Now())},
								{"last_offline_time", xtime.ToStr(time.Now())},
							},
						},
					}

					_, err := db.Node2Collection.UpdateOne(context.TODO(), filter, update)
					if err != nil {
						log.UpdateLog.Println(err)
					}
				} else {
					t := xtime.FromStr(res.LastOfflineTime)
					delta := time.Now().Sub(t)
					if delta.Hours() > 24 {
						block.Nodes2Mutex.Lock()
						Delete(res.Ip)
						block.Nodes2Mutex.Unlock()
					}
				}
			}
		}

		time.Sleep(10 * time.Second)
	}
}
