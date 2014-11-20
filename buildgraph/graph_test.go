package buildgraph_test

import (
	"testing"

	"github.com/sylphon/graph-builder/buildgraph"
)

func TestValidGraph(test *testing.T) {
	// Valid graph
	a := buildgraph.Job{
		Name: "Block-A",
	}
	b := buildgraph.Job{
		Name:     "Block-B",
		Requires: []*buildgraph.Job{&a},
	}
	c := buildgraph.Job{
		Name:     "Block-C",
		Requires: []*buildgraph.Job{&b},
	}
	validJobs := []*buildgraph.Job{&a, &b, &c}

	graph, err := buildgraph.NewGraph(validJobs)
	if err != nil {
		test.Error(err)
	}

	if !graph.Validate() {
		test.Fail()
	}
}

func TestInvalidGraph(test *testing.T) {
	// Will fail due to cycle
	d := buildgraph.Job{
		Name: "Block-D",
	}
	e := buildgraph.Job{
		Name:     "Block-E",
		Requires: []*buildgraph.Job{&d},
	}
	d.Requires = []*buildgraph.Job{&e}
	invalidJobs := []*buildgraph.Job{&d, &e}
	_, err := buildgraph.NewGraph(invalidJobs)
	if err == nil {
		test.Fail()
	}
}
