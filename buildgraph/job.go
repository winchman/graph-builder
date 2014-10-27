package buildgraph

type Credentials struct {
	Username string
	Password string
}

type PushInfo struct {
	Image       string
	Credentials Credentials
}

type Job struct {
	Name         string
	Requires     []*Job
	Dockerfile   string
	ImageName    string
	DisableCache bool
	SkipPush     bool
	Tags         string
	PushInfo     PushInfo
}
