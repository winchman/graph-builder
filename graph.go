package buildgraph

import (
	"errors"
)

// Graph represents the dependancy graph of Jobs. Use NewGraph to construct one.
type Graph struct {
	Jobs []*Job

	sort      []*Job
	validated bool
}

// NewGraph accepts a list of Job structs and convert that into a Graph object.
// If the resulting graph is not valid, an error is returned.
func NewGraph(jobs []*Job) (g *Graph, err error) {
	g = new(Graph)

	g.Jobs = jobs

	if !g.Validate() {
		return g, errors.New("graph is invalid")
	}

	return g, nil
}

// Validate returns true if the graph has a valid topological
// sort. This ensures that the graph is a DAG. Doing this also means
// that the graph has a sort that can be accessed later.
func (g *Graph) Validate() bool {
	sorted, err := g.topologicalSort()
	if err != nil {
		return false
	}

	g.sort = sorted
	g.validated = true

	return true
}

// GetDependants returns all dependant jobs for the given jobs. Note
// that you must pass all prerequisite jobs, not just some of the
// completed jobs. To get the initial list of jobs with no
// dependancies, just pass an empty list.
func (g *Graph) GetDependants(completedJobs []*Job) ([]*Job, error) {
	var output []*Job

	if !g.validated {
		return output, errors.New("graph not sorted yet, call graph.Validate()")
	}

	completedMap := make(map[*Job]bool)
	for _, j := range completedJobs {
		completedMap[j] = true
	}
SortLoop:
	for _, n := range g.sort {
		for _, r := range n.Requires {
			if !completedMap[r] {
				break SortLoop
			}
		}
		if !completedMap[n] {
			output = append(output, n)
		}
	}
	return output, nil
}

// jobsDependantOnSingleJob returns all dependant jobs for a single job j.
func (g *Graph) jobsDependantOnSingleJob(j *Job) []*Job {
	var output []*Job

	for _, n := range g.Jobs {
		for _, r := range n.Requires {
			if r == j {
				output = append(output, n)
			}
		}
	}
	return output
}

// Citation: http://en.wikipedia.org/wiki/Topological_sorting
func (g *Graph) topologicalSort() (sorted []*Job, err error) {
	sorted = make([]*Job, 0)

	tempMark := make(map[*Job]bool)
	permMark := make(map[*Job]bool)
	// Visit function to complete at each node in DFS
	var visit func(*Job) error
	visit = func(j *Job) error {
		if tempMark[j] {
			return errors.New("Not a DAG!")
		}
		if (!tempMark[j]) && (!permMark[j]) {
			tempMark[j] = true
			for _, n := range g.jobsDependantOnSingleJob(j) {
				err := visit(n)
				if err != nil {
					return err
				}
			}
			permMark[j] = true
			tempMark[j] = false
			sorted = append([]*Job{j}, sorted...)
		}
		return nil
	}

	// Complete for each node.
	for _, n := range g.Jobs {
		if (!tempMark[n]) && (!permMark[n]) {
			err := visit(n)
			if err != nil {
				return nil, err
			}
		}
	}
	return sorted, nil
}
