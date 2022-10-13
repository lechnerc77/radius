//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package generated

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// VolumesClient contains the methods for the Volumes group.
// Don't use this type directly, use NewVolumesClient() instead.
type VolumesClient struct {
	host string
	rootScope string
	pl runtime.Pipeline
}

// NewVolumesClient creates a new instance of VolumesClient with the specified values.
// rootScope - The scope in which the resource is present. For Azure resource this would be /subscriptions/{subscriptionID}/resourceGroup/{resourcegroupID}
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewVolumesClient(rootScope string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VolumesClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &VolumesClient{
		rootScope: rootScope,
		host: ep,
pl: pl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update an Volume.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// volumeName - The name of the Volume.
// volumeResource - Volume details
// options - VolumesClientCreateOrUpdateOptions contains the optional parameters for the VolumesClient.CreateOrUpdate method.
func (client *VolumesClient) CreateOrUpdate(ctx context.Context, volumeName string, volumeResource VolumeResource, options *VolumesClientCreateOrUpdateOptions) (VolumesClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, volumeName, volumeResource, options)
	if err != nil {
		return VolumesClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VolumesClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusNoContent) {
		return VolumesClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *VolumesClient) createOrUpdateCreateRequest(ctx context.Context, volumeName string, volumeResource VolumeResource, options *VolumesClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/volumes/{volumeName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if volumeName == "" {
		return nil, errors.New("parameter volumeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{volumeName}", url.PathEscape(volumeName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, volumeResource)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *VolumesClient) createOrUpdateHandleResponse(resp *http.Response) (VolumesClientCreateOrUpdateResponse, error) {
	result := VolumesClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VolumeResource); err != nil {
		return VolumesClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete an Volume.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// volumeName - The name of the Volume.
// options - VolumesClientDeleteOptions contains the optional parameters for the VolumesClient.Delete method.
func (client *VolumesClient) Delete(ctx context.Context, volumeName string, options *VolumesClientDeleteOptions) (VolumesClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, volumeName, options)
	if err != nil {
		return VolumesClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VolumesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return VolumesClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return VolumesClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *VolumesClient) deleteCreateRequest(ctx context.Context, volumeName string, options *VolumesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/volumes/{volumeName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if volumeName == "" {
		return nil, errors.New("parameter volumeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{volumeName}", url.PathEscape(volumeName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the properties of an Volume.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// volumeName - The name of the Volume.
// options - VolumesClientGetOptions contains the optional parameters for the VolumesClient.Get method.
func (client *VolumesClient) Get(ctx context.Context, volumeName string, options *VolumesClientGetOptions) (VolumesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, volumeName, options)
	if err != nil {
		return VolumesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VolumesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VolumesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *VolumesClient) getCreateRequest(ctx context.Context, volumeName string, options *VolumesClientGetOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/volumes/{volumeName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if volumeName == "" {
		return nil, errors.New("parameter volumeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{volumeName}", url.PathEscape(volumeName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VolumesClient) getHandleResponse(resp *http.Response) (VolumesClientGetResponse, error) {
	result := VolumesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VolumeResource); err != nil {
		return VolumesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByScopePager - List all volumes in the given scope.
// Generated from API version 2022-03-15-privatepreview
// options - VolumesClientListByScopeOptions contains the optional parameters for the VolumesClient.ListByScope method.
func (client *VolumesClient) NewListByScopePager(options *VolumesClientListByScopeOptions) (*runtime.Pager[VolumesClientListByScopeResponse]) {
	return runtime.NewPager(runtime.PagingHandler[VolumesClientListByScopeResponse]{
		More: func(page VolumesClientListByScopeResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VolumesClientListByScopeResponse) (VolumesClientListByScopeResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByScopeCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VolumesClientListByScopeResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return VolumesClientListByScopeResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VolumesClientListByScopeResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByScopeHandleResponse(resp)
		},
	})
}

// listByScopeCreateRequest creates the ListByScope request.
func (client *VolumesClient) listByScopeCreateRequest(ctx context.Context, options *VolumesClientListByScopeOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/volumes"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByScopeHandleResponse handles the ListByScope response.
func (client *VolumesClient) listByScopeHandleResponse(resp *http.Response) (VolumesClientListByScopeResponse, error) {
	result := VolumesClientListByScopeResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VolumeResourceList); err != nil {
		return VolumesClientListByScopeResponse{}, err
	}
	return result, nil
}

// Update - Update the properties of an existing Volume.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-15-privatepreview
// volumeName - The name of the Volume.
// volumeResource - Volume details
// options - VolumesClientUpdateOptions contains the optional parameters for the VolumesClient.Update method.
func (client *VolumesClient) Update(ctx context.Context, volumeName string, volumeResource VolumeResource, options *VolumesClientUpdateOptions) (VolumesClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, volumeName, volumeResource, options)
	if err != nil {
		return VolumesClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VolumesClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusNoContent) {
		return VolumesClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *VolumesClient) updateCreateRequest(ctx context.Context, volumeName string, volumeResource VolumeResource, options *VolumesClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/volumes/{volumeName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if volumeName == "" {
		return nil, errors.New("parameter volumeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{volumeName}", url.PathEscape(volumeName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, volumeResource)
}

// updateHandleResponse handles the Update response.
func (client *VolumesClient) updateHandleResponse(resp *http.Response) (VolumesClientUpdateResponse, error) {
	result := VolumesClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VolumeResource); err != nil {
		return VolumesClientUpdateResponse{}, err
	}
	return result, nil
}

