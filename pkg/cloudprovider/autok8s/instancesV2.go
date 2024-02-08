package autok8s

import (
	"context"
	"fmt"
	"net/http"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"
)

type instancesv2 struct {
	client *http.Client
}

func newInstancesV2(c *http.Client) cloudprovider.InstancesV2 {
	return &instancesv2{
		client: c,
	}
}
func (i *instancesv2) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	return false, nil
}
func (i *instancesv2) InstanceMetadata(ctx context.Context, node *v1.Node) (*cloudprovider.InstanceMetadata, error) {

	return &cloudprovider.InstanceMetadata{}, nil
}

// NodeAddresses returns the addresses of the specified instance.
func (i *instancesv2) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	klog.V(5).Infof("NodeAddresses(%v)", name)
	node, _ := ReturnJson(string(name))
	klog.V(5).Infof("NodeAddresses(%v) Data:(%v)", name, node)
	var addrs []v1.NodeAddress

	klog.V(5).Infof("NodeAddresses(%v) , Internal ip: (%v)", name, node.IP.Private)
	klog.V(5).Infof("NodeAddresses(%v) , External ip: (%v)", name, node.IP.Public)

	nodeAddr := v1.NodeAddress{
		Type:    v1.NodeInternalIP,
		Address: node.IP.Private,
	}
	nodeExternalAddr := v1.NodeAddress{
		Type:    v1.NodeExternalIP,
		Address: node.IP.Public,
	}
	nodeHostName := v1.NodeAddress{
		Type:    v1.NodeHostName,
		Address: fmt.Sprintf("%v", node.Name),
	}
	addrs = append(addrs, nodeAddr)
	addrs = append(addrs, nodeExternalAddr)
	addrs = append(addrs, nodeHostName)

	return addrs, nil
}
func (i *instancesv2) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {

	return true, nil
}

// NodeAddressesByProviderID returns the addresses of the specified instance.
// The instance is specified using the providerID of the node. The
// ProviderID is a unique identifier of the node. This will not be called
// from the node whose nodeaddresses are being queried. i.e. local metadata
// services cannot be used in this method to obtain nodeaddresses
func (i *instancesv2) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	klog.V(5).Infof("NodeAddressesByProviderID(%v)", providerID)

	// if providerID == "autok8s://worker01" {

	node, _ := ReturnJson_by_provider(providerID)
	klog.V(5).Infof("NodeAddressesByProviderID(%v) Data:(%v)", providerID, node)
	var addrs []v1.NodeAddress
	klog.V(5).Infof("NodeAddressesByProviderID(%v) , Internal ip: (%v)", providerID, node.IP.Private)
	klog.V(5).Infof("NodeAddressesByProviderID(%v) , External ip: (%v)", providerID, node.IP.Public)
	nodeAddr := v1.NodeAddress{
		Type:    v1.NodeInternalIP,
		Address: node.IP.Private,
	}
	nodeExternalAddr := v1.NodeAddress{
		Type:    v1.NodeExternalIP,
		Address: node.IP.Public,
	}
	nodeHostName := v1.NodeAddress{
		Type:    v1.NodeHostName,
		Address: fmt.Sprintf("%v", node.Name),
	}
	addrs = append(addrs, nodeAddr)
	addrs = append(addrs, nodeExternalAddr)
	addrs = append(addrs, nodeHostName)

	return addrs, nil
}

// InstanceID returns the cloud provider ID of the node with the specified NodeName.
// Note that if the instance does not exist, we must return ("", cloudprovider.InstanceNotFound)
// cloudprovider.InstanceNotFound should NOT be returned for instances that exist but are stopped/sleeping
func (i *instancesv2) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	klog.V(5).Infof("InstanceID(%v)", nodeName)

	node, _ := ReturnJson(string(nodeName))
	klog.V(5).Infof("InstanceID(%v) Data:(%v)", string(nodeName), node)
	instanceID := "autok8s://" + fmt.Sprintf("%v", node.Name)

	return instanceID, nil
}

// InstanceType returns the type of the specified instance.
func (i *instancesv2) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	klog.V(5).Infof("InstanceType(%v)", name)

	node, _ := ReturnJson(string(name))
	klog.V(5).Infof("InstanceType(%v) Data:(%v)", string(name), node)
	instanceType := node.Type

	return instanceType, nil
}

// InstanceTypeByProviderID returns the type of the specified instance.
func (i *instancesv2) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	klog.V(5).Infof("InstanceTypeByProviderID(%v)", providerID)

	node, _ := ReturnJson_by_provider(providerID)
	instanceType := node.Type

	return instanceType, nil
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (i *instancesv2) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	klog.V(5).Info("AddSSHKeyToAllInstances(%v, %v)", user, keyData)
	return cloudprovider.NotImplemented
}

// CurrentNodeName returns the name of the node we are currently running on
// On most clouds (e.g. GCE) this is the hostname, so we provide the hostname
func (i *instancesv2) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	klog.V(5).Infof("CurrentNodeName(%v)", hostname)

	node, _ := ReturnJson(hostname)

	return types.NodeName(node.Name), nil
}

// InstanceExistsByProviderID returns true if the instance for the given provider exists.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
// This method should still return true for instances that exist but are stopped/sleeping.
func (i *instancesv2) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.V(5).Infof("InstanceExistsByProviderID(%v)", providerID)

	_, exists := ReturnJson_by_provider(providerID)

	return exists, nil
}

// InstanceShutdownByProviderID returns true if the instance is shutdown in cloudprovider
func (i *instancesv2) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.V(5).Infof("InstanceShutdownByProviderID(%v)", providerID)

	// node,_ :=ReturnJson(providerID)

	return false, nil
}