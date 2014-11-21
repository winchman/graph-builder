package buildgraph_test

import (
	"testing"

	buildgraph "github.com/sylphon/graph-builder"
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

const exampleConfig = `
blocks:
  - name: block-A
    skip_push: true
    disable_cache: true
  - name: block-B
    requires:
      - block-A
    disable_cache: true
    tags:
      - latest
    push_info:
      image: quay.io/namespace/repo:latest
      credentials:
        username: fakeuser
        password: fakepass
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

func TestParseGraphExampleYAML(test *testing.T) {
	graph, err := buildgraph.ParseGraphFromYAML([]byte(exampleConfig))
	if err != nil {
		test.Error(err)
	}

	if len(graph.Jobs) != 2 {
		test.Fail()
	}

	a := graph.Jobs[0]
	b := graph.Jobs[1]
	if a.Name != "block-A" ||
		a.SkipPush != true ||
		a.DisableCache != true {
		test.Fail()
	}

	if b.Name != "block-B" ||
		b.Requires[0] != a ||
		len(b.Tags) != 1 ||
		b.Tags[0] != "latest" ||
		b.DisableCache != true ||
		b.PushInfo.Image != "quay.io/namespace/repo:latest" ||
		b.PushInfo.Credentials.Username != "fakeuser" ||
		b.PushInfo.Credentials.Password != "fakepass" {
		test.Fail()
	}
}
