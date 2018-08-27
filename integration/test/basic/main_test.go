// +build k8srequired

package basic

import (
	"fmt"
	"testing"

	"github.com/giantswarm/apprclient"
	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/e2etests/managedservices"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/env"
	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/setup"
	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/templates"
)

const (
	chartName          = "kubernetes-nginx-ingress-controller"
	controllerName     = "nginx-ingress-controller"
	defaultBackendName = "default-http-backend"
	testName           = "basic"
)

var (
	a          *apprclient.Client
	ms         *managedservices.ManagedServices
	h          *framework.Host
	helmClient *helmclient.Client
	l          micrologger.Logger
)

func init() {
	var err error

	{
		c := micrologger.Config{}
		l, err = micrologger.New(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := apprclient.Config{
			Fs:     afero.NewOsFs(),
			Logger: l,

			Address:      "https://quay.io",
			Organization: "giantswarm",
		}
		a, err = apprclient.New(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := framework.HostConfig{
			Logger:     l,
			ClusterID:  "na",
			VaultToken: "na",
		}
		h, err = framework.NewHost(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := helmclient.Config{
			Logger:          l,
			K8sClient:       h.K8sClient(),
			RestConfig:      h.RestConfig(),
			TillerNamespace: "giantswarm",
		}
		helmClient, err = helmclient.New(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := managedservices.Config{
			ApprClient:    a,
			HelmClient:    helmClient,
			HostFramework: h,
			Logger:        l,

			ChartConfig: managedservices.ChartConfig{
				ChannelName: fmt.Sprintf("%s-%s", env.CircleSHA(), testName),
				ChartName:   chartName,
				ChartValues: templates.NginxIngressControllerBasicValues,
				Namespace:   metav1.NamespaceSystem,
			},
			ChartResources: managedservices.ChartResources{
				Deployments: []managedservices.Deployment{
					{
						Name:      controllerName,
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							"app": controllerName,
							"giantswarm.io/service-type": "managed",
							"k8s-app":                    controllerName,
						},
						MatchLabels: map[string]string{
							"k8s-app": controllerName,
						},
						Replicas: 3,
					},
					{
						Name:      defaultBackendName,
						Namespace: metav1.NamespaceSystem,
						Labels: map[string]string{
							"app": defaultBackendName,
							"giantswarm.io/service-type": "managed",
							"k8s-app":                    defaultBackendName,
						},
						MatchLabels: map[string]string{
							"k8s-app": defaultBackendName,
						},
						Replicas: 2,
					},
				},
			},
		}
		ms, err = managedservices.New(c)
		if err != nil {
			panic(err.Error())
		}
	}
}

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	{
		c := setup.Config{
			HelmClient: helmClient,
			Host:       h,
		}

		setup.Setup(m, c)
	}
}
