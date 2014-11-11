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

	_, err := buildgraph.NewGraph(validJobs)
	if err != nil {
		test.Error(err)
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

func TestGetDependants(test *testing.T) {
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

	output1, err := graph.GetDependants([]*buildgraph.Job{})
	if err != nil {
		test.Error(err)
	}

	if len(output1) != 1 || output1[0] != &a {
		test.Fail()
	}

	output2, err := graph.GetDependants([]*buildgraph.Job{&a})
	if err != nil {
		test.Error(err)
	}

	if len(output2) != 1 || output2[0] != &b {
		test.Fail()
	}

	output3, err := graph.GetDependants([]*buildgraph.Job{&a, &b})
	if err != nil {
		test.Error(err)
	}

	if len(output3) != 1 || output3[0] != &c {
		test.Fail()
	}

	output4, err := graph.GetDependants([]*buildgraph.Job{&a, &b, &c})
	if err != nil {
		test.Error(err)
	}

	if len(output4) != 0 {
		test.Fail()
	}

}
