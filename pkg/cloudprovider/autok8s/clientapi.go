package autok8s

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Worker struct {
	Name      string
	IpPrivate string
	IPPublic  string
	Type      string
}
type Data struct {
	Workers []Worker
}

func ReturnJson_by_provider(name string) (Worker, bool) {
	split_name := strings.Split(name, "//")
	return ReturnJson(split_name[1])

}

func ReturnJson(name string) (Worker, bool) {

	resp, err := http.Get("http://10.17.3.1:8050/api/cloud/vms") // DOTO fix static addres
	body, err := ioutil.ReadAll(resp.Body)
	var payload Data
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	for _, node := range payload.Workers {
		if node.Name == name {
			fmt.Println("%v", node.IPPublic)
			return node, false
		}
	}
	return Worker{}, true
}
