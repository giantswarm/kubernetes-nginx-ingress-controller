// +build k8srequired

package setup

import (
	"log"
	"os"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/helmclient"

	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/teardown"
)

func WrapTestMain(h *framework.Host, helmClient *helmclient.Client, m *testing.M) {
	var v int
	var err error

	err = h.CreateNamespace("giantswarm")
	if err != nil {
		log.Printf("%#v\n", err)
		v = 1
	}

	err = helmClient.EnsureTillerInstalled()
	if err != nil {
		log.Printf("%#v\n", err)
		v = 1
	}

	if v == 0 {
		v = m.Run()
	}

	if os.Getenv("KEEP_RESOURCES") != "true" {
		// only do full teardown when not on CI
		if os.Getenv("CIRCLECI") != "true" {
			err := teardown.Teardown(h, helmClient)
			if err != nil {
				log.Printf("%#v\n", err)
				v = 1
			}
			// TODO there should be error handling for the framework teardown.
			h.Teardown()
		}
	}

	os.Exit(v)
}
