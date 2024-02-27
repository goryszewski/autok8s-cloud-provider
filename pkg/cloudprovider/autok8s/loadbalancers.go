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

// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (lb *loadbalancers) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	klog.V(5).Infof("GetLoadBalancer ---")
	klog.V(5).Infof("GetLoadBalancer ---")
	klog.V(5).Infof("GetLoadBalancer service(%v)", service.Namespace+service.Name)
	status1 := lb.client.UnBindLB(service.Namespace + service.Name)
	klog.V(5).Infof("GetLoadBalancer STATUS(%v)", status1)
	return nil, false, nil
}

// GetLoadBalancerName returns the name of the load balancer. Implementations must treat the
// *v1.Service parameter as read-only and not modify it.
func (lb *loadbalancers) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	klog.V(5).Infof("GetLoadBalancerName ---")
	klog.V(5).Infof("GetLoadBalancerName service(%v)", service.Namespace+service.Name)
	status := lb.client.UnBindLB(service.Namespace + service.Name)
	klog.V(5).Infof("GetLoadBalancerName STATUS(%v)", status)
	return ""
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one. Returns the status of the balancer
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (lb *loadbalancers) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {

	var LoadBalancerI []v1.LoadBalancerIngress
	lbfree, _ := lb.client.GetFreeLB()
	LoadBalancerI = append(LoadBalancerI, v1.LoadBalancerIngress{IP: lbfree.Ip})
	test := lb.client.BindLB(lbfree.Ip, service.Namespace+service.Name, "test")
	klog.V(5).Infof("EnsureLoadBalancer test (%v)", test)
	klog.V(5).Infof("EnsureLoadBalancer(%v)", clusterName)
	klog.V(5).Infof("EnsureLoadBalancer ---")
	klog.V(5).Infof("EnsureLoadBalancer service(%v)", service.Name)
	klog.V(5).Infof("EnsureLoadBalancer service(%v)", service.Namespace)
	klog.V(5).Infof("EnsureLoadBalancer service(%v)", service.Labels)
	klog.V(5).Infof("EnsureLoadBalancer ---")
	// klog.V(5).Infof("EnsureLoadBalancer nodes(%v)", nodes)
	klog.V(5).Infof("EnsureLoadBalancer ---")
	return &v1.LoadBalancerStatus{Ingress: LoadBalancerI}, nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (lb *loadbalancers) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	klog.V(5).Infof("UpdateLoadBalancer ---")
	klog.V(5).Infof("UpdateLoadBalancer(%v)", clusterName)
	klog.V(5).Infof("UpdateLoadBalancer ---")
	klog.V(5).Infof("UpdateLoadBalancer service(%v)", service.Name)

	// service(%v)&Service{ObjectMeta:{
	// nginx-service  default  9a68da3b-cacd-414a-9b50-fd724dd19291 15302 0 2024-02-05 21:30:09 +0100 CET <nil> <nil> map[] map[kubectl.kubernetes.io/last-applied-configuration:{"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"name":"nginx-service","namespace":"default"},"spec":{"ports":[{"name":"name-of-service-port","port":80,"protocol":"TCP","targetPort":"http-web-svc"}],"selector":{"app.kubernetes.io/name":"proxy"},"type":"LoadBalancer"}}
	// ] [] []  [{kubectl-client-side-apply Update v1 2024-02-05 21:30:09 +0100 CET FieldsV1 FieldsV1{Raw:*[123 34 102 58 109 101 116 97 100 97 116 97 34 58 123 34 102 58 97 110 110 111 116 97 116 105 111 110 115 34 58 123 34 46 34 58 123 125 44 34 102 58 107 117 98 101 99 116 108 46 107 117 98 101 114 110 101 116 101 115 46 105 111 47 108 97 115 116 45 97 112 112 108 105 101 100 45 99 111 110 102 105 103 117 114 97 116 105 111 110 34 58 123 125 125 125 44 34 102 58 115 112 101 99 34 58 123 34 102 58 97 108 108 111 99 97 116 101 76 111 97 100 66 97 108 97 110 99 101 114 78 111 100 101 80 111 114 116 115 34 58 123 125 44 34 102 58 101 120 116 101 114 110 97 108 84 114 97 102 102 105 99 80 111 108 105 99 121 34 58 123 125 44 34 102 58 105 110 116 101 114 110 97 108 84 114 97 102 102 105 99 80 111 108 105 99 121 34 58 123 125 44 34 102 58 112 111 114 116 115 34 58 123 34 46 34 58 123 125 44 34 107 58 123 92 34 112 111 114 116 92 34 58 56 48 44 92 34 112 114 111 116 111 99 111 108 92 34 58 92 34 84 67 80 92 34 125 34 58 123 34 46 34 58 123 125 44 34 102 58 110 97 109 101 34 58 123 125 44 34 102 58 112 111 114 116 34 58 123 125 44 34 102 58 112 114 111 116 111 99 111 108 34 58 123 125 44 34 102 58 116 97 114 103 101 116 80 111 114 116 34 58 123 125 125 125 44 34 102 58 115 101 108 101 99 116 111 114 34 58 123 125 44 34 102 58 115 101 115 115 105 111 110 65 102 102 105 110 105 116 121 34 58 123 125 44 34 102 58 116 121 112 101 34 58 123 125 125 125
	// ],}}]},
	// Spec:ServiceSpec{Ports:[]ServicePort{ServicePort{Name:name-of-service-port,Protocol:TCP,Port:80,TargetPort:{1 0 http-web-svc},NodePort:30510,AppProtocol:nil,},},Selector:map[string]string{app.kubernetes.io/name: proxy,},ClusterIP:10.32.0.227,Type:LoadBalancer,ExternalIPs:[],SessionAffinity:None,LoadBalancerIP:,LoadBalancerSourceRanges:[],ExternalName:,ExternalTrafficPolicy:Cluster,HealthCheckNodePort:0,PublishNotReadyAddresses:false,SessionAffinityConfig:nil,IPFamily:nil,TopologyKeys:[],},Status:ServiceStatus{LoadBalancer:LoadBalancerStatus{Ingress:[]LoadBalancerIngress{},},},}
	// klog.V(5).Infof("UpdateLoadBalancer ---")
	// klog.V(5).Infof("UpdateLoadBalancer nodes(%v)", nodes)
	// klog.V(5).Infof("UpdateLoadBalancer ---")
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
	klog.V(5).Infof("EnsureLoadBalancerDeleted ---")
	klog.V(5).Infof("EnsureLoadBalancerDeleted(%v)", clusterName)
	klog.V(5).Infof("EnsureLoadBalancerDeleted ---")
	klog.V(5).Infof("EnsureLoadBalancerDeleted - service(%v)", service.Name)
	status := lb.client.UnBindLB(service.Namespace + service.Name)
	klog.V(5).Infof("EnsureLoadBalancerDeleted STATUS(%v)", status)
	return nil
}
