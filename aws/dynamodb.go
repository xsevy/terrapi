package aws

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type dynamoDB struct {
	client *dynamodb.DynamoDB
}

type DynamoDB interface {
	ListTables() ([]string, error)
}

func NewDynamoDB(aws AWS) DynamoDB {
	return &dynamoDB{
		client: dynamodb.New(aws.GetSession()),
	}
}

func (d *dynamoDB) ListTables() ([]string, error) {
	var tables []string

	result, err := d.client.ListTables(nil)
	if err != nil {
		return tables, err
	}

	for _, t := range result.TableNames {
		tables = append(tables, *t)
	}

	return tables, nil
}
