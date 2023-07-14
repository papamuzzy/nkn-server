package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

func CatchAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	reqBody, _ := io.ReadAll(r.Body)
	/*var req RequestData
	json.Unmarshal(reqBody, &req)*/

	ip := gjson.Get(string(reqBody), "ip").String()
	exists := gjson.Get(string(reqBody), "exists").Bool()
	generation := gjson.Get(string(reqBody), "generation_id").Int()

	var resp map[string]interface{}
	resp = make(map[string]interface{})

	resp["Route"] = "CatchAll"
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
