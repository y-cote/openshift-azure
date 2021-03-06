package containerinstance

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// ContainerLogsClient is the client for the ContainerLogs methods of the Containerinstance service.
type ContainerLogsClient struct {
	BaseClient
}

// NewContainerLogsClient creates an instance of the ContainerLogsClient client.
func NewContainerLogsClient(subscriptionID string) ContainerLogsClient {
	return NewContainerLogsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewContainerLogsClientWithBaseURI creates an instance of the ContainerLogsClient client.
func NewContainerLogsClientWithBaseURI(baseURI string, subscriptionID string) ContainerLogsClient {
	return ContainerLogsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// List get the logs for a specified container instance in a specified resource group and container group.
// Parameters:
// resourceGroupName - the name of the resource group.
// containerGroupName - the name of the container group.
// containerName - the name of the container instance.
// tail - the number of lines to show from the tail of the container instance log. If not provided, all
// available logs are shown up to 4mb.
func (client ContainerLogsClient) List(ctx context.Context, resourceGroupName string, containerGroupName string, containerName string, tail *int32) (result Logs, err error) {
	req, err := client.ListPreparer(ctx, resourceGroupName, containerGroupName, containerName, tail)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerLogsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerLogsClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerLogsClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client ContainerLogsClient) ListPreparer(ctx context.Context, resourceGroupName string, containerGroupName string, containerName string, tail *int32) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"containerGroupName": autorest.Encode("path", containerGroupName),
		"containerName":      autorest.Encode("path", containerName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2017-10-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if tail != nil {
		queryParameters["tail"] = autorest.Encode("query", *tail)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerInstance/containerGroups/{containerGroupName}/containers/{containerName}/logs", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client ContainerLogsClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ContainerLogsClient) ListResponder(resp *http.Response) (result Logs, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
