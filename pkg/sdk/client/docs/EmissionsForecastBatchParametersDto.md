# EmissionsForecastBatchParametersDto

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RequestedAt** | [**time.Time**](time.Time.md) | For historical forecast requests, this value is the timestamp used to access the most  recently generated forecast as of that time. | [optional] [default to null]
**Location** | **string** | The location of the forecast | [optional] [default to null]
**DataStartAt** | [**time.Time**](time.Time.md) | Start time boundary of forecasted data points.Ignores current forecast data points before this time.  Defaults to the earliest time in the forecast data. | [optional] [default to null]
**DataEndAt** | [**time.Time**](time.Time.md) | End time boundary of forecasted data points. Ignores current forecast data points after this time.  Defaults to the latest time in the forecast data. | [optional] [default to null]
**WindowSize** | **int32** | The estimated duration (in minutes) of the workload.  Defaults to the duration of a single forecast data point. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

