package autok8s

import (
	"encoding/json"
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
