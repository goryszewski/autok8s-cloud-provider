package autok8s

import (
	v1 "k8s.io/api/core/v1"
	// libvirtApiClient "github.com/goryszewski/libvirtApi-client/libvirtApiClient"
)

func PreperbindServiceLB(service *v1.Service) libvirtApiClient.bindServiceLB {
	var ports []libvirtApiClient.portlb

	for _, port := range service.Spec.Ports {
		pre_port := libvirtApiClient.portlb{Name: port.Name, Protocol: string(port.Protocol), Port: int(port.Port), NodePort: int(port.NodePort)}
		ports = append(ports, pre_port)
	}

	bind_payload := libvirtApiClient.bindServiceLB{
		ports: ports,
	}

	return bind_payload
}
