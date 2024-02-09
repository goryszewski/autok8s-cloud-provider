package autok8s

import "autok8s.io/autok8s/pkg/cloudprovider/autok8s/internal_client"

func newAutok8sClient() *internal_client.Autok8sClient {
	return &internal_client.Autok8sClient{URL: "http://10.17.3.1:8050"} // move to config file
}
