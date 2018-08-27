package managedservices

import "context"

// ChartConfig is the chart to test.
type ChartConfig struct {
	ChannelName string
	ChartName   string
	ChartValues string
	Namespace   string
}

// ChartResources are the key resources deployed by the chart.
type ChartResources struct {
	DaemonSets  []DaemonSet
	Deployments []Deployment
}

// DaemonSet is a daemonset to be tested.
type DaemonSet struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	MatchLabels map[string]string
	Replicas    int
}

// Deployment is a deployment to be tested.
type Deployment struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	MatchLabels map[string]string
	Replicas    int
}

type Interface interface {
	// Test executes the test of a managed services chart with basic
	// functionality that applies to all managed services charts.
	//
	// - Install chart.
	// - Check chart is deployed.
	// - Check key resources are correct.
	// - Run helm release tests.
	//
	Test(ctx context.Context, chartConfig ChartConfig, chartResources ChartResources) error
}
