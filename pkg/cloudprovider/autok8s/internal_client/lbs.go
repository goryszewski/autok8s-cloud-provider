package internal_client

import (
	"encoding/json"
	"log"
)

func (i *Autok8sClient) GetFreeLB() LB {

	payload := map[string]string{"function": "GetFirstFreeLB"}

	body := call_api(i.URL+"/api/v1/k8s/lb", payload)

	var lb LB

	err := json.Unmarshal(body, &lb)

	if err != nil {
		log.Fatal("Error 003")
	}

	return lb
}

func (i *Autok8sClient) BindLB(ip string, service string) bool {
	payload := map[string]string{"function": "BindLB", "ip": ip, "service": service}
	body := call_api(i.URL+"/api/v1/k8s/lb", payload)
	var lb LB

	err := json.Unmarshal(body, &lb)

	if err != nil {
		log.Fatal("Error 003")
	}

	return true
}
func (i *Autok8sClient) UnBindLB(service string) bool {
	payload := map[string]string{"function": "UnBindLB", "service": service}
	body := call_api(i.URL+"/api/v1/k8s/lb", payload)
	var lb LB

	err := json.Unmarshal(body, &lb)

	if err != nil {
		log.Fatal("Error 003")
	}
	return true
}
