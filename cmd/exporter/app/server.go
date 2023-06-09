/*
MIT License
Copyright (c) Microsoft Corporation.
*/
package app

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // enable pprof in the server
	"time"

	"github.com/Azure/kubernetes-carbon-intensity-exporter/pkg/sdk/client"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/server/healthz"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/component-base/term"
	"k8s.io/component-base/version/verflag"
	"k8s.io/klog/v2"

	exporterconfig "github.com/Azure/kubernetes-carbon-intensity-exporter/cmd/exporter/app/config"
	"github.com/Azure/kubernetes-carbon-intensity-exporter/cmd/exporter/app/options"
	"github.com/Azure/kubernetes-carbon-intensity-exporter/pkg/exporter"
)

var (
	//exporter command args
	configMapName  = flag.String("configmap-name", "carbon-intensity", "Configmap name - Default 'carbonIntensity'")
	patrolInterval = flag.String("patrol-interval", "12h", "Patrol interval in hours - Default every 12 hours")
	region         = flag.String("region", "", "Region to get carbon intensity for - Required")
)

func NewExporterCommand(stopChan <-chan struct{}) *cobra.Command {

	s, err := options.NewExporterOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}

	cmd := &cobra.Command{
		Use:  "carbon-data-exporter",
		Long: `The carbon-data-exporter is a controller that pulls carbon intensity data from GSF API server`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var c *exporterconfig.Config
			verflag.PrintAndExitIfRequested()

			c, err = s.Config()
			if err != nil {
				klog.Fatalf("unable to initialize command configs: %s", err.Error())
			}
			if err := Run(c.Complete(), stopChan); err != nil {
				klog.Fatalf("unable to execute command : %s", err.Error())
			}
		},
	}

	fs := cmd.Flags()
	namedFlagSets := s.Flags()
	verflag.AddFlags(namedFlagSets.FlagSet("global"))
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())

	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})

	return cmd
}

func Run(cc *exporterconfig.CompletedConfig, stopCh <-chan struct{}) error {
	// Init client SDK and exporter
	apiClient := client.NewAPIClient(client.NewConfiguration())
	e, err := exporter.New(cc.ClusterClient, apiClient, cc.Recorder)
	if err != nil {
		return fmt.Errorf("new syncer: %v", err)
	}

	// Prepare the event broadcaster.
	if cc.Broadcaster != nil && cc.ClusterClient != nil {
		cc.Broadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: cc.ClusterClient.CoreV1().Events("")})
	}

	// Start all informers.
	cc.ClusterInformerFactory.Start(stopCh)

	// Wait for all caches to sync before resource sync.
	cc.ClusterInformerFactory.WaitForCacheSync(stopCh)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// Prepare a reusable runCommand function.
	run := startExporter(e, stopCh)

	go func() {
		select {
		case <-stopCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	go func() {
		// start a pprof http server
		klog.Fatal(http.ListenAndServe(":6060", nil))
	}()

	go func() {
		// start a health http server.
		mux := http.NewServeMux()
		healthz.InstallHandler(mux)
		klog.Fatal(http.ListenAndServe(":8080", mux))
	}()

	run(ctx)
	return fmt.Errorf("finished without leader elect")
}

func startExporter(p *exporter.Exporter, stopCh <-chan struct{}) func(context.Context) {
	// Parse patrolInterval to time.Duration
	ptDuration, err := time.ParseDuration(*patrolInterval)
	if err != nil {
		return func(ctx context.Context) {
			klog.Fatalf("an error while parsing patrol-interval, err: %s", err.Error())
			ctx.Err()
		}
	}

	return func(ctx context.Context) {
		p.Run(ctx, *configMapName, *region, ptDuration, stopCh)
		<-ctx.Done()
	}
}
