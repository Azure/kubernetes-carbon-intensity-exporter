package app

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // enable pprof in the server
	"os"

	"github.com/spf13/cobra"
	"k8s.io/apiserver/pkg/server/healthz"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/component-base/term"
	"k8s.io/component-base/version/verflag"
	"k8s.io/klog/v2"

	providerconfig "github.com/Azure/sustainability/carbon-aware/cmd/carbon-data-provider/app/config"
	"github.com/Azure/sustainability/carbon-aware/cmd/carbon-data-provider/app/options"
	"github.com/Azure/sustainability/carbon-aware/pkg/provider"
)

func NewProviderCommand(stopChan <-chan struct{}) *cobra.Command {
	s, err := options.NewProviderOptions()
	if err != nil {
		klog.Fatalf("unable to initialize command options: %v", err)
	}

	cmd := &cobra.Command{
		Use:  "carbon-data-provider",
		Long: `The carbon-data-provider is a controller that pulls carbon intensity data from ADLS`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var c *providerconfig.Config
			verflag.PrintAndExitIfRequested()

			c, err = s.Config()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if err := Run(c.Complete(), stopChan); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
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

func Run(cc *providerconfig.CompletedConfig, stopCh <-chan struct{}) error {
	p, err := provider.New(cc.ClusterClient, cc.Recorder)

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
	run := startProvider(p, stopCh)

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

func startProvider(p *provider.Provider, stopCh <-chan struct{}) func(context.Context) {
	return func(ctx context.Context) {
		p.Run(stopCh)
		<-ctx.Done()
	}
}
