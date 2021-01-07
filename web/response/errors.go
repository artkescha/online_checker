package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseParamsError struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Value    string `json:"value"`
	Message  string `json:"msg"`
}

func NewResponseParamsErrors(location string, param string, value string, err error) ResponseParamsError {
	return ResponseParamsError{
		Location: location,
		Param:    param,
		Value:    value,
		Message:  err.Error(),
	}
}

func WriteParamsErrors(w http.ResponseWriter, statusCode int, errors ...ResponseParamsError) {
	data := make(map[string]interface{})
	data["errors"] = errors

	response, err := json.Marshal(data)

	if err != nil {
		log.Printf("json marshal error: %v", errors)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(response)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	data := make(map[string]interface{})
	data["message"] = err.Error()
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(statusCode)
	w.Write(response)
}
