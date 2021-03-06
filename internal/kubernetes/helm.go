package kubernetes

import (
	"os"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

// Chart is helm chart info
type Chart struct {
	Name    string
	Version string
}

// GetHelmChartsFromNamespaces fetches all charts from the namespaces
func GetHelmChartsFromNamespaces(namespaces []string, useLocally bool) []Chart {
	namespaces = getNamespaces(namespaces, getKubernetesClient(useLocally))

	var charts []Chart
	for _, namespace := range namespaces {
		settings := cli.New()
		actionConfig := new(action.Configuration)

		err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), log.Infof)
		if err != nil {
			log.WithError(err).Error("Failed to get Helm action config")
			continue
		}

		client := action.NewList(actionConfig)
		chartsInNamespace, err := client.Run()
		if err != nil {
			log.Errorf("Failed to run helm command: [%v]", err)
			continue
		}
		for _, chart := range chartsInNamespace {
			charts = append(charts, Chart{
				Name:    chart.Name,
				Version: chart.Chart.Metadata.Version,
			})
		}
	}
	return charts
}
