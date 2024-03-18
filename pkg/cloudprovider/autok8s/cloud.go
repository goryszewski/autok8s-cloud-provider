package autok8s

import (
	"fmt"
	"io"
	"net/http"
	"time"

	libvirtApiClient "github.com/goryszewski/libvirtApi-client/libvirtApiClient"
	gcfg "gopkg.in/gcfg.v1"
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
type Config struct {
	Global libvirtApiClient.Config
}

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, func(config io.Reader) (cloudprovider.Interface, error) {
		cfg, err := ReadConfig(config)
		if err != nil {
			klog.Warningf("failed to read config: %v", err)
			return nil, err
		}
		return newCloud(cfg)
	})
}

func ReadConfig(config io.Reader) (libvirtApiClient.Config, error) {
	if config == nil {
		return libvirtApiClient.Config{}, fmt.Errorf("no autok8s cloud provider config file given")
	}
	var cfg Config

	err := gcfg.FatalOnly(gcfg.ReadInto(&cfg, config))
	if err != nil {

		return libvirtApiClient.Config{}, err
	}
	klog.Infof("Config URL: %v and User: %v", *cfg.Global.Url, *cfg.Global.Username)
	return cfg.Global, nil
}

// newCloud returns a cloudprovider.Interface
func newCloud(conf libvirtApiClient.Config) (cloudprovider.Interface, error) {
	// Bootstrap HTTP client here
	// cc := newAutok8sClient()

	cc, err := libvirtApiClient.NewClient(conf, &http.Client{Timeout: 10 * time.Second})

	if err != nil {
		return &autok8s{}, err
	}

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
