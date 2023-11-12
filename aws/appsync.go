package aws

type appsync struct {
	aws AWS
}

type AppSync interface {
	Regions() ([]string, error)
}

// NewAppSync creates a new AppSync client
func NewAppSync(aws AWS) AppSync {
	return &appsync{
		aws: aws,
	}
}

// Regions returns a list of regions that the AppSync service is
func (a appsync) Regions() ([]string, error) {
	return a.aws.Regions("appsync")
}
