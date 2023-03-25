package main

import (
	"math/rand"
	"os"
	"time"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"

	exporterapp "github.com/Azure/kubernetes-carbon-intensity-exporter/cmd/exporter/app"
)

func mainMethod() error {
	rand.Seed(time.Now().UTC().UnixNano())

	logs.InitLogs()
	defer logs.FlushLogs()

	stopChan := genericapiserver.SetupSignalHandler()

	return exporterapp.NewExporterCommand(stopChan).Execute()
}

func main() {
	if mainMethod() != nil {
		os.Exit(1)
	}
}
