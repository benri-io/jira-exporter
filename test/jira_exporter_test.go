package test

import (
	"encoding/json"
	"testing"

        "io/ioutil"
	"github.com/benri-io/jira-exporter/exporter"
)

// reads the file from a path. helper
func readFile(path string) string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// To test:
// Task  15
// Story 21
// Epic 6
// types=( "Task" "Story" "Epic" )
// for f in $types; do echo $f; ff=$(cat testdata/jira_response.json | jq '.issues[] | [.key , .fields.issuetype.name]' | grep $f | wc -l);  echo $ff; done
func TestUnmarshalIssuesResponse(t *testing.T) {
	var data = readFile("testdata/jira_response.json")
	var sr exporter.SearchResponse
	err := json.Unmarshal([]byte(data), &sr)
	if err != nil {
		panic(err)
	}
	if len(sr.Issues) != 50 {
		t.Fatalf("Failed to parse the correct number of issues. Got: %v", len(sr.Issues))
	}

	type TType struct {
		key      string
		expected int
	}

	// Check we parsed the correct types
	var types = []TType{TType{"Task", 15}, TType{"Story", 21}, TType{"Epic", 6}}
	for _, tf := range types {
		var taskFilter = exporter.IssueFilter{IssueType: tf.key}
		var filteredIssues = taskFilter.Filter(sr.Issues)
		if len(filteredIssues) != tf.expected {
			t.Fatalf("Failed to parse the correct number of task issues. Got: %v", len(filteredIssues))
		}

	}
}
