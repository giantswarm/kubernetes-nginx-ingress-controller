// +build k8srequired

package migration

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/e2e-harness/pkg/framework/resource"
	e2esetup "github.com/giantswarm/e2esetup/chart"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/micrologger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	h          *framework.Host
	helmClient *helmclient.Client
	l          micrologger.Logger
	r          *resource.Resource
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

	resourceConfig := resource.ResourceConfig{
		Logger:     l,
		HelmClient: helmClient,
		Namespace:  metav1.NamespaceSystem,
	}
	r, err = resource.New(resourceConfig)
	if err != nil {
		panic(err.Error())
	}
}

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	ctx := context.Background()

	{
		c := e2esetup.Config{
			HelmClient: helmClient,
			Host:       h,
		}

		err := e2esetup.Setup(ctx, m, c)
		if err != nil {
			l.LogCtx(ctx, "level", "error", "message", "e2e test failed", "stack", fmt.Sprintf("%#v\n", err))
			os.Exit(1)
		}
	}
}
