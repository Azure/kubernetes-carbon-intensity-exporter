package options

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/record"
	cliflag "k8s.io/component-base/cli/flag"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/klog/v2"

	providerconfig "github.com/Azure/sustainability/carbon-aware/cmd/carbon-data-provider/app/config"
)

type ProviderOptions struct {
	ClientConnection componentbaseconfig.ClientConnectionConfiguration
	Timeout          string

	Name     string
	Address  string
	Port     string
	CertFile string
	KeyFile  string
}

// NewResourceSyncerOptions creates a new resource syncer with a default config.
func NewProviderOptions() (*ProviderOptions, error) {
	return &ProviderOptions{
		ClientConnection: componentbaseconfig.ClientConnectionConfiguration{},
		Name:             "carbon-data-provider",
		Address:          "",
		Port:             "80",
		CertFile:         "",
		KeyFile:          "",
	}, nil
}

func (o *ProviderOptions) Flags() cliflag.NamedFlagSets {
	fss := cliflag.NamedFlagSets{}

	fs := fss.FlagSet("server")
	fs.StringVar(&o.ClientConnection.Kubeconfig, "master-kubeconfig", o.ClientConnection.Kubeconfig, "Path to kubeconfig file with authorization and control plane location information.")

	serverFlags := fss.FlagSet("metricsServer")
	serverFlags.StringVar(&o.Address, "address", o.Address, "The server address.")
	serverFlags.StringVar(&o.Port, "port", o.Port, "The server port.")
	serverFlags.StringVar(&o.CertFile, "cert-file", o.CertFile, "CertFile is the file containing x509 Certificate for HTTPS.")
	serverFlags.StringVar(&o.KeyFile, "key-file", o.KeyFile, "KeyFile is the file containing x509 private key matching certFile.")

	return fss
}

// Config return a syncer config object
func (o *ProviderOptions) Config() (*providerconfig.Config, error) {
	c := &providerconfig.Config{}

	// Prepare kube clients
	var (
		restConfig *restclient.Config
		err        error
	)
	restConfig, err = getClientConfig(o.ClientConnection, "", o.Timeout)
	if err != nil {
		return nil, err
	}
	clusterClient, err := clientset.NewForConfig(restclient.AddUserAgent(restConfig, "carbon-data-provider"))
	if err != nil {
		return nil, err
	}

	// Prepare event clients.
	eventBroadcaster := record.NewBroadcaster()
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "carbon-data-provider"})

	c.Kubeconfig = restConfig
	c.ClusterClient = clusterClient
	c.ClusterInformerFactory = informers.NewSharedInformerFactory(clusterClient, 0)
	c.Broadcaster = eventBroadcaster
	c.Recorder = recorder

	c.Address = o.Address
	c.Port = o.Port
	c.CertFile = o.CertFile
	c.KeyFile = o.KeyFile

	return c, nil
}

// getClientConfig creates a Kubernetes client rest config from the given config and serverAddrOverride.
func getClientConfig(config componentbaseconfig.ClientConnectionConfiguration, serverAddrOverride, timeout string) (*restclient.Config, error) {
	// This creates a client, first loading any specified kubeconfig
	// file, and then overriding the serverAddr flag, if non-empty.
	var (
		restConfig *restclient.Config
		err        error
	)
	if len(config.Kubeconfig) == 0 && len(serverAddrOverride) == 0 {
		klog.Info("Neither kubeconfig file nor control plane URL was specified. Falling back to in-cluster config.")
		restConfig, err = restclient.InClusterConfig()
	} else {
		// This creates a client, first loading any specified kubeconfig
		// file, and then overriding the serverAddr flag, if non-empty.
		restConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: config.Kubeconfig},
			&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: serverAddrOverride}}).ClientConfig()
	}

	if err != nil {
		return nil, err
	}

	// Allow Syncer CLI Flag timeout override
	if len(timeout) == 0 {
		if restConfig.Timeout == 0 {
			restConfig.Timeout = 30 * time.Second
		}
	} else {
		timeoutDuration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, err
		}

		restConfig.Timeout = timeoutDuration
	}

	restConfig.ContentConfig.ContentType = config.AcceptContentTypes
	restConfig.QPS = config.QPS
	if restConfig.QPS == 0 {
		restConfig.QPS = 100
	}
	restConfig.Burst = int(config.Burst)
	if restConfig.Burst == 0 {
		restConfig.Burst = 200
	}

	return restConfig, nil
}
