package webeh

import (
	"encoding/json"
	"net/http"
)

func httpResponseWithCode(w http.ResponseWriter, data []byte, code int) (int, error) {
	// WriteHeader sends an HTTP response header with status code.
	// If WriteHeader is not called explicitly, the first call to Write
	// will trigger an implicit WriteHeader(http.StatusOK).
	// Thus explicit calls to WriteHeader are mainly used to
	// send error codes.
	w.WriteHeader(code)
	return w.Write(data)
}

func ReplyJson(w http.ResponseWriter, ret interface{}) (int, error) {
	js, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json")
	return httpResponseWithCode(w, js, http.StatusOK)
}

func ReplyJsonWithCode(w http.ResponseWriter, ret interface{}, code int) (int, error) {
	js, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json")
	return httpResponseWithCode(w, js, code)
}
