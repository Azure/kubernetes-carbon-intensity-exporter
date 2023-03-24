package config

import (
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

type Config struct {
	// the k8s cluster client
	ClusterClient          clientset.Interface
	ClusterInformerFactory informers.SharedInformerFactory

	// the rest config for the k8s cluster
	Kubeconfig *restclient.Config

	// the event sink
	Recorder    record.EventRecorder
	Broadcaster record.EventBroadcaster

	// server config.
	Address  string
	Port     string
	CertFile string
	KeyFile  string
}

type completedConfig struct {
	*Config
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

func (c *Config) Complete() *CompletedConfig {
	cc := completedConfig{c}
	return &CompletedConfig{&cc}
}
