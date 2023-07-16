package node

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
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

		cursor, err := db.NodeCollection.Find(context.TODO(), bson.D{})
		if err != nil {
			log.UpdateLog.Println(err)
		}

		if err = cursor.All(context.TODO(), &nodes); err != nil {
			log.UpdateLog.Println(err)
		}
		for _, res := range nodes {
			now := time.Now()
			actualTime := xtime.ToStr(now)
			nowDay := now.Day()

			if res.LastUpdate == "-" {
				res.LastUpdate = actualTime
			}
			lastUpdate := xtime.FromStr(res.LastUpdate)
			lastUpdateDay := lastUpdate.Day()

			filter := bson.D{{"ip", res.Ip}}
			if CheckConnection(res.Ip) {
				nodeData := GetData("getnodestate", res.Ip)

				nodeState := gjson.Get(nodeData, "result").String()
				height := int(gjson.Get(nodeData, "result.height").Int())
				totalBlocks := int(gjson.Get(nodeData, "result.proposalSubmitted").Int())

				version := gjson.Get(GetData("getversion", res.Ip), "result").String()

				var blocksForToday int
				if res.LastBlockNumber > 0 {
					blocksForToday = totalBlocks - res.LastBlockNumber
				} else {
					blocksForToday = 0
				}

				state := gjson.Get(nodeState, "syncState").String()
				uptime := gjson.Get(nodeState, "uptime").Float()
				workTime := ""
				uptime /= 3600
				if uptime < 24 {
					workTime = fmt.Sprintf("%.1f h", uptime)
				} else {
					workTime = fmt.Sprintf("%.1f d", uptime/24)
				}

				if nowDay != lastUpdateDay {
					update := bson.D{
						{"$set",
							bson.D{
								{"height", height},
								{"version", version},
								{"work_time", workTime},
								{"mined_ever", totalBlocks},
								{"mined_today", 0},
								{"node_status", state},
								{"last_block_number", totalBlocks},
								{"last_update", actualTime},
							},
						},
					}

					_, err := db.NodeCollection.UpdateOne(context.TODO(), filter, update)
					if err != nil {
						log.UpdateLog.Println(err)
					}
				} else {
					update := bson.D{
						{"$set",
							bson.D{
								{"height", height},
								{"version", version},
								{"work_time", workTime},
								{"mined_ever", totalBlocks},
								{"mined_today", blocksForToday},
								{"node_status", state},
								{"last_update", actualTime},
							},
						},
					}

					_, err := db.NodeCollection.UpdateOne(context.TODO(), filter, update)
					if err != nil {
						log.UpdateLog.Println(err)
					}
				}
			} else {
				if res.NodeStatus != "OFFLINE" {
					update := bson.D{
						{"$set",
							bson.D{
								{"height", 0},
								{"version", "-"},
								{"work_time", "-"},
								{"mined_ever", 0},
								{"mined_today", 0},
								{"node_status", "OFFLINE"},
								{"last_update", xtime.ToStr(time.Now())},
								{"last_offline_time", xtime.ToStr(time.Now())},
							},
						},
					}

					_, err := db.NodeCollection.UpdateOne(context.TODO(), filter, update)
					if err != nil {
						log.UpdateLog.Println(err)
					}
				} else {
					t := xtime.FromStr(res.LastOfflineTime)
					delta := now.Sub(t)
					if delta.Hours() > 24 {
						block.NodesMutex.Lock()
						Delete(res.Ip)
						block.NodesMutex.Unlock()
					}
				}
			}
		}

		time.Sleep(30 * time.Second)
	}
}
