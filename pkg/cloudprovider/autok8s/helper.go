package autok8s

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"k8s.io/klog"
)

type Worker struct {
    Name string
    Ip   string
	Type string
}
type Data struct {
    Workers []Worker
}

func ReturnJson_by_provider ( name string) (Worker, bool) {
	split_name :=strings.Split(name, "//")
	return ReturnJson( split_name[1])

}

func ReturnJson( name string) (Worker, bool) {
	content, err := ioutil.ReadFile("./payload.json")

	if err != nil {
		klog.V(5).Infof("Error import json: %v", err )
	}

	var payload Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
        log.Fatal("Error during Unmarshal(): ", err)
    }
	for _, node := range payload.Workers {
		if node.Name == name{
			return node,false
		}
		// fmt.Println("%v", node)
	}
	return Worker{} , true
}
