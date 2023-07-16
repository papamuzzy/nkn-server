package node2

import (
	"fmt"
	"net/http"
)

func CheckConnection(ip string) bool {
	url := fmt.Sprintf("http://%s:4449", ip)

	_, err := http.Get(url)

	return err == nil
}
