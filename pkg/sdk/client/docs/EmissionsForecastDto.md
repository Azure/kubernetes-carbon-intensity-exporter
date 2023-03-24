# EmissionsForecastDto

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**GeneratedAt** | [**time.Time**](time.Time.md) | Timestamp when the forecast was generated. | [optional] [default to null]
**RequestedAt** | [**time.Time**](time.Time.md) | For current requests, this value is the timestamp the request for forecast data was made.  For historical forecast requests, this value is the timestamp used to access the most   recently generated forecast as of that time. | [optional] [default to null]
**Location** | **string** | The location of the forecast | [optional] [default to null]
**DataStartAt** | [**time.Time**](time.Time.md) | Start time boundary of forecasted data points. Ignores forecast data points before this time.  Defaults to the earliest time in the forecast data. | [optional] [default to null]
**DataEndAt** | [**time.Time**](time.Time.md) | End time boundary of forecasted data points. Ignores forecast data points after this time.  Defaults to the latest time in the forecast data. | [optional] [default to null]
**WindowSize** | **int32** | The estimated duration (in minutes) of the workload.  Defaults to the duration of a single forecast data point. | [optional] [default to null]
**OptimalDataPoints** | [**[]EmissionsDataDto**](EmissionsDataDTO.md) | The optimal forecasted data point within the &#x27;forecastData&#x27; array.  Null if &#x27;forecastData&#x27; array is empty. | [optional] [default to null]
**ForecastData** | [**[]EmissionsDataDto**](EmissionsDataDTO.md) | The forecasted data points transformed and filtered to reflect the specified time and window parameters.  Points are ordered chronologically; Empty array if all data points were filtered out.  E.G. dataStartAt and dataEndAt times outside the forecast period; windowSize greater than total duration of forecast data; | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

