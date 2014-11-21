package buildgraph

// PushCredentials are used to authenticate with the Docker registry that is
// being pushed to.
type PushCredentials struct {
	Username string
	Password string
}

// PushInfo represents the Image name and credentials that should be used
// when pushing the result of this a Job.
type PushInfo struct {
	Image       string
	Credentials PushCredentials
}

// Job is the basic unit of computation in the graph. It must have a unique name.
type Job struct {
	Name         string
	Requires     []*Job `yaml:"-"`
	Dockerfile   string
	ImageName    string `yaml:"image_name"`
	DisableCache bool   `yaml:"disable_cache"`
	SkipPush     bool   `yaml:"skip_push"`
	Tags         string
	PushInfo     PushInfo `yaml:"push_info"`
}
