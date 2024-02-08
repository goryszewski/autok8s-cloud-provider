package autok8s

import (
	"io"

	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

const ProviderName = "autok8s"

type autok8s struct {
	providerName  string
	instances     cloudprovider.Instances
	zones         cloudprovider.Zones
	loadbalancers cloudprovider.LoadBalancer
}

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, func(io.Reader) (cloudprovider.Interface, error) {
		return newCloud()
	})
}

// newCloud returns a cloudprovider.Interface
func newCloud() (cloudprovider.Interface, error) {
	// Bootstrap HTTP client here
	cc := newAutok8sClient()

	return &autok8s{
		instances:     newInstances(cc),
		zones:         newZones(cc),
		loadbalancers: newLoadBalancers(cc),
		// instancesv2:   newInstancesV2(cc),
	}, nil
}

func (c *autok8s) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
	// Start your own controllers here
	klog.V(5).Info("Initialize()")
}

func (c *autok8s) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	klog.V(5).Info("LoadBalancer()")
	return c.loadbalancers, true
}

func (c *autok8s) Instances() (cloudprovider.Instances, bool) {
	klog.V(5).Info("Instances()")
	return c.instances, true
}

// InstancesV2 is not implemented
func (c *autok8s) InstancesV2() (cloudprovider.InstancesV2, bool) {
	klog.V(5).Info("Instancesv2()")
	return nil, true
}

func (c *autok8s) Zones() (cloudprovider.Zones, bool) {
	klog.V(5).Info("Zones()")
	return c.zones, true
}

// Clusters is not implemented
func (c *autok8s) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

// Routes is not implemented
func (c *autok8s) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

// ProviderName returns this cloud providers name
func (c *autok8s) ProviderName() string {
	klog.V(5).Infof("ProviderName() returned %s", ProviderName)
	return ProviderName
}

func (c *autok8s) HasClusterID() bool {
	klog.V(5).Info("HasClusterID()")
	return true
}
