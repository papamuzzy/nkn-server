package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"nkn-server/log"
	"nkn-server/node2"
	"strings"
)

func Nodes2Get(w http.ResponseWriter, r *http.Request) {
	log.MyLog.Println("Route NodesGet Started")

	var resp map[string]interface{}
	resp = make(map[string]interface{})

	resp["Nodes"] = node2.GetAll()

	response := UserResponse{Status: http.StatusOK, Message: "NodesGet2 success", Data: map[string]interface{}{"resp": resp}}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func Node2Add(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)

	ip := gjson.Get(string(reqBody), "ip").String()
	generation := int(gjson.Get(string(reqBody), "generation_id").Int())

	if ip != "" && generation > 0 {
		node2.Add(ip, generation)
	}

	response := UserResponse{Status: http.StatusOK, Message: "NodeAdd2 success", Data: map[string]interface{}{}}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func Node2Make(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)

	ip := strings.Split(strings.TrimSpace(gjson.Get(string(reqBody), "ip").String()), " ")[0]

	go node2.Make(ip)

	var resp map[string]interface{}
	resp = make(map[string]interface{})

	resp["Route"] = "NodeMake2"
	resp["Request"] = fmt.Sprintf("%#v", r)
	resp["Ip"] = ip

	response := UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"resp": resp}}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func Node2Delete(w http.ResponseWriter, r *http.Request) {
	log.MyLog.Println("Route NodeDelete2 Started")

	reqBody, _ := io.ReadAll(r.Body)
	ip := strings.Split(strings.TrimSpace(gjson.Get(string(reqBody), "ip").String()), " ")[0]

	node2.Delete(ip)

	var resp map[string]interface{}
	resp = make(map[string]interface{})

	resp["Route"] = "NodeDelete2"
	resp["Request"] = fmt.Sprintf("%#v", r)
	resp["Ip"] = ip

	response := UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"resp": resp}}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
