package script

import (
	"log"
	"nkn-server/config"
	"os"
	"strconv"
	"strings"
)

func GetString(fName string, generationId int, logger *log.Logger) string {
	file := config.DirRoot + "/scriptfiles" + fName
	dat, err := os.ReadFile(file)
	if err != nil {
		logger.Println(err)
	}

	return strings.ReplaceAll(string(dat), "%GenId", strconv.Itoa(generationId))
}
