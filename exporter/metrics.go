package exporter

import (
	"strconv"

	log "github.com/benri-io/jira-exporter/logger"
	"github.com/prometheus/client_golang/prometheus"
)

// AddMetrics - Add's all of the metrics to a map of strings, returns the map.
func AddMetrics() map[string]*prometheus.Desc {
	log.GetDefaultLogger().Info("Setting Up Metrics")
	APIMetrics := make(map[string]*prometheus.Desc)
	APIMetrics["Issues"] = prometheus.NewDesc(
		prometheus.BuildFQName("jira", "project", "issue"),
		"A reference to a ticket on a JIRA metric",
		[]string{"project", "creator",
			"assignee", "priority", "status", "status_key",
			"reporter", "issue_type", "hierarchy"}, nil,
	)
	log.GetDefaultLogger().Info("Finished Adding Metrics")
	return APIMetrics
}

// processMetrics - processes the response data and sets the metrics using it as a source
func (e *Exporter) processMetrics(data []*Datum, rates *RateLimits, ch chan<- prometheus.Metric) error {

	log.GetDefaultLogger().Info("Processing Metrics")
	defer log.GetDefaultLogger().Info("Done Processing Metrics")

	// APIMetrics - range through the data slice
	for _, x := range data {

		log.GetDefaultLogger().Infof("Processing %d issues.", len(x.Issues))
		for _, issue := range x.Issues {

			ch <- prometheus.MustNewConstMetric(e.APIMetrics["Issues"],
				prometheus.CounterValue,
				1,
				issue.Fields.Project.Name,
				issue.Fields.Creator.DisplayName,
				issue.Fields.Assignee.DisplayName,
				issue.Fields.Priority.Name,
				issue.Fields.Status.StatusCategory.Name,
				issue.Fields.Status.StatusCategory.Key,
				issue.Fields.Reporter.DisplayName,
				issue.Fields.IssueType.Name,
				strconv.Itoa(issue.Fields.IssueType.HeirarchyLevel),
			)
		}
	}
	return nil
}
