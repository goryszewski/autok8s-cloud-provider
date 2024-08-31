package autok8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/goryszewski/libvirtApi-client/libvirtApiClient"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

type instancesv2 struct {
	client *libvirtApiClient.Client
}

func GetIP(node *libvirtApiClient.NodeV2) []v1.NodeAddress {
	klog.V(5).Infof("[V2] - GetIP_NodeAddresses(%v) , IP: (%#+v)", node.Name, node.Interface)
	var addrs []v1.NodeAddress = prepAddress(node)
	return addrs
}

func newInstancesV2(c *libvirtApiClient.Client) cloudprovider.InstancesV2 {
	klog.V(5).Infof("call newInstancesV2")
	return &instancesv2{
		client: c,
	}
}
func (i *instancesv2) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	klog.V(5).Infof("[V2] - InstanceShutdown(%v)", node.Name)
	return false, nil
}
func (i *instancesv2) InstanceMetadata(ctx context.Context, node *v1.Node) (*cloudprovider.InstanceMetadata, error) {
	klog.V(5).Infof("[V2] - InstanceMetadata(%v)", node.Name)

	node2, err := i.client.GetNodeByName(node.Name)
	if err != nil {
		return nil, err
	}
	return &cloudprovider.InstanceMetadata{
		ProviderID:    fmt.Sprintf("%s://%s", ProviderName, node2.Name),
		InstanceType:  node2.Type,
		NodeAddresses: GetIP(node2),
		Region:        "PL1",
		Zone:          "A",
	}, nil
}

// NodeAddresses returns the addresses of the specified instance.
func (i *instancesv2) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	klog.V(5).Infof("call NodeAddresses")
	node, err := i.client.GetNodeByName(string(name))
	if err != nil {
		return nil, err
	}
	klog.V(5).Infof("[V2] - NodeAddresses(%v) Data:(%#+v)", name, node)
	var addrs []v1.NodeAddress = GetIP(node)

	return addrs, nil
}
func (i *instancesv2) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	klog.V(5).Infof("[V2] - InstanceExists(%v)", node.Name)
	_, err := i.client.GetNodeByName(node.Name)
	if err != nil {
		return false, err
	}
	return true, nil
}

// NodeAddressesByProviderID returns the addresses of the specified instance.
// The instance is specified using the providerID of the node. The
// ProviderID is a unique identifier of the node. This will not be called
// from the node whose nodeaddresses are being queried. i.e. local metadata
// services cannot be used in this method to obtain nodeaddresses
func (i *instancesv2) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	klog.V(5).Infof("[V2] - NodeAddressesByProviderID(%v)", providerID)

	// if providerID == "autok8s://worker01" {

	name := strings.Split(providerID, "//")[1]
	node, err := i.client.GetNodeByName(string(name))
	if err != nil {
		return nil, err
	}
	var addrs []v1.NodeAddress = GetIP(node)

	return addrs, nil
}

// InstanceID returns the cloud provider ID of the node with the specified NodeName.
// Note that if the instance does not exist, we must return ("", cloudprovider.InstanceNotFound)
// cloudprovider.InstanceNotFound should NOT be returned for instances that exist but are stopped/sleeping
func (i *instancesv2) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	klog.V(5).Infof("[V2] - InstanceID(%v)", nodeName)

	node, err := i.client.GetNodeByName(string(nodeName))
	if err != nil {
		return "", err
	}

	klog.V(5).Infof("[V2] - InstanceID(%v) Data:(%v)", string(nodeName), node)
	instanceID := "autok8s://" + fmt.Sprintf("%v", node.Name)

	return instanceID, nil
}

// InstanceType returns the type of the specified instance.
func (i *instancesv2) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	klog.V(5).Infof("[V2] - InstanceType(%v)", name)

	node, err := i.client.GetNodeByName(string(name))
	if err != nil {
		return "", err
	}

	klog.V(5).Infof("[V2] - InstanceType(%v) Data:(%v)", string(name), node)
	instanceType := node.Type

	return instanceType, nil
}

// InstanceTypeByProviderID returns the type of the specified instance.
func (i *instancesv2) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	klog.V(5).Infof("[V2] - InstanceTypeByProviderID(%v)", providerID)

	name := strings.Split(providerID, "//")[1]
	node, err := i.client.GetNodeByName(string(name))
	if err != nil {
		return "", err
	}

	instanceType := node.Type

	return instanceType, nil
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (i *instancesv2) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	klog.V(5).Info("AddSSHKeyToAllInstances(%v, %v)", user, keyData)

	nodes, err := i.client.GetNodes()
	if err != nil {
		klog.Errorf("Failed to retrieve nodes: %v", err)
		return err
	}
	for _, node := range *nodes {
		err := i.client.addSSHKeyToInstance(node, user, keyData)
		if err != nil {
			klog.Errorf("Failed to add SSH key to instance %v: %v", node.Name, err)
			return err
		}
	}
	klog.Infof("Successfully added SSH key to all instances for user %v", user)
	return nil
}

// CurrentNodeName returns the name of the node we are currently running on
// On most clouds (e.g. GCE) this is the hostname, so we provide the hostname
func (i *instancesv2) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	klog.V(5).Infof("[V2] - CurrentNodeName(%v)", hostname)

	node, err := i.client.GetNodeByName(string(hostname))
	if err != nil {
		return "", err
	}

	return types.NodeName(node.Name), nil
}

// InstanceExistsByProviderID returns true if the instance for the given provider exists.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
// This method should still return true for instances that exist but are stopped/sleeping.
func (i *instancesv2) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.V(5).Infof("[V2] - InstanceExistsByProviderID(%v)", providerID)

	name := strings.Split(providerID, "//")[1]
	node, err := i.client.GetNodeByName(string(name))
	if err != nil {
		return false, err
	}
	klog.V(5).Infof("[V2] - InstanceExistsByProviderID(%v):exists", node)
	return true, nil
}

// InstanceShutdownByProviderID returns true if the instance is shutdown in cloudprovider
func (i *instancesv2) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.V(5).Infof("[V2] - InstanceShutdownByProviderID(%v)", providerID)

	name := strings.Split(providerID, "//")[1]
	_, err := i.client.GetNodeByName(string(name))
	if err != nil {
		return false, err
	}

	return true, nil
}
