package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
)

type aws struct {
	Sess *session.Session
}

type AWS interface {
	GetSession() *session.Session
	Regions(service string) ([]string, error)
}

// NewAWS returns a new AWS interface
func NewAWS() AWS {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return aws{
		Sess: sess,
	}
}

// GetSession returns the session
func (a aws) GetSession() *session.Session {
	return a.Sess
}

// Regions returns a list of regions that the given service is
func (a aws) Regions(service string) ([]string, error) {
	sr, exists := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, service)
	if !exists {
		return nil, fmt.Errorf("service %s does not exist", service)
	}

	regions := make([]string, 0, len(sr))
	for _, region := range sr {
		r := region.ID()
		regions = append(regions, r)
	}

	return regions, nil
}
