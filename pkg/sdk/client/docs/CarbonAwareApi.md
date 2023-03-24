# {{classname}}

All URIs are relative to *https://virtserver.swaggerhub.com/Microsoft-hela/carbonaware/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BatchForecastDataAsync**](CarbonAwareApi.md#BatchForecastDataAsync) | **Post** /emissions/forecasts/batch | Given an array of historical forecasts, retrieves the data that contains  forecasts metadata, the optimal forecast and a range of forecasts filtered by the attributes [start...end] if provided.
[**GetAverageCarbonIntensity**](CarbonAwareApi.md#GetAverageCarbonIntensity) | **Get** /emissions/average-carbon-intensity | Retrieves the measured carbon intensity data between the time boundaries and calculates the average carbon intensity during that period.
[**GetAverageCarbonIntensityBatch**](CarbonAwareApi.md#GetAverageCarbonIntensityBatch) | **Post** /emissions/average-carbon-intensity/batch | Given an array of request objects, each with their own location and time boundaries, calculate the average carbon intensity for that location and time period   and return an array of carbon intensity objects.
[**GetBestEmissionsDataForLocationsByTime**](CarbonAwareApi.md#GetBestEmissionsDataForLocationsByTime) | **Get** /emissions/bylocations/best | Calculate the best emission data by list of locations for a specified time period.
[**GetCurrentForecastData**](CarbonAwareApi.md#GetCurrentForecastData) | **Get** /emissions/forecasts/current | Retrieves the most recent forecasted data and calculates the optimal marginal carbon intensity window.
[**GetEmissionsDataForLocationByTime**](CarbonAwareApi.md#GetEmissionsDataForLocationByTime) | **Get** /emissions/bylocation | Calculate the best emission data by location for a specified time period.
[**GetEmissionsDataForLocationsByTime**](CarbonAwareApi.md#GetEmissionsDataForLocationsByTime) | **Get** /emissions/bylocations | Calculate the observed emission data by list of locations for a specified time period.

# **BatchForecastDataAsync**
> []EmissionsForecastDto BatchForecastDataAsync(ctx, optional)
Given an array of historical forecasts, retrieves the data that contains  forecasts metadata, the optimal forecast and a range of forecasts filtered by the attributes [start...end] if provided.

This endpoint takes a batch of requests for historical forecast data, fetches them, and calculates the optimal   marginal carbon intensity windows for each using the same parameters available to the '/emissions/forecasts/current'  endpoint.                This endpoint is useful for back-testing what one might have done in the past, if they had access to the   current forecast at the time.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***CarbonAwareApiBatchForecastDataAsyncOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CarbonAwareApiBatchForecastDataAsyncOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of []EmissionsForecastBatchParametersDto**](EmissionsForecastBatchParametersDTO.md)| Array of requested forecasts. | 

### Return type

[**[]EmissionsForecastDto**](EmissionsForecastDTO.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, text/json, application/_*+json
 - **Accept**: application/json, application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAverageCarbonIntensity**
> CarbonIntensityDto GetAverageCarbonIntensity(ctx, location, startTime, endTime)
Retrieves the measured carbon intensity data between the time boundaries and calculates the average carbon intensity during that period.

This endpoint is useful for reporting the measured carbon intensity for a specific time period in a specific location.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **location** | **string**| The location name where workflow is run | 
  **startTime** | **time.Time**| The time at which the workflow we are measuring carbon intensity for started | 
  **endTime** | **time.Time**| The time at which the workflow we are measuring carbon intensity for ended | 

### Return type

[**CarbonIntensityDto**](CarbonIntensityDTO.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAverageCarbonIntensityBatch**
> []CarbonIntensityDto GetAverageCarbonIntensityBatch(ctx, optional)
Given an array of request objects, each with their own location and time boundaries, calculate the average carbon intensity for that location and time period   and return an array of carbon intensity objects.

The application only supports batching across a single location with different time boundaries. If multiple locations are provided, an error is returned.  For each item in the request array, the application returns a corresponding object containing the location, time boundaries, and average marginal carbon intensity.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***CarbonAwareApiGetAverageCarbonIntensityBatchOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CarbonAwareApiGetAverageCarbonIntensityBatchOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of []CarbonIntensityBatchParametersDto**](CarbonIntensityBatchParametersDTO.md)| Array of inputs where each contains a &quot;location&quot;, &quot;startDate&quot;, and &quot;endDate&quot; for which to calculate average marginal carbon intensity. | 

### Return type

[**[]CarbonIntensityDto**](CarbonIntensityDTO.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, text/json, application/_*+json
 - **Accept**: application/json, application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetBestEmissionsDataForLocationsByTime**
> []EmissionsData GetBestEmissionsDataForLocationsByTime(ctx, location, optional)
Calculate the best emission data by list of locations for a specified time period.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **location** | [**[]string**](string.md)| String array of named locations | 
 **optional** | ***CarbonAwareApiGetBestEmissionsDataForLocationsByTimeOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CarbonAwareApiGetBestEmissionsDataForLocationsByTimeOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **time** | **optional.Time**| [Optional] Start time for the data query. | 
 **toTime** | **optional.Time**| [Optional] End time for the data query. | 

### Return type

[**[]EmissionsData**](EmissionsData.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCurrentForecastData**
> []EmissionsForecastDto GetCurrentForecastData(ctx, location, optional)
Retrieves the most recent forecasted data and calculates the optimal marginal carbon intensity window.

This endpoint fetches only the most recently generated forecast for all provided locations.  It uses the \"dataStartAt\" and   \"dataEndAt\" parameters to scope the forecasted data points (if available for those times). If no start or end time   boundaries are provided, the entire forecast dataset is used. The scoped data points are used to calculate average marginal   carbon intensities of the specified \"windowSize\" and the optimal marginal carbon intensity window is identified.                The forecast data represents what the data source predicts future marginal carbon intesity values to be,   not actual measured emissions data (as future values cannot be known).                This endpoint is useful for determining if there is a more carbon-optimal time to use electicity predicted in the future.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **location** | [**[]string**](string.md)| String array of named locations | 
 **optional** | ***CarbonAwareApiGetCurrentForecastDataOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CarbonAwareApiGetCurrentForecastDataOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **dataStartAt** | **optional.Time**| Start time boundary of forecasted data points.Ignores current forecast data points before this time.  Defaults to the earliest time in the forecast data. | 
 **dataEndAt** | **optional.Time**| End time boundary of forecasted data points. Ignores current forecast data points after this time.  Defaults to the latest time in the forecast data. | 
 **windowSize** | **optional.Int32**| The estimated duration (in minutes) of the workload.  Defaults to the duration of a single forecast data point. | 

### Return type

[**[]EmissionsForecastDto**](EmissionsForecastDTO.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetEmissionsDataForLocationByTime**
> []EmissionsData GetEmissionsDataForLocationByTime(ctx, location, optional)
Calculate the best emission data by location for a specified time period.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **location** | **string**| String named location. | 
 **optional** | ***CarbonAwareApiGetEmissionsDataForLocationByTimeOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CarbonAwareApiGetEmissionsDataForLocationByTimeOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **startTime** | **optional.Time**| [Optional] Start time for the data query. | 
 **endTime** | **optional.Time**| [Optional] End time for the data query. | 

### Return type

[**[]EmissionsData**](EmissionsData.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetEmissionsDataForLocationsByTime**
> []EmissionsData GetEmissionsDataForLocationsByTime(ctx, location, optional)
Calculate the observed emission data by list of locations for a specified time period.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **location** | [**[]string**](string.md)| String array of named locations | 
 **optional** | ***CarbonAwareApiGetEmissionsDataForLocationsByTimeOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CarbonAwareApiGetEmissionsDataForLocationsByTimeOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **time** | **optional.Time**| [Optional] Start time for the data query. | 
 **toTime** | **optional.Time**| [Optional] End time for the data query. | 

### Return type

[**[]EmissionsData**](EmissionsData.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

