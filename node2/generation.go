package node2

import (
	"context"
	"fmt"
	"io/fs"
	mainlog "log"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"nkn-server/config"
	"nkn-server/db"
)

func GetGenerationId() int {
	count, err := db.Node2Collection.CountDocuments(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Println(err)
	}

	return int(count)
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
		mainlog.Println(err)
	}

	return count
}
