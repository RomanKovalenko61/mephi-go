package resp

import (
	"encoding/json"
	"net/http"
)

func ResponseJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Contetnt-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
