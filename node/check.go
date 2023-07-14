package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nkn-server/log"
)

func CheckConnection(ip string) bool {
	url := fmt.Sprintf("http://%s:30003", ip)

	_, err := http.Get(url)

	return err == nil
}

func GetData(method, ip string) string {
	if CheckConnection(ip) {
		bytesRepresentation, err := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  method,
			"params":  map[string]interface{}{},
			"id":      1,
		})

		if err != nil {
			log.UpdateLog.Println(err)
		}

		url := fmt.Sprintf("http://%s:30003", ip)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))

		if err != nil {
			log.UpdateLog.Println(err)
		}

		jsn, err := io.ReadAll(resp.Body)
		if err != nil {
			log.UpdateLog.Println(err)
		}

		return string(jsn)
	}

	return `{"result": "-"}`
}
