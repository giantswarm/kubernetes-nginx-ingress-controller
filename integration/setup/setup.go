// +build k8srequired

package setup

import (
	"log"
	"os"
	"testing"

	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/env"
)

func Setup(m *testing.M, config Config) {
	var v int
	var err error

	err = config.Host.CreateNamespace("giantswarm")
	if err != nil {
		log.Printf("%#v\n", err)
		v = 1
	}

	err = config.HelmClient.EnsureTillerInstalled()
	if err != nil {
		log.Printf("%#v\n", err)
		v = 1
	}

	if v == 0 {
		v = m.Run()
	}

	if env.KeepResources() != "true" {
		// Only do full teardown when not on CI.
		if env.CircleCI() != "true" {
			err := teardown(config)
			if err != nil {
				log.Printf("%#v\n", err)
				v = 1
			}
			// TODO there should be error handling for the framework teardown.
			config.Host.Teardown()
		}
	}

	os.Exit(v)
}
