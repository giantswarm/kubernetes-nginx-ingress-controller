// +build k8srequired

package basic

import (
	"fmt"
	"testing"

	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/env"
	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/templates"
)

const (
	testName = "basic"
)

func TestHelm(t *testing.T) {
	channel := fmt.Sprintf("%s-%s", env.CircleSHA(), testName)
	releaseName := "kubernetes-nginx-ingress-controller"

	err := r.InstallResource(releaseName, templates.NginxIngressControllerValues, channel)
	if err != nil {
		t.Fatalf("could not install %q %v", releaseName, err)
	}

	err = r.WaitForStatus(releaseName, "DEPLOYED")
	if err != nil {
		t.Fatalf("could not get release status of %q %v", releaseName, err)
	}
	l.Log("level", "debug", "message", fmt.Sprintf("%s succesfully deployed", releaseName))

	controllerName := "nginx-ingress-controller"
	controllerLabels := map[string]string{
		"app": controllerName,
		"giantswarm.io/service-type": "managed",
		"k8s-app":                    controllerName,
	}
	controllerMatchLabels := map[string]string{
		"k8s-app": controllerName,
	}
	err = d.Check(controllerName, 3, controllerLabels, controllerMatchLabels)
	if err != nil {
		t.Fatalf("controller manifest is incorrect: %v", err)
	}

	backendName := "default-http-backend"
	backendLabels := map[string]string{
		"giantswarm.io/service-type": "managed",
		"k8s-app":                    backendName,
	}
	backendMatchLabels := map[string]string{
		"k8s-app": backendName,
	}
	err = d.Check(backendName, 2, backendLabels, backendMatchLabels)
	if err != nil {
		t.Fatalf("default backend manifest is incorrect: %v", err)
	}

	err = helmClient.RunReleaseTest(releaseName)
	if err != nil {
		t.Fatalf("unexpected error during test of the chart: %v", err)
	}
}
