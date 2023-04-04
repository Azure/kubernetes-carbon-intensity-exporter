/*
MIT License
Copyright (c) Microsoft Corporation.
*/
package exporter

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/Azure/kubernetes-carbon-intensity-exporter/pkg/sdk/client"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog/v2"
)

var (
	isImmutable = true
)

func (e *Exporter) CreateOrUpdateConfigMap(ctx context.Context, configMapName string, emissionForecast []client.EmissionsForecastDto) error {
	if emissionForecast == nil {
		return errors.New("emission forecast cannot be nil")
	}
	forecast := emissionForecast[0]
	binaryData, err := json.Marshal(forecast.ForecastData)
	if err != nil {
		return err
	}

	if forecast.ForecastData == nil || len(forecast.ForecastData) == 0 {
		return errors.New("forecast data cannot be nil or empty")
	}

	minForecast, maxForeCast := getMinMaxForecast(ctx, forecast.ForecastData)

	configMap := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Immutable: &isImmutable,
		Data: map[string]string{
			ConfigMapLastHeartbeatTime: time.Now().String(),                      // The latest time that the data exporter controller sends the data.
			ConfigMapMessage:           "",                                       // Additional information for user notification, if any.
			ConfigMapNumOfRecords:      strconv.Itoa(len(forecast.ForecastData)), // The number can be any value between 0 (no records for the current location) and 24(hours) * 12(5 min interval per hour).
			ConfigMapForecastDateTime:  forecast.DataStartAt.String(),            // The time when the data was started by the GSF SDK.
			ConfigMapMinForecast:       fmt.Sprintf("%f", minForecast),           // min forecast in the forecastData.
			ConfigMapMaxForecast:       fmt.Sprintf("%f", maxForeCast),           // max forecast in the forecastData.
		},
		BinaryData: map[string][]byte{
			"data": binaryData, // json marshal of the EmissionsData array.
		},
	}

	_, err = e.clusterClient.CoreV1().
		ConfigMaps(namespace).
		Create(ctx, configMap, v1.CreateOptions{})
	if err != nil {
		return err
	}
	klog.Infof("configMap %s has been created", configMapName)

	return nil
}

func (e *Exporter) GetGonfigMapWatch(ctx context.Context, configMapName string) watch.Interface {
	watch, err := e.clusterClient.CoreV1().
		ConfigMaps(namespace).
		Watch(ctx, v1.ListOptions{
			FieldSelector: "metadata.name=" + configMapName,
		})
	if err != nil {
		klog.Fatalf("unable to watch configMap %s, err: %v", configMapName, err)
	}
	return watch
}

func (e *Exporter) DeleteConfigmap(ctx context.Context, configMapName string) error {
	_, err := e.clusterClient.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, v1.GetOptions{})
	if err != nil {
		klog.Errorf("unable to get configMap %s", configMapName)
		return err
	}

	err = e.clusterClient.CoreV1().
		ConfigMaps(namespace).
		Delete(ctx, configMapName, v1.DeleteOptions{})
	if err != nil {
		klog.Errorf("unable to delete configMap %s", configMapName)
		return err
	}

	klog.Infof("configMap %s has been deleted", configMapName)
	return nil
}

func getMinMaxForecast(ctx context.Context, forecastData []client.EmissionsDataDto) (float64, float64) {
	values := make([]float64, len(forecastData))
	for index := range forecastData {
		values[index] = forecastData[index].Value
	}

	// Sort values
	sort.Float64s(values)
	return values[0], values[len(values)-1]
}
