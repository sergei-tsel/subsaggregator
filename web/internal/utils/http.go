package utils

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, data interface{}, status ...int) {
	code := http.StatusOK

	if len(status) > 0 {
		code = status[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}
