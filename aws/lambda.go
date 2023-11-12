package aws

import (
	lambda_sdk "github.com/aws/aws-sdk-go/service/lambda"
)

type lambda struct {
	client *lambda_sdk.Lambda
}

type Lambda interface {
	ListFunctions() ([]string, error)
	ListRuntimes() ([]string, error)
}

// NewLambda creates a new Lambda client
func NewLambda(aws AWS) Lambda {
	return &lambda{
		client: lambda_sdk.New(aws.GetSession()),
	}
}

// ListFunctions lists all lambda functions
func (l *lambda) ListFunctions() ([]string, error) {
	var funcs []string

	result, err := l.client.ListFunctions(nil)
	if err != nil {
		return funcs, err
	}

	for _, f := range result.Functions {
		funcs = append(funcs, *f.FunctionName)
	}

	return funcs, nil
}

func (l *lambda) ListRuntimes() ([]string, error) {
	return []string{"python3.11"}, nil
}
