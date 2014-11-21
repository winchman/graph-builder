package buildgraph_test

import (
	"testing"

	"github.com/sylphon/graph-builder/buildgraph"
)

const validConfig = `
blocks:
  - name: block-A
  - name: block-B
    requires:
      - block-A
  - name: block-C
    requires:
      - block-A
      - block-B
`

const cycleConfig = `
blocks:
  - name: block-A
    requires:
      - block-B
  - name: block-B
    requires:
      - block-A
`

const selfLoopConfig = `
blocks:
  - name: block-A
    requires:
      - block-A
`

const nonUniqueRequiresConfig = `
blocks:
  - name: block-A
  - name: block-B
    requires:
      - block-A
      - block-A
`

const nonUniqueBlockConfig = `
blocks:
  - name: block-A
  - name: block-A
`

const requiresNonExistConfig = `
blocks:
  - name: block-A
  - name: block-B
    requires:
      - block-C
`

func TestParseGraphValid(test *testing.T) {
	graph, err := buildgraph.ParseGraphFromYAML([]byte(validConfig))
	if err != nil {
		test.Error(err)
	}
	if graph == nil {
		test.Fail()
	}
	if len(graph.Jobs) != 3 {
		test.Fail()
	}
}

func TestParseGraphCycles(test *testing.T) {
	_, err := buildgraph.ParseGraphFromYAML([]byte(cycleConfig))
	if err == nil {
		test.Fail()
	}

	_, err = buildgraph.ParseGraphFromYAML([]byte(selfLoopConfig))
	if err == nil {
		test.Fail()
	}
}

func TestParseGraphNonUnique(test *testing.T) {
	_, err := buildgraph.ParseGraphFromYAML([]byte(nonUniqueBlockConfig))
	if err == nil {
		test.Fail()
	}

	_, err = buildgraph.ParseGraphFromYAML([]byte(nonUniqueRequiresConfig))
	if err == nil {
		test.Fail()
	}
}

func TestParseGraphNonExist(test *testing.T) {
	_, err := buildgraph.ParseGraphFromYAML([]byte(requiresNonExistConfig))
	if err == nil {
		test.Fail()
	}
}
