package azureclient

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-06-01/network"
	"github.com/Azure/go-autorest/autorest"
)

// VirtualNetworksClient is a minimal interface for azure VirtualNetworkClient
type VirtualNetworksClient interface {
	Get(ctx context.Context, resourceGroupName string, virtualNetworkName string, expand string) (result network.VirtualNetwork, err error)
	List(ctx context.Context, resourceGroupName string) (result network.VirtualNetworkListResultPage, err error)
	CreateOrUpdate(ctx context.Context, resourceGroupName string, virtualNetworkName string, parameters network.VirtualNetwork) (network.VirtualNetworksCreateOrUpdateFuture, error)
	Delete(ctx context.Context, resourceGroupName string, virtualNetworkName string) (network.VirtualNetworksDeleteFuture, error)
	Client
}

type virtualNetworksClient struct {
	network.VirtualNetworksClient
}

var _ VirtualNetworksClient = &virtualNetworksClient{}

// NewVirtualNetworkClient creates a new VirtualNetworkClient
func NewVirtualNetworkClient(subscriptionID string, authorizer autorest.Authorizer, languages []string) VirtualNetworksClient {
	client := network.NewVirtualNetworksClient(subscriptionID)
	client.Authorizer = authorizer
	client.RequestInspector = addAcceptLanguages(languages)

	return &virtualNetworksClient{
		VirtualNetworksClient: client,
	}
}

func (c *virtualNetworksClient) Client() autorest.Client {
	return c.VirtualNetworksClient.Client
}

// VirtualNetworksPeeringClient is a minimal interface for azure NewVirtualNetworkPeeringsClient
type VirtualNetworksPeeringsClient interface {
	Delete(ctx context.Context, resourceGroupName string, virtualNetworkName string, instanceID string) (network.VirtualNetworkPeeringsDeleteFuture, error)
	List(ctx context.Context, resourceGroupName string, virtualNetworkName string) (network.VirtualNetworkPeeringListResultPage, error)
	Client
}

type virtualNetworkPeeringsClient struct {
	network.VirtualNetworkPeeringsClient
}

var _ VirtualNetworksPeeringsClient = &virtualNetworkPeeringsClient{}

// NewVirtualNetworksPeeringsClient creates a new VirtualMachineScaleSetVMsClient
func NewVirtualNetworksPeeringsClient(subscriptionID string, authorizer autorest.Authorizer, languages []string) VirtualNetworksPeeringsClient {
	client := network.NewVirtualNetworkPeeringsClient(subscriptionID)
	client.Authorizer = authorizer
	client.RequestInspector = addAcceptLanguages(languages)

	return &virtualNetworkPeeringsClient{
		VirtualNetworkPeeringsClient: client,
	}
}

func (c *virtualNetworkPeeringsClient) List(ctx context.Context, resourceGroupName string, virtualNetworkName string) (network.VirtualNetworkPeeringListResultPage, error) {
	nwRes, err := c.VirtualNetworkPeeringsClient.List(ctx, resourceGroupName, virtualNetworkName)
	return nwRes, err
}

func (c *virtualNetworkPeeringsClient) Client() autorest.Client {
	return c.VirtualNetworkPeeringsClient.Client
}
