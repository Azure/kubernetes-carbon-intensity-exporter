/*
 * CarbonAware.WebApi, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package client

import (
	"time"
)

type EmissionsForecastDto struct {
	// Timestamp when the forecast was generated.
	GeneratedAt time.Time `json:"generatedAt,omitempty"`
	// For current requests, this value is the timestamp the request for forecast data was made.  For historical forecast requests, this value is the timestamp used to access the most   recently generated forecast as of that time.
	RequestedAt time.Time `json:"requestedAt,omitempty"`
	// The location of the forecast
	Location string `json:"location,omitempty"`
	// Start time boundary of forecasted data points. Ignores forecast data points before this time.  Defaults to the earliest time in the forecast data.
	DataStartAt time.Time `json:"dataStartAt,omitempty"`
	// End time boundary of forecasted data points. Ignores forecast data points after this time.  Defaults to the latest time in the forecast data.
	DataEndAt time.Time `json:"dataEndAt,omitempty"`
	// The estimated duration (in minutes) of the workload.  Defaults to the duration of a single forecast data point.
	WindowSize int32 `json:"windowSize,omitempty"`
	// The optimal forecasted data point within the 'forecastData' array.  Null if 'forecastData' array is empty.
	OptimalDataPoints []EmissionsDataDto `json:"optimalDataPoints,omitempty"`
	// The forecasted data points transformed and filtered to reflect the specified time and window parameters.  Points are ordered chronologically; Empty array if all data points were filtered out.  E.G. dataStartAt and dataEndAt times outside the forecast period; windowSize greater than total duration of forecast data;
	ForecastData []EmissionsDataDto `json:"forecastData,omitempty"`
}
