package buildgraph_test

import (
	"testing"

	buildgraph "github.com/sylphon/graph-builder"
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

func TestGetDependants(test *testing.T) {
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

	successJobs := []*buildgraph.Job{}
	failedJobs := []*buildgraph.Job{}
	nextJobs, err := graph.GetDependants(successJobs, failedJobs)
	if err != nil {
		test.Error(err)
	}
	if len(nextJobs) != 1 || nextJobs[0] != &a {
		test.Fail()
	}
	successJobs = append(successJobs, &a)
	nextJobs, err = graph.GetDependants(successJobs, failedJobs)
	if err != nil {
		test.Error(err)
	}
	if len(nextJobs) != 1 || nextJobs[0] != &b {
		test.Fail()
	}
	successJobs = append(successJobs, &b)
	nextJobs, err = graph.GetDependants(successJobs, failedJobs)
	if err != nil {
		test.Error(err)
	}
	if len(nextJobs) != 1 || nextJobs[0] != &c {
		test.Fail()
	}
	successJobs = append(successJobs, &c)
	nextJobs, err = graph.GetDependants(successJobs, failedJobs)
	if err != nil {
		test.Error(err)
	}
	if len(nextJobs) != 0 {
		test.Fail()
	}

}

func TestGetDependantsWithFailures(test *testing.T) {
	// Graph with two chains
	// a -> b Success
	// c -> d Fail at c
	a := buildgraph.Job{
		Name: "Block-A",
	}
	b := buildgraph.Job{
		Name:     "Block-B",
		Requires: []*buildgraph.Job{&a},
	}
	c := buildgraph.Job{
		Name: "Block-C",
	}
	d := buildgraph.Job{
		Name:     "Block-D",
		Requires: []*buildgraph.Job{&c},
	}
	validJobs := []*buildgraph.Job{&a, &b, &c, &d}

	graph, err := buildgraph.NewGraph(validJobs)
	if err != nil {
		test.Error(err)
	}

	if !graph.Validate() {
		test.Fail()
	}

	successJobs := []*buildgraph.Job{}
	failedJobs := []*buildgraph.Job{}
	nextJobs, err := graph.GetDependants(successJobs, failedJobs)
	if err != nil {
		test.Error(err)
	}
	if len(nextJobs) != 2 {
		test.Fail()
	}
	successJobs = append(successJobs, &a)
	failedJobs = append(failedJobs, &c)
	nextJobs, err = graph.GetDependants(successJobs, failedJobs)
	if err != nil {
		test.Error(err)
	}
	if len(nextJobs) != 1 || nextJobs[0] != &b {
		test.Fail()
	}
	successJobs = append(successJobs, &b)
	failedJobs = append(failedJobs, &c)
	nextJobs, err = graph.GetDependants(successJobs, failedJobs)
	if err != nil {
		test.Error(err)
	}
	if len(nextJobs) != 0 {
		test.Fail()
	}
}

func TestGetDependantsWithFailSuccessJob(test *testing.T) {
	// a -> b
	a := buildgraph.Job{
		Name: "Block-A",
	}
	b := buildgraph.Job{
		Name:     "Block-B",
		Requires: []*buildgraph.Job{&a},
	}
	validJobs := []*buildgraph.Job{&a, &b}

	graph, err := buildgraph.NewGraph(validJobs)
	if err != nil {
		test.Error(err)
	}

	if !graph.Validate() {
		test.Fail()
	}

	successJobs := []*buildgraph.Job{&a}
	failedJobs := []*buildgraph.Job{&a}
	_, err = graph.GetDependants(successJobs, failedJobs)
	if err == nil {
		test.Fail()
	}
}

func TestGetDependantsWithNonValidatedGraph(test *testing.T) {
	// a -> b
	a := buildgraph.Job{
		Name: "Block-A",
	}
	b := buildgraph.Job{
		Name:     "Block-B",
		Requires: []*buildgraph.Job{&a},
	}
	validJobs := []*buildgraph.Job{&a, &b}

	graph := buildgraph.Graph{Jobs: validJobs}

	successJobs := []*buildgraph.Job{&a}
	failedJobs := []*buildgraph.Job{&a}
	_, err := graph.GetDependants(successJobs, failedJobs)
	if err == nil {
		test.Fail()
	}
}
