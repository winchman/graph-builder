package main

import (
	"errors"
)

type Graph struct {
	Jobs []*Job

	sort []*Job
}

func NewGraph(config string) (g *Graph, err error) {
	g = new(Graph)
	g.Jobs = make([]*Job, 1)

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

// http://en.wikipedia.org/wiki/Topological_sorting
func (g *Graph) topologicalSort() (sorted []*Job, err error) {
	sorted = make([]*Job, len(g.Jobs))

	tempMark := make(map[*Job]bool)
	permMark := make(map[*Job]bool)
	var visit func(*Job) error
	visit = func(j *Job) error {
		if tempMark[j] {
			return errors.New("Not a DAG!")
		}
		if (!tempMark[j]) && (!permMark[j]) {
			tempMark[j] = true
			for _, n := range j.Requires {
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

	for _, n := range g.Jobs {
		if permMark[n] == false {
			err := visit(n)
			if err != nil {
				return nil, err
			}
		}
	}
	return sorted, nil
}
