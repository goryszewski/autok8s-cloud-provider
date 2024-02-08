package autok8s

// moving to next step learing - internal_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"k8s.io/klog"
)

type Ip struct {
	Public  string
	Private string
}

type Worker struct {
	Name string
	IP   Ip
	Type string
}
type Data struct {
	Workers []Worker
}
type LB struct {
	Id           int
	Ip           string
	service_name string
}

func ReturnJson_by_provider(name string) (Worker, bool) {
	split_name := strings.Split(name, "//")
	return ReturnJson(split_name[1])

}

func ReturnJson(nodeName string) (Worker, bool) {

	resp, err := http.Get("http://10.17.3.1:8050/api/cloud/vms") // DOTO fix static addres
	body, err := ioutil.ReadAll(resp.Body)

	var payload Data
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	for _, node := range payload.Workers {
		if node.Name == nodeName {
			klog.V(5).Infof("InstanceID(%v) Data(%v)", nodeName, node)
			return node, false
		}
	}
	return Worker{}, true
}

func lbFree() LB {
	resp, err := http.Get("http://10.17.3.1:8050/api/lbs?filter=2") // DOTO fix static addres
	body, err := ioutil.ReadAll(resp.Body)
	var payload LB
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return payload
}

func lbUnbind(ip string) {
	data := []byte(`{"Operation": "Unbind","ip": ip}`)

	url := fmt.Sprintf("http://10.17.3.1:8050/api/lbs")

	resp, err := http.Post(url, "", bytes.NewBuffer(data))
	body, err := ioutil.ReadAll(resp.Body)
	var payload LB
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

}
func lbBind(id int, service_name string) string {
	data := []byte(`{"Operation": "Bind","Id":id ,"Service_name":service_name}`)

	url := fmt.Sprintf("http://10.17.3.1:8050/api/lbs")
	klog.V(5).Info("LB URL(%v)", url)
	resp, err := http.Post(url, "", bytes.NewBuffer(data))
	body, err := ioutil.ReadAll(resp.Body)
	var payload LB
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload.Ip
}
