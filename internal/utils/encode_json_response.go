package utils

import (
	"encoding/json"
	"net/http"
)

func EncodeJSONRequest[Model interface{}](w http.ResponseWriter, data Model) error {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(data)

	if err != nil {
		return err
	}

	_, err = w.Write(resp)

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
