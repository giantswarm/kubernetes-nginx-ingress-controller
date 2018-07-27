// +build k8srequired

package basic

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/kubernetes-nginx-ingress-controller/integration/templates"
)

const (
	testName = "basic"
)

func TestHelm(t *testing.T) {
	channel := fmt.Sprintf("%s-%s", os.Getenv("CIRCLE_SHA1"), testName)
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
		"giantswarm.io/service.type": "managed",
		"k8s-app":                    controllerName,
	}
	controllerMatchLabels := map[string]string{
		"app":     controllerName,
		"k8s-app": controllerName,
	}
	err = checkDeployment(controllerName, 3, controllerLabels, controllerMatchLabels)
	if err != nil {
		t.Fatalf("controller manifest is incorrect: %v", err)
	}

	backendName := "default-http-backend"
	backendLabels := map[string]string{
		"app": backendName,
		"giantswarm.io/service.type": "managed",
		"k8s-app":                    backendName,
	}
	backendMatchLabels := map[string]string{
		"k8s-app": backendName,
	}
	err = checkDeployment(backendName, 2, backendLabels, backendMachLabels)
	if err != nil {
		t.Fatalf("default backend manifest is incorrect: %v", err)
	}

	err = helmClient.RunReleaseTest(releaseName)
	if err != nil {
		t.Fatalf("unexpected error during test of the chart: %v", err)
	}
}

// checkDeployment ensures that key properties of the deployment are correct.
func checkDeployment(name string, replicas int, labels, matchLabels) error {
	c := f.K8sClient()
	ds, err := c.Apps().Deployments(metav1.NamespaceSystem).Get(name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return microerror.Newf("could not find deployment: '%s' %v", name, err)
	} else if err != nil {
		return microerror.Newf("unexpected error getting deployment: %v", err)
	}

	// Check deployment labels.
	if !reflect.DeepEqual(labels, ds.ObjectMeta.Labels) {
		return microerror.Newf("expected labels: %v got: %v", expectedLabels, ds.ObjectMeta.Labels)
	}

	// Check selector match labels.
	if !reflect.DeepEqual(matchLabels, ds.Spec.Selector.MatchLabels) {
		return microerror.Newf("expected match labels: %v got: %v", expectedMatchLabels, ds.Spec.Selector.MatchLabels)
	}

	// Check pod labels.
	if !reflect.DeepEqual(labels, ds.Spec.Template.ObjectMeta.Labels) {
		return microerror.Newf("expected pod labels: %v got: %v", expectedLabels, ds.Spec.Template.ObjectMeta.Labels)
	}

	// Check replica count.
	if *ds.Spec.Replicas != int32(replicas) {
		return microerror.Newf("expected replicas: %d got: %d", replicas, ds.Spec.Replicas)
	}

	return nil
}
