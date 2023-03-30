package exporter

import (
	"context"
	"errors"
	"strconv"

	"github.com/Azure/kubernetes-carbon-intensity-exporter/pkg/sdk/client"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

var (
	configMapName = "carbonintensity"
	isImmutable   = true
)

func (e *Exporter) CreateOrUpdateConfigMap(ctx context.Context, emissionForecast []client.EmissionsForecastDto) error {
	if emissionForecast == nil {
		return errors.New("emission forecast cannot be nil")
	}
	forecast := emissionForecast[0]
	binaryData, err := json.Marshal(forecast.ForecastData)
	if err != nil {
		return err
	}
	forecastRecordNum := 0
	if forecast.ForecastData != nil {
		forecastRecordNum = len(forecast.ForecastData)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name: configMapName,
		},
		Immutable: &isImmutable,
		Data: map[string]string{
			"lastHeartbeatTime": forecast.DataEndAt.String(),     // The latest time that the data exporter controller sends the data.
			"message":           "",                              // Additional information for user notification, if any.
			"numOfRecords":      strconv.Itoa(forecastRecordNum), // The number can be any value between 0 (no records for the current location) and 24(hours) * 12(5 min interval per hour).
			"forecastDateTime":  forecast.DataStartAt.String(),   // The time when the data was started by the GSF SDK.
			"minForecast":       "",                              // min forecast in the binarydata.
			"maxForecast":       "",                              // max forecast in the binarydata.
		},
		BinaryData: map[string][]byte{
			"data": binaryData, // json marshal of the EmissionsData array.
		},
	}

	currentConfig, err := e.clusterClient.CoreV1().ConfigMaps("default").Get(ctx, configMapName, v1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			klog.Info("configmap is not found, creating a new one")
		} else {
			return err
		}
	}

	// Delete the old configmap if any.
	if currentConfig != nil && currentConfig.Name != "" || !apierrors.IsNotFound(err) {
		// Delete it first (as it is immutable)
		err = e.clusterClient.CoreV1().ConfigMaps("default").Delete(ctx, configMapName, v1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	_, err = e.clusterClient.CoreV1().ConfigMaps("default").Create(ctx, configMap, v1.CreateOptions{})
	if err != nil {
		return err
	}
	klog.Info("configmap has been Created")

	return nil
}
