/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"os"

	operatorsv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorsv1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	integrationv1 "github.com/redhat-integration/integration-operator/api/v1"
	"github.com/redhat-integration/integration-operator/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(operatorsv1.AddToScheme(scheme))
	utilruntime.Must(operatorsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(integrationv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	config, err := createInstallationConfig()
	if err != nil {
		setupLog.Error(err, "unable to read configuration")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "2d830203.redhat.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.InstallationReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Installation"),
		Scheme: mgr.GetScheme(),
		Config: config,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Installation")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func createInstallationConfig() (*controllers.InstallationConfig, error) {
	config := &controllers.InstallationConfig{}
	var ok bool

	if config.Channel3scale, ok = os.LookupEnv("CHANNEL_3SCALE"); !ok {
		return nil, fmt.Errorf("missing CHANNEL_3SCALE environment variable")
	}
	if config.Channel3scaleAPIcast, ok = os.LookupEnv("CHANNEL_3SCALE_APICAST"); !ok {
		return nil, fmt.Errorf("missing CHANNEL_3SCALE_APICAST environment variable")
	}
	if config.ChannelAMQBroker, ok = os.LookupEnv("CHANNEL_AMQ_BROKER"); !ok {
		return nil, fmt.Errorf("missing CHANNEL_AMQ_BROKER environment variable")
	}
	if config.ChannelAMQInterconnect, ok = os.LookupEnv("CHANNEL_AMQ_INTERCONNECT"); !ok {
		return nil, fmt.Errorf("missing CHANNEL_AMQ_INTERCONNECT environment variable")
	}
	if config.ChannelAMQStreams, ok = os.LookupEnv("CHANNEL_AMQ_STREAMS"); !ok {
		return nil, fmt.Errorf("missing CHANNEL_AMQ_STREAMS environment variable")
	}
	if config.ChannelAPIDesigner, ok = os.LookupEnv("CHANNEL_API_DESIGNER"); !ok {
		return nil, fmt.Errorf("missing CHANNEL_API_DESIGNER environment variable")
	}
	if config.ChannelCamelK, ok = os.LookupEnv("CHANNEL_CAMEL_K"); !ok {
		return nil, fmt.Errorf("missing CHANNEL_CAMEL_K environment variable")
	}
	if config.ChannelFuseConsole, ok = os.LookupEnv("CHANNEL_FUSE_CONSOLE"); !ok {
		return nil, fmt.Errorf("missing CHANNEL_FUSE_CONSOLE environment variable")
	}
	if config.ChannelFuseOnline, ok = os.LookupEnv("CHANNEL_FUSE_ONLINE"); !ok {
		return nil, fmt.Errorf("missing CHANNEL_FUSE_ONLINE environment variable")
	}
	if config.ChannelServiceRegistry, ok = os.LookupEnv("CHANNEL_SERVICE_REGISTRY"); !ok {
		return nil, fmt.Errorf("missing CHANNEL_SERVICE_REGISTRY environment variable")
	}

	return config, nil
}
