package internal_client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (i *autok8sClient) GetIPByNodeName(nodename string) Worker {
	r, err := http.Get(i.URL + "/api/clud/vms")
	if err != nil {
		log.Fatal("Error 001:")
	}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal("Error 002:")
	}

	var node Worker

	err = json.Unmarshal(body, &node)
	if err != nil {
		log.Fatal("Error 003")
	}
	return node
}
