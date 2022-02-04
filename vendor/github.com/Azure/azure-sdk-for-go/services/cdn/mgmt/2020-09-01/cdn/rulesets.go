package cdn

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// RuleSetsClient is the cdn Management Client
type RuleSetsClient struct {
	BaseClient
}

// NewRuleSetsClient creates an instance of the RuleSetsClient client.
func NewRuleSetsClient(subscriptionID string) RuleSetsClient {
	return NewRuleSetsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewRuleSetsClientWithBaseURI creates an instance of the RuleSetsClient client using a custom endpoint.  Use this
// when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewRuleSetsClientWithBaseURI(baseURI string, subscriptionID string) RuleSetsClient {
	return RuleSetsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Create creates a new rule set within the specified profile.
// Parameters:
// resourceGroupName - name of the Resource group within the Azure subscription.
// profileName - name of the CDN profile which is unique within the resource group.
// ruleSetName - name of the rule set under the profile which is unique globally
func (client RuleSetsClient) Create(ctx context.Context, resourceGroupName string, profileName string, ruleSetName string) (result RuleSetsCreateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/RuleSetsClient.Create")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("cdn.RuleSetsClient", "Create", err.Error())
	}

	req, err := client.CreatePreparer(ctx, resourceGroupName, profileName, ruleSetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "Create", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "Create", result.Response(), "Failure sending request")
		return
	}

	return
}

// CreatePreparer prepares the Create request.
func (client RuleSetsClient) CreatePreparer(ctx context.Context, resourceGroupName string, profileName string, ruleSetName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"profileName":       autorest.Encode("path", profileName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"ruleSetName":       autorest.Encode("path", ruleSetName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/ruleSets/{ruleSetName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client RuleSetsClient) CreateSender(req *http.Request) (future RuleSetsCreateFuture, err error) {
	var resp *http.Response
	future.FutureAPI = &azure.Future{}
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client RuleSetsClient) CreateResponder(resp *http.Response) (result RuleSet, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes an existing AzureFrontDoor rule set with the specified rule set name under the specified
// subscription, resource group and profile.
// Parameters:
// resourceGroupName - name of the Resource group within the Azure subscription.
// profileName - name of the CDN profile which is unique within the resource group.
// ruleSetName - name of the rule set under the profile which is unique globally.
func (client RuleSetsClient) Delete(ctx context.Context, resourceGroupName string, profileName string, ruleSetName string) (result RuleSetsDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/RuleSetsClient.Delete")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("cdn.RuleSetsClient", "Delete", err.Error())
	}

	req, err := client.DeletePreparer(ctx, resourceGroupName, profileName, ruleSetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "Delete", result.Response(), "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client RuleSetsClient) DeletePreparer(ctx context.Context, resourceGroupName string, profileName string, ruleSetName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"profileName":       autorest.Encode("path", profileName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"ruleSetName":       autorest.Encode("path", ruleSetName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/ruleSets/{ruleSetName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client RuleSetsClient) DeleteSender(req *http.Request) (future RuleSetsDeleteFuture, err error) {
	var resp *http.Response
	future.FutureAPI = &azure.Future{}
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client RuleSetsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets an existing AzureFrontDoor rule set with the specified rule set name under the specified subscription,
// resource group and profile.
// Parameters:
// resourceGroupName - name of the Resource group within the Azure subscription.
// profileName - name of the CDN profile which is unique within the resource group.
// ruleSetName - name of the rule set under the profile which is unique globally.
func (client RuleSetsClient) Get(ctx context.Context, resourceGroupName string, profileName string, ruleSetName string) (result RuleSet, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/RuleSetsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("cdn.RuleSetsClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, profileName, ruleSetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client RuleSetsClient) GetPreparer(ctx context.Context, resourceGroupName string, profileName string, ruleSetName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"profileName":       autorest.Encode("path", profileName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"ruleSetName":       autorest.Encode("path", ruleSetName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/ruleSets/{ruleSetName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client RuleSetsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client RuleSetsClient) GetResponder(resp *http.Response) (result RuleSet, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByProfile lists existing AzureFrontDoor rule sets within a profile.
// Parameters:
// resourceGroupName - name of the Resource group within the Azure subscription.
// profileName - name of the CDN profile which is unique within the resource group.
func (client RuleSetsClient) ListByProfile(ctx context.Context, resourceGroupName string, profileName string) (result RuleSetListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/RuleSetsClient.ListByProfile")
		defer func() {
			sc := -1
			if result.rslr.Response.Response != nil {
				sc = result.rslr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("cdn.RuleSetsClient", "ListByProfile", err.Error())
	}

	result.fn = client.listByProfileNextResults
	req, err := client.ListByProfilePreparer(ctx, resourceGroupName, profileName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "ListByProfile", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByProfileSender(req)
	if err != nil {
		result.rslr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "ListByProfile", resp, "Failure sending request")
		return
	}

	result.rslr, err = client.ListByProfileResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "ListByProfile", resp, "Failure responding to request")
		return
	}
	if result.rslr.hasNextLink() && result.rslr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListByProfilePreparer prepares the ListByProfile request.
func (client RuleSetsClient) ListByProfilePreparer(ctx context.Context, resourceGroupName string, profileName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"profileName":       autorest.Encode("path", profileName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/ruleSets", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByProfileSender sends the ListByProfile request. The method will close the
// http.Response Body if it receives an error.
func (client RuleSetsClient) ListByProfileSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByProfileResponder handles the response to the ListByProfile request. The method always
// closes the http.Response Body.
func (client RuleSetsClient) ListByProfileResponder(resp *http.Response) (result RuleSetListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByProfileNextResults retrieves the next set of results, if any.
func (client RuleSetsClient) listByProfileNextResults(ctx context.Context, lastResults RuleSetListResult) (result RuleSetListResult, err error) {
	req, err := lastResults.ruleSetListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "listByProfileNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByProfileSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "listByProfileNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByProfileResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "listByProfileNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByProfileComplete enumerates all values, automatically crossing page boundaries as required.
func (client RuleSetsClient) ListByProfileComplete(ctx context.Context, resourceGroupName string, profileName string) (result RuleSetListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/RuleSetsClient.ListByProfile")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByProfile(ctx, resourceGroupName, profileName)
	return
}

// ListResourceUsage checks the quota and actual usage of endpoints under the given CDN profile.
// Parameters:
// resourceGroupName - name of the Resource group within the Azure subscription.
// profileName - name of the CDN profile which is unique within the resource group.
// ruleSetName - name of the rule set under the profile which is unique globally.
func (client RuleSetsClient) ListResourceUsage(ctx context.Context, resourceGroupName string, profileName string, ruleSetName string) (result UsagesListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/RuleSetsClient.ListResourceUsage")
		defer func() {
			sc := -1
			if result.ulr.Response.Response != nil {
				sc = result.ulr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("cdn.RuleSetsClient", "ListResourceUsage", err.Error())
	}

	result.fn = client.listResourceUsageNextResults
	req, err := client.ListResourceUsagePreparer(ctx, resourceGroupName, profileName, ruleSetName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "ListResourceUsage", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListResourceUsageSender(req)
	if err != nil {
		result.ulr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "ListResourceUsage", resp, "Failure sending request")
		return
	}

	result.ulr, err = client.ListResourceUsageResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "ListResourceUsage", resp, "Failure responding to request")
		return
	}
	if result.ulr.hasNextLink() && result.ulr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListResourceUsagePreparer prepares the ListResourceUsage request.
func (client RuleSetsClient) ListResourceUsagePreparer(ctx context.Context, resourceGroupName string, profileName string, ruleSetName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"profileName":       autorest.Encode("path", profileName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"ruleSetName":       autorest.Encode("path", ruleSetName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/ruleSets/{ruleSetName}/usages", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListResourceUsageSender sends the ListResourceUsage request. The method will close the
// http.Response Body if it receives an error.
func (client RuleSetsClient) ListResourceUsageSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResourceUsageResponder handles the response to the ListResourceUsage request. The method always
// closes the http.Response Body.
func (client RuleSetsClient) ListResourceUsageResponder(resp *http.Response) (result UsagesListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listResourceUsageNextResults retrieves the next set of results, if any.
func (client RuleSetsClient) listResourceUsageNextResults(ctx context.Context, lastResults UsagesListResult) (result UsagesListResult, err error) {
	req, err := lastResults.usagesListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "listResourceUsageNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListResourceUsageSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "listResourceUsageNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResourceUsageResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RuleSetsClient", "listResourceUsageNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListResourceUsageComplete enumerates all values, automatically crossing page boundaries as required.
func (client RuleSetsClient) ListResourceUsageComplete(ctx context.Context, resourceGroupName string, profileName string, ruleSetName string) (result UsagesListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/RuleSetsClient.ListResourceUsage")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListResourceUsage(ctx, resourceGroupName, profileName, ruleSetName)
	return
}
