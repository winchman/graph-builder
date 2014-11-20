package buildgraph

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type jobSpec struct {
	Job           `yaml:",inline"`
	RequiresBlock []string `yaml:"requires"`
}

type configDocument struct {
	Blocks []*jobSpec
}

func uniqueStrings(strings []string) bool {
	check := make(map[string]bool)
	for _, s := range strings {
		if check[s] == true {
			return false
		}
		check[s] = true
	}
	return true
}

// ParseGraphFromYAML takes a YAML document describing a set of jobs
// to execute. It returns a graph parsed from that config. Note that
// Job names must be unique and no cycles can exist in the graph.
func ParseGraphFromYAML(input []byte) (g *Graph, err error) {
	config := configDocument{}

	err = yaml.Unmarshal(input, &config)
	if err != nil {
		return
	}

	jobMap := make(map[string]*Job)
	jobList := make([]*Job, len(config.Blocks))
	// Ensure job names are unique
	for i, b := range config.Blocks {
		if jobMap[b.Name] != nil {
			err = fmt.Errorf("More than one job with the name %s",
				b.Name)
			return
		}
		jobMap[b.Name] = &(b.Job)
		jobList[i] = &(b.Job)
	}

	// Setup dependancies
	for _, b := range config.Blocks {
		if !uniqueStrings(b.RequiresBlock) {
			err = fmt.Errorf("Job %s must contain unique requirements.",
				b.Name)
			return
		}
		for _, requiredName := range b.RequiresBlock {
			if jobMap[requiredName] == nil {
				err = fmt.Errorf("Job %s required by %s does not exist",
					requiredName, b.Name)
				return
			}
			jobMap[b.Name].Requires =
				append(jobMap[b.Name].Requires, jobMap[requiredName])
		}
	}

	return NewGraph(jobList)
}
