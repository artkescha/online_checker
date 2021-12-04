package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, statusCode int, value interface{}, tags ...string) {
	var err error
	var response_ []byte
	if len(tags) > 0 {
		data := make(map[string]interface{})
		for _, tag := range tags {
			data[tag] = value
		}
		response_, err = json.Marshal(data)
	} else {
		response_, err = json.Marshal(value)
	}
	if err != nil {
		log.Printf("json marshal error: %s", err)
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(statusCode)
	_, err = w.Write(response_)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
