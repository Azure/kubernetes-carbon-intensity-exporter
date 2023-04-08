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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

var (
	isImmutable = true
)

func (e *Exporter) CreateConfigMapFromEmissionForecast(ctx context.Context, configMapName string, emissionForecast []client.EmissionsForecastDto) error {
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

	return e.CreateConfigMapFromProperties(ctx, configMapName,
		map[string]string{
			ConfigMapLastHeartbeatTime: time.Now().String(),                      // The latest time that the data exporter controller sends the data.
			ConfigMapMessage:           "",                                       // Additional information for user notification, if any.
			ConfigMapNumOfRecords:      strconv.Itoa(len(forecast.ForecastData)), // The number can be any value between 0 (no records for the current location) and 24(hours) * 12(5 min interval per hour).
			ConfigMapForecastDateTime:  forecast.DataStartAt.String(),            // The time when the data was started by the GSF SDK.
			ConfigMapMinForecast:       fmt.Sprintf("%f", minForecast),           // min forecast in the forecastData.
			ConfigMapMaxForecast:       fmt.Sprintf("%f", maxForeCast),           // max forecast in the forecastData.
		}, binaryData)
}

func (e *Exporter) CreateConfigMapFromProperties(ctx context.Context, configMapName string, data map[string]string, binaryData []byte) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      configMapName,
			Namespace: client.Namespace,
		},
		Immutable: &isImmutable,
		Data:      data,
		BinaryData: map[string][]byte{
			BinaryData: binaryData, // json marshal of the EmissionsData array.
		},
	}
	_, err := e.clusterClient.CoreV1().
		ConfigMaps(client.Namespace).
		Create(ctx, configMap, v1.CreateOptions{})
	if err != nil {
		return err
	}
	klog.Infof("configMap %s has been created", configMapName)
	return nil
}

func (e *Exporter) DeleteConfigMap(ctx context.Context, configMapName string) error {
	currentConfigMap, err := e.GetConfigMap(ctx, configMapName)
	if err != nil {
		return err
	}

	if currentConfigMap == nil {
		return nil // configMap is not found, delete will not be called.
	}

	err = e.clusterClient.CoreV1().
		ConfigMaps(client.Namespace).
		Delete(ctx, configMapName, v1.DeleteOptions{})
	if err != nil {
		klog.Errorf("unable to delete configMap %s", configMapName)
		return err
	}

	klog.Infof("configMap %s has been deleted", configMapName)
	return nil
}

func (e *Exporter) GetConfigMap(ctx context.Context, configMapName string) (*corev1.ConfigMap, error) {
	currentConfigMap, err := e.clusterClient.CoreV1().ConfigMaps(client.Namespace).Get(ctx, configMapName, v1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) { // if configMap is not found, no errors will be returned.
			return nil, nil
		}
		return nil, err
	}

	return currentConfigMap, nil
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
