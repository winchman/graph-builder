package main

import (
	"fmt"
	"github.com/sylphon/graph-builder/buildgraph"
)

var (
	a = buildgraph.Job{
		Name: "Block-A",
	}
	b = buildgraph.Job{
		Name:     "Block-B",
		Requires: []*buildgraph.Job{&a},
	}
	c = buildgraph.Job{
		Name:     "Block-C",
		Requires: []*buildgraph.Job{&b},
	}

	validJobs = []*buildgraph.Job{&a, &b, &c}

	d = buildgraph.Job{
		Name: "Block-D",
	}
	e = buildgraph.Job{
		Name:     "Block-E",
		Requires: []*buildgraph.Job{&d},
	}
)

func main() {
	// Valid graph
	graph, err := buildgraph.NewGraph(validJobs)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Graph valid? ", graph.Validate())

	// Will fail due to cycle
	d.Requires = []*buildgraph.Job{&e}
	invalidJobs := []*buildgraph.Job{&d, &e}
	graph, err = buildgraph.NewGraph(invalidJobs)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Graph valid? ", graph.Validate())

}
