//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20220315privatepreview

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// RabbitMQMessageQueuesClient contains the methods for the RabbitMQMessageQueues group.
// Don't use this type directly, use NewRabbitMQMessageQueuesClient() instead.
type RabbitMQMessageQueuesClient struct {
	ep string
	pl runtime.Pipeline
	rootScope string
	subscriptionID string
}

// NewRabbitMQMessageQueuesClient creates a new instance of RabbitMQMessageQueuesClient with the specified values.
func NewRabbitMQMessageQueuesClient(con *arm.Connection, rootScope string, subscriptionID string) *RabbitMQMessageQueuesClient {
	return &RabbitMQMessageQueuesClient{ep: con.Endpoint(), pl: con.NewPipeline(module, version), rootScope: rootScope, subscriptionID: subscriptionID}
}

// CreateOrUpdate - Creates or updates a RabbitMQMessageQueue resource
// If the operation fails it returns the *ErrorResponse error type.
func (client *RabbitMQMessageQueuesClient) CreateOrUpdate(ctx context.Context, rabbitMQMessageQueueName string, rabbitMQMessageQueueParameters RabbitMQMessageQueueResource, options *RabbitMQMessageQueuesCreateOrUpdateOptions) (RabbitMQMessageQueuesCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, rabbitMQMessageQueueName, rabbitMQMessageQueueParameters, options)
	if err != nil {
		return RabbitMQMessageQueuesCreateOrUpdateResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return RabbitMQMessageQueuesCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return RabbitMQMessageQueuesCreateOrUpdateResponse{}, client.createOrUpdateHandleError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *RabbitMQMessageQueuesClient) createOrUpdateCreateRequest(ctx context.Context, rabbitMQMessageQueueName string, rabbitMQMessageQueueParameters RabbitMQMessageQueueResource, options *RabbitMQMessageQueuesCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Connector/rabbitMQMessageQueues/{rabbitMQMessageQueueName}"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", url.PathEscape(client.rootScope))
	if rabbitMQMessageQueueName == "" {
		return nil, errors.New("parameter rabbitMQMessageQueueName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rabbitMQMessageQueueName}", url.PathEscape(rabbitMQMessageQueueName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, rabbitMQMessageQueueParameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *RabbitMQMessageQueuesClient) createOrUpdateHandleResponse(resp *http.Response) (RabbitMQMessageQueuesCreateOrUpdateResponse, error) {
	result := RabbitMQMessageQueuesCreateOrUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.RabbitMQMessageQueueResource); err != nil {
		return RabbitMQMessageQueuesCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *RabbitMQMessageQueuesClient) createOrUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Delete - Deletes an existing rabbitMQMessageQueue resource
// If the operation fails it returns the *ErrorResponse error type.
func (client *RabbitMQMessageQueuesClient) Delete(ctx context.Context, rabbitMQMessageQueueName string, options *RabbitMQMessageQueuesDeleteOptions) (RabbitMQMessageQueuesDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, rabbitMQMessageQueueName, options)
	if err != nil {
		return RabbitMQMessageQueuesDeleteResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return RabbitMQMessageQueuesDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return RabbitMQMessageQueuesDeleteResponse{}, client.deleteHandleError(resp)
	}
	return RabbitMQMessageQueuesDeleteResponse{RawResponse: resp}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *RabbitMQMessageQueuesClient) deleteCreateRequest(ctx context.Context, rabbitMQMessageQueueName string, options *RabbitMQMessageQueuesDeleteOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Connector/rabbitMQMessageQueues/{rabbitMQMessageQueueName}"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", url.PathEscape(client.rootScope))
	if rabbitMQMessageQueueName == "" {
		return nil, errors.New("parameter rabbitMQMessageQueueName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rabbitMQMessageQueueName}", url.PathEscape(rabbitMQMessageQueueName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *RabbitMQMessageQueuesClient) deleteHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Retrieves information about a rabbitMQMessageQueue resource
// If the operation fails it returns the *ErrorResponse error type.
func (client *RabbitMQMessageQueuesClient) Get(ctx context.Context, rabbitMQMessageQueueName string, options *RabbitMQMessageQueuesGetOptions) (RabbitMQMessageQueuesGetResponse, error) {
	req, err := client.getCreateRequest(ctx, rabbitMQMessageQueueName, options)
	if err != nil {
		return RabbitMQMessageQueuesGetResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return RabbitMQMessageQueuesGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return RabbitMQMessageQueuesGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *RabbitMQMessageQueuesClient) getCreateRequest(ctx context.Context, rabbitMQMessageQueueName string, options *RabbitMQMessageQueuesGetOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Connector/rabbitMQMessageQueues/{rabbitMQMessageQueueName}"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", url.PathEscape(client.rootScope))
	if rabbitMQMessageQueueName == "" {
		return nil, errors.New("parameter rabbitMQMessageQueueName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rabbitMQMessageQueueName}", url.PathEscape(rabbitMQMessageQueueName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *RabbitMQMessageQueuesClient) getHandleResponse(resp *http.Response) (RabbitMQMessageQueuesGetResponse, error) {
	result := RabbitMQMessageQueuesGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.RabbitMQMessageQueueResource); err != nil {
		return RabbitMQMessageQueuesGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *RabbitMQMessageQueuesClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// ListByRootScope - Lists information about all rabbitMQMessageQueue resources in the given root scope
// If the operation fails it returns the *ErrorResponse error type.
func (client *RabbitMQMessageQueuesClient) ListByRootScope(options *RabbitMQMessageQueuesListByRootScopeOptions) (*RabbitMQMessageQueuesListByRootScopePager) {
	return &RabbitMQMessageQueuesListByRootScopePager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listByRootScopeCreateRequest(ctx, options)
		},
		advancer: func(ctx context.Context, resp RabbitMQMessageQueuesListByRootScopeResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.RabbitMQMessageQueueList.NextLink)
		},
	}
}

// listByRootScopeCreateRequest creates the ListByRootScope request.
func (client *RabbitMQMessageQueuesClient) listByRootScopeCreateRequest(ctx context.Context, options *RabbitMQMessageQueuesListByRootScopeOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Connector/rabbitMQMessageQueues"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", url.PathEscape(client.rootScope))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByRootScopeHandleResponse handles the ListByRootScope response.
func (client *RabbitMQMessageQueuesClient) listByRootScopeHandleResponse(resp *http.Response) (RabbitMQMessageQueuesListByRootScopeResponse, error) {
	result := RabbitMQMessageQueuesListByRootScopeResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.RabbitMQMessageQueueList); err != nil {
		return RabbitMQMessageQueuesListByRootScopeResponse{}, err
	}
	return result, nil
}

// listByRootScopeHandleError handles the ListByRootScope error response.
func (client *RabbitMQMessageQueuesClient) listByRootScopeHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// ListSecrets - Lists secrets values for the specified RabbitMQMessageQueue resource
// If the operation fails it returns the *ErrorResponse error type.
func (client *RabbitMQMessageQueuesClient) ListSecrets(ctx context.Context, rabbitMQMessageQueueName string, options *RabbitMQMessageQueuesListSecretsOptions) (RabbitMQMessageQueuesListSecretsResponse, error) {
	req, err := client.listSecretsCreateRequest(ctx, rabbitMQMessageQueueName, options)
	if err != nil {
		return RabbitMQMessageQueuesListSecretsResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return RabbitMQMessageQueuesListSecretsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return RabbitMQMessageQueuesListSecretsResponse{}, client.listSecretsHandleError(resp)
	}
	return client.listSecretsHandleResponse(resp)
}

// listSecretsCreateRequest creates the ListSecrets request.
func (client *RabbitMQMessageQueuesClient) listSecretsCreateRequest(ctx context.Context, rabbitMQMessageQueueName string, options *RabbitMQMessageQueuesListSecretsOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Connector/rabbitMQMessageQueues/{rabbitMQMessageQueueName}/listSecrets"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", url.PathEscape(client.rootScope))
	if rabbitMQMessageQueueName == "" {
		return nil, errors.New("parameter rabbitMQMessageQueueName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rabbitMQMessageQueueName}", url.PathEscape(rabbitMQMessageQueueName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listSecretsHandleResponse handles the ListSecrets response.
func (client *RabbitMQMessageQueuesClient) listSecretsHandleResponse(resp *http.Response) (RabbitMQMessageQueuesListSecretsResponse, error) {
	result := RabbitMQMessageQueuesListSecretsResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.RabbitMQSecrets); err != nil {
		return RabbitMQMessageQueuesListSecretsResponse{}, err
	}
	return result, nil
}

// listSecretsHandleError handles the ListSecrets error response.
func (client *RabbitMQMessageQueuesClient) listSecretsHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

