package exporter

import (
	"fmt"
	"hash/fnv"
	"strconv"

	log "github.com/benri-io/jira-exporter/logger"
	"github.com/prometheus/client_golang/prometheus"
)

// AddMetrics - Add's all of the metrics to a map of strings, returns the map.
func AddMetrics() map[string]*prometheus.Desc {
	APIMetrics := make(map[string]*prometheus.Desc)
	APIMetrics["Issues"] = prometheus.NewDesc(
		prometheus.BuildFQName("jira", "project", "issue"),
		"A counter for jira issues",
		[]string{"project", "creator",
			"assignee", "priority", "status", "status_key",
			"reporter", "issue_type", "hierarchy"}, nil,
	)
	return APIMetrics
}

type fieldCounterPair struct {
	v     int
	issue Issue
}

type IssueCounter struct {
	m map[uint32]fieldCounterPair
}

// Aggregates based on hash of labels
func (ic *IssueCounter) add(i Issue) {

	var s = fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s:%d",
		i.Fields.Project.Name,
		i.Fields.Creator.DisplayName,
		i.Fields.Assignee.DisplayName,
		i.Fields.Priority.Name,
		i.Fields.Status.StatusCategory.Name,
		i.Fields.Status.StatusCategory.Key,
		i.Fields.Reporter.DisplayName,
		i.Fields.IssueType.Name,
		i.Fields.IssueType.HeirarchyLevel)

	hs := hash(s)
	if entry, ok := ic.m[hs]; !ok {
		ic.m[hs] = fieldCounterPair{1, i}
	} else {
		entry.v += 1
		ic.m[hs] = entry
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// processMetrics - processes the response data and sets the metrics using it as a source
func (e *Exporter) processMetrics(data []*Datum, rates *RateLimits, ch chan<- prometheus.Metric) error {
	ic := IssueCounter{
		m: make(map[uint32]fieldCounterPair),
	}

	log.GetDefaultLogger().Info("Processing Metrics")
	defer log.GetDefaultLogger().Info("Done Processing Metrics")
	// APIMetrics - range through the data slice
	for _, x := range data {
		log.GetDefaultLogger().Infof("Processing %d issues.", len(x.Issues))
		for _, issue := range x.Issues {
			ic.add(issue)
		}
		for _, k := range ic.m {
			ch <- prometheus.MustNewConstMetric(e.APIMetrics["Issues"],
				prometheus.CounterValue,
				float64(k.v),
				k.issue.Fields.Project.Name,
				k.issue.Fields.Creator.DisplayName,
				k.issue.Fields.Assignee.DisplayName,
				k.issue.Fields.Priority.Name,
				k.issue.Fields.Status.StatusCategory.Name,
				k.issue.Fields.Status.StatusCategory.Key,
				k.issue.Fields.Reporter.DisplayName,
				k.issue.Fields.IssueType.Name,
				strconv.Itoa(k.issue.Fields.IssueType.HeirarchyLevel),
			)
		}
	}
	return nil
}
