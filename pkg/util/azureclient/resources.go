package azureclient

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
)

// DeploymentsClient is a minimal interface for azure DeploymentsClient
type DeploymentsClient interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, deploymentName string, parameters resources.Deployment) (result resources.DeploymentsCreateOrUpdateFuture, err error)
	Client
	DeploymentClient() resources.DeploymentsClient
}

type deploymentsClient struct {
	resources.DeploymentsClient
}

var _ DeploymentsClient = &deploymentsClient{}

// NewDeploymentsClient creates a new DeploymentsClient
func NewDeploymentsClient(subscriptionID string, authorizer autorest.Authorizer, languages []string) DeploymentsClient {
	client := resources.NewDeploymentsClient(subscriptionID)
	setupClient(&client.Client, authorizer, languages)

	return &deploymentsClient{
		DeploymentsClient: client,
	}
}

func (c *deploymentsClient) Client() autorest.Client {
	return c.DeploymentsClient.Client
}

func (c *deploymentsClient) DeploymentClient() resources.DeploymentsClient {
	return c.DeploymentsClient
}

// ResourcesClient is a minimal interface for azure Resources Client
type ResourcesClient interface {
	DeleteByID(ctx context.Context, resourceID string) (result resources.DeleteByIDFuture, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string, filter string, expand string, top *int32) (result resources.ListResultPage, err error)
}

type resourcesClient struct {
	resources.Client
}

var _ ResourcesClient = &resourcesClient{}

// NewResourcesClient creates a new ResourcesClient
func NewResourcesClient(subscriptionID string, authorizer autorest.Authorizer, languages []string) ResourcesClient {
	client := resources.NewClient(subscriptionID)
	setupClient(&client.Client, authorizer, languages)

	return &resourcesClient{
		Client: client,
	}
}

// GroupsClient is a minimal interface for azure Resources Client
type GroupsClient interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, parameters resources.Group) (result resources.Group, err error)
	CheckExistence(ctx context.Context, resourceGroupName string) (exists bool, err error)
	List(ctx context.Context, filter string, top *int32) (result resources.GroupListResultPage, err error)
	Delete(ctx context.Context, resourceGroupName string) (result resources.GroupsDeleteFuture, err error)
	Client
}

type groupsClient struct {
	resources.GroupsClient
}

var _ GroupsClient = &groupsClient{}

// NewGroupsClient creates a new ResourcesClient
func NewGroupsClient(subscriptionID string, authorizer autorest.Authorizer, languages []string) GroupsClient {
	client := resources.NewGroupsClient(subscriptionID)
	client.Authorizer = authorizer
	client.RequestInspector = addAcceptLanguages(languages)

	return &groupsClient{
		GroupsClient: client,
	}
}

func (c *groupsClient) Client() autorest.Client {
	return c.GroupsClient.Client
}

func (c *groupsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, parameters resources.Group) (resources.Group, error) {
	return c.GroupsClient.CreateOrUpdate(ctx, resourceGroupName, parameters)
}

func (c *groupsClient) CheckExistence(ctx context.Context, resourceGroupName string) (bool, error) {
	resp, err := c.GroupsClient.CheckExistence(ctx, resourceGroupName)
	switch resp.StatusCode {
	case 404:
		return false, nil
	case 200:
		return true, nil
	default:
		return false, err
	}
}
