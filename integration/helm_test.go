// +build k8srequired

package integration

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	resourceNamespace = "kube-system"
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

	err := framework.HelmCmd(fmt.Sprintf("registry install --wait quay.io/giantswarm/kubernetes-nginx-ingress-controller-chart:%s -n test-deploy", channel))
	if err != nil {
		t.Errorf("unexpected error during installation of the chart: %v", err)
	}

	err = checkDeployment("nginx-ingress-controller", 3)
	if err != nil {
		t.Fatalf("controller manifest is incorrect: %v", err)
	}

	err = checkDeployment("default-http-backend", 2)
	if err != nil {
		t.Fatalf("default backend manifest is incorrect: %v", err)
	}

	err = framework.HelmCmd("test --debug test-deploy")
	if err != nil {
		t.Errorf("unexpected error during test of the chart: %v", err)
	}
}

// checkDeployment ensures that key properties of the deployment are correct.
func checkDeployment(name string, replicas int) error {
	expectedMatchLabels := map[string]string{
		"k8s-app": name,
	}
	expectedLabels := map[string]string{
		"k8s-app":                    name,
		"giantswarm.io/service-type": "managed",
	}

	c := f.K8sClient()
	ds, err := c.Apps().Deployments(resourceNamespace).Get(name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return microerror.Newf("could not find daemonset: '%s' %v", name, err)
	} else if err != nil {
		return microerror.Newf("unexpected error getting daemonset: %v", err)
	}

	// Check daemonset labels.
	if !reflect.DeepEqual(expectedLabels, ds.ObjectMeta.Labels) {
		return microerror.Newf("expected labels: %v got: %v", expectedLabels, ds.ObjectMeta.Labels)
	}

	// Check selector match labels.
	if !reflect.DeepEqual(expectedMatchLabels, ds.Spec.Selector.MatchLabels) {
		return microerror.Newf("expected match labels: %v got: %v", expectedMatchLabels, ds.Spec.Selector.MatchLabels)
	}

	// Check pod labels.
	if !reflect.DeepEqual(expectedLabels, ds.Spec.Template.ObjectMeta.Labels) {
		return microerror.Newf("expected pod labels: %v got: %v", expectedLabels, ds.Spec.Template.ObjectMeta.Labels)
	}

	// Check replica count.
	if *ds.Spec.Replicas != int32(replicas) {
		return microerror.Newf("expected replicas: %d got: %d", replicas, ds.Spec.Replicas)
	}

	return nil
}
