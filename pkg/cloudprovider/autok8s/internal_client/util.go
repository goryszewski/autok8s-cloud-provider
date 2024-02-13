package internal_client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (i *Autok8sClient) get_token(user string, password string) []byte {
	payload := map[string]string{"user": user, "password": password}
	json_payload, _ := json.Marshal(payload)
	url := i.URL + "/api/v1/auth"
	req, err := http.NewRequest("Post", url, bytes.NewBuffer(json_payload))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, erro := client.Do(req)

	if err != nil || erro != nil {
		log.Fatal("Error|token|001:")
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal("Error|token|002:")
	}

	return body
}

func call_api(url string, payload map[string]string, token string) []byte {
	json_payload, err := json.Marshal(payload)

	req, err := http.NewRequest("Post", url, bytes.NewBuffer(json_payload))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, erro := client.Do(req)

	if err != nil || erro != nil {
		log.Fatal("Error|call_api|001:")
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal("Error|call_api|002:")
	}

	return body
}

func get_token(url string)
