// +build k8srequired

package integration

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
)

var (
	f *framework.Host
)

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	var v int
	var err error

	f, err = framework.NewHost(framework.HostConfig{})
	if err != nil {
		panic(err.Error())
	}

	err = f.CreateNamespace("giantswarm")
	if err != nil {
		log.Printf("%#v\n", err)
		v = 1
	}

	if v == 0 {
		v = m.Run()
	}

	if os.Getenv("KEEP_RESOURCES") != "true" {
		f.Teardown()
	}

	os.Exit(v)
}

func TestHelm(t *testing.T) {
	channel := os.Getenv("CIRCLE_SHA1")

	err := framework.HelmCmd(fmt.Sprintf("registry install quay.io/giantswarm/kubernetes-nginx-ingress-controller-chart:%s -n test-deploy -- --wait", channel))
	if err != nil {
		t.Errorf("unexpected error during installation of the chart: %v", err)
	}

	err = framework.HelmCmd("test --debug test-deploy")
	if err != nil {
		t.Errorf("unexpected error during test of the chart: %v", err)
	}
}
