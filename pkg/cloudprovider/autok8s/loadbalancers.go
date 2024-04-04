package autok8s

import (
	"context"

	libvirtApiClient "github.com/goryszewski/libvirtApi-client/libvirtApiClient"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

type loadbalancers struct {
	client *libvirtApiClient.Client
}

func newLoadBalancers(c *libvirtApiClient.Client) cloudprovider.LoadBalancer {
	return &loadbalancers{
		c,
	}
}

func prepServiceLoadBalancerPayload(service *v1.Service, nodes []*v1.Node) libvirtApiClient.LoadBalancer {
	var ports []libvirtApiClient.Port_Service
	var _nodes []libvirtApiClient.Node

	for _, port := range service.Spec.Ports {
		pre_port := libvirtApiClient.Port_Service{Name: port.Name, Protocol: string(port.Protocol), Port: int(port.Port), NodePort: int(port.NodePort)}
		ports = append(ports, pre_port)
	}

	if nodes != nil {
		for _, item := range nodes {
			node := libvirtApiClient.Node{Name: item.Name}
			for _, address := range item.Status.Addresses {

				if address.Type == "InternalIP" {
					node.Internal = address.Address
				} else if address.Type == "ExternalIP" {
					node.External = address.Address
				}
			}
			_nodes = append(_nodes, node)

		}
	} else {
		node := libvirtApiClient.Node{Name: "", External: "", Internal: ""}
		_nodes = append(_nodes, node)
	}

	bind_payload := libvirtApiClient.LoadBalancer{
		Ports:     ports,
		Name:      service.Name,
		Namespace: service.Namespace,
		Nodes:     _nodes,
	}

	return bind_payload
}

// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (lb *loadbalancers) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	klog.V(5).Infof("[GetLoadBalancer]")
	prepServiceLoadBalancer := prepServiceLoadBalancerPayload(service, nil)
	loadbalancer, exist, _ := lb.client.GetLoadBalancer(prepServiceLoadBalancer)
	if exist {
		var LoadBalancerI []v1.LoadBalancerIngress = []v1.LoadBalancerIngress{v1.LoadBalancerIngress{IP: loadbalancer.Ip}}
		return &v1.LoadBalancerStatus{Ingress: LoadBalancerI}, true, nil
	}
	return nil, false, nil
}

// GetLoadBalancerName returns the name of the load balancer. Implementations must treat the
// *v1.Service parameter as read-only and not modify it.
func (lb *loadbalancers) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	klog.V(5).Infof("GetLoadBalancerName ---")
	prepServiceLoadBalancer := prepServiceLoadBalancerPayload(service, nil)
	loadbalancer, _, _ := lb.client.GetLoadBalancer(prepServiceLoadBalancer)
	return loadbalancer.Ip
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one. Returns the status of the balancer
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (lb *loadbalancers) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	klog.V(5).Infof("EnsureLoadBalancer")
	ip := ""
	prepServiceLoadBalancer := prepServiceLoadBalancerPayload(service, nodes)

	loadbalancer, exist, err := lb.client.GetLoadBalancer(prepServiceLoadBalancer)

	if err != nil {
		klog.V(5).Infof("EnsureLoadBalancer : GetLoadBalancer (%v)", err)
		return nil, err
	}
	if exist {
		ip = loadbalancer.Ip
		klog.V(5).Infof("EnsureLoadBalancer - LB - do update")

		err = lb.client.UpdateLoadBalancer(prepServiceLoadBalancer)
		if err != nil {
			klog.V(5).Infof("ERROR: UpdateLoadBalancer (%v)", err)
			return nil, err
		}
	} else {
		klog.V(5).Infof("EnsureLoadBalancer - LB not exist")

		ip, err = lb.client.CreateLoadBalancer(prepServiceLoadBalancer)
		if err != nil {
			klog.V(5).Infof("EnsureLoadBalancer : CreateLoadBalancer (%v)", err)
			return nil, err
		}
	}

	var LoadBalancerI []v1.LoadBalancerIngress = []v1.LoadBalancerIngress{v1.LoadBalancerIngress{IP: ip}}

	return &v1.LoadBalancerStatus{Ingress: LoadBalancerI}, nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (lb *loadbalancers) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	klog.V(5).Infof("[UpdateLoadBalancer]")

	prepServiceLoadBalancer := prepServiceLoadBalancerPayload(service, nil)

	err := lb.client.UpdateLoadBalancer(prepServiceLoadBalancer)

	if err != nil {
		klog.V(5).Infof("[UpdateLoadBalancer][ERROR] (%v)", err)
		return err
	}

	return nil
}

// EnsureLoadBalancerDeleted deletes the specified load balancer if it
// exists, returning nil if the load balancer specified either didn't exist or
// was successfully deleted.
// This construction is useful because many cloud providers' load balancers
// have multiple underlying components, meaning a Get could say that the LB
// doesn't exist even if some part of it is still laying around.
// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (lb *loadbalancers) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	klog.V(5).Infof("[EnsureLoadBalancerDeleted]")
	prepServiceLoadBalancer := prepServiceLoadBalancerPayload(service, nil)

	ip, exist, err := lb.client.GetLoadBalancer(prepServiceLoadBalancer)
	if err != nil {
		klog.V(5).Infof("[EnsureLoadBalancerDeleted][GetLoadBalancer][ERROR] (%v)", err)
		return err
	}
	if exist {
		klog.V(5).Infof("[EnsureLoadBalancerDeleted] - LB exist:%v - call DeleteLoadBalancer", ip)
		err := lb.client.DeleteLoadBalancer(prepServiceLoadBalancer)
		if err != nil {
			klog.V(5).Infof("[EnsureLoadBalancerDeleted][DeleteLoadBalancer][ERROR] (%v)", err)
			return err
		}
	}
	return nil
}
