package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"nkn-server/log"
	"nkn-server/node"
	"strconv"
)

type RequestData struct {
	Ip           string `json:"ip"`
	GenerationId int    `json:"generation_id"`
	Exists       bool   `json:"exists"`
}

func NodesGet(w http.ResponseWriter, r *http.Request) {
	log.MyLog.Println("Route NodesGet Started")

	var resp map[string]interface{}
	resp = make(map[string]interface{})

	resp["Nodes"] = node.GetAll()

	response := UserResponse{Status: http.StatusOK, Message: "NodesGet success", Data: map[string]interface{}{"resp": resp}}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func NodeAdd(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)

	ip := gjson.Get(string(reqBody), "ip").String()
	generation := int(gjson.Get(string(reqBody), "generation_id").Int())

	if ip != "" && generation > 0 {
		node.Add(ip, generation)
	}

	response := UserResponse{Status: http.StatusOK, Message: "NodeAdd success", Data: map[string]interface{}{}}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func NodeMake(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	reqBody, _ := io.ReadAll(r.Body)
	/*var req RequestData
	json.Unmarshal(reqBody, &req)*/

	ip := gjson.Get(string(reqBody), "ip").String()

	if ip != "" {
		go node.Make(ip)
	}

	var resp map[string]interface{}
	resp = make(map[string]interface{})

	resp["Route"] = "NodeMake"
	resp["Request"] = fmt.Sprintf("%#v", r)
	resp["Ip"] = ip
	resp["Vars"] = vars

	response := UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"resp": resp}}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func NodeDelete(w http.ResponseWriter, r *http.Request) {
	log.MyLog.Println("Route NodeDelete Started")

	vars := mux.Vars(r)
	generationId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.MyLog.Println(err)
	}

	node.Delete(generationId)

	reqBody, _ := io.ReadAll(r.Body)
	/*var req RequestData
	json.Unmarshal(reqBody, &req)*/

	ip := gjson.Get(string(reqBody), "ip").String()
	exists := gjson.Get(string(reqBody), "exists").Bool()
	generation := gjson.Get(string(reqBody), "generation_id").Int()

	var resp map[string]interface{}
	resp = make(map[string]interface{})

	resp["Route"] = "NodeDelete"
	resp["Request"] = fmt.Sprintf("%#v", r)
	resp["Ip"] = ip
	resp["Generation"] = generation
	resp["Exists"] = exists
	resp["Vars"] = vars

	response := UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"resp": resp}}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
