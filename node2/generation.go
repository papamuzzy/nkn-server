package node2

import (
	"context"
	"io/fs"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"nkn-server/config"
	"nkn-server/db"
	"nkn-server/log"
)

func GetGenerationId() int {
	var id int
	for id = 1; true; id++ {
		filter := bson.D{{"generation", id}}
		count, err := db.Node2Collection.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.MyLog.Println(err)
		}

		if count == 0 {
			break
		}
	}

	return id
}

func GetGenerationsCount() int {
	dir := config.DirRoot + "/public/generations/2/"
	count := 0

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			count++
		}

		return nil
	})

	if err != nil {
		log.MyLog.Println(err)
	}

	return count
}
