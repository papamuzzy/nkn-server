package handlers

import (
	"encoding/json"
	"net/http"
	"nkn-server/node"
)

func GenerationCount(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	resp = make(map[string]interface{})

	resp["Count"] = node.GetGenerationsCount()

	response := UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"resp": resp}}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
