package buildgraph

import (
	"errors"
)

type Graph struct {
	Jobs []*Job

	sort []*Job
}

type Config struct {
	Jobs []*Job `yaml:"blocks"`
}

func NewGraph(jobs []*Job) (g *Graph, err error) {
	g = new(Graph)

	g.Jobs = jobs

	return g, nil
}

// A Graph is valid iff it has a topological sort.
func (g *Graph) Validate() bool {
	sorted, err := g.topologicalSort()
	if err != nil {
		return false
	}
	g.sort = sorted

	return true
}

func (g *Graph) getDependants(j *Job) []*Job {
	output := make([]*Job, 0)
	for _, n := range g.Jobs {
		for _, r := range n.Requires {
			if r == j {
				output = append(output, n)
			}
		}
	}
	return output
}

// http://en.wikipedia.org/wiki/Topological_sorting
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
			for _, n := range g.getDependants(j) {
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
