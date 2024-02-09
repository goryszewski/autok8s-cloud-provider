package internal_client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func call_api(url string, payload map[string]string) []byte {
	json_payload, err := json.Marshal(payload)
	r, err := http.Post(url, "application/json", bytes.NewBuffer(json_payload))

	if err != nil {
		log.Fatal("Error|call_api|001:")
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal("Error|call_api|002:")
	}

	return body
}
