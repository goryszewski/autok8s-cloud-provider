package internal_client

import (
	"encoding/json"
	"log"
)

func (i *Autok8sClient) GetIPByNodeName(nodename string) Worker {

	payload := map[string]string{"function": "GetNodeByHostname", "hostname": nodename}

	body := call_api(i.URL+"/api/v1/k8s/node", payload)

	var node Worker

	err := json.Unmarshal(body, &node)
	if err != nil {
		log.Fatal("Error Unmarshal")
	}
	return node
}
