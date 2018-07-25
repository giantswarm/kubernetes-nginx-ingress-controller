// +build k8srequired

package basic

import (
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/e2e-harness/pkg/framework/resource"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/setup"
)

var (
	f          *framework.Host
	helmClient *helmclient.Client
	l          micrologger.Logger
	r          *resource.Resource
)

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	var err error

	{
		c := micrologger.Config{}
		l, err = micrologger.New(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := framework.HostConfig{
			Logger:     l,
			ClusterID:  "someval",
			VaultToken: "someval",
		}
		f, err = framework.NewHost(c)
		if err != nil {
			panic(err.Error())
		}
	}

	{
		c := helmclient.Config{
			Logger:     l,
			K8sClient:  f.K8sClient(),
			RestConfig: f.RestConfig(),
		}
		helmClient, err = helmclient.New(c)
		if err != nil {
			panic(err.Error())
		}
	}

	resourceConfig := resource.ResourceConfig{
		Logger:     l,
		HelmClient: helmClient,
		Namespace:  "giantswarm",
	}
	r, err = resource.New(resourceConfig)
	if err != nil {
		panic(err.Error())
	}

	setup.WrapTestMain(f, helmClient, m)
}
