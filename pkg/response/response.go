package response

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Content string `json:"message"`
}

func Error(w http.ResponseWriter, err error, code int) {
	if err := json.NewEncoder(w).Encode(&Message{Content: err.Error()}); err != nil {
		http.Error(w, `{"message":"internal server error"}`, http.StatusInternalServerError)
	}
}

func JSON(w http.ResponseWriter, body interface{}, code int) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		Error(w, err, http.StatusInternalServerError)
	}
}
