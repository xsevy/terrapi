package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xsevy/terrapi/aws"
	"github.com/xsevy/terrapi/models/main_model"
	"github.com/xsevy/terrapi/models/menu"
	"github.com/xsevy/terrapi/models/select_column"
	"github.com/xsevy/terrapi/models/select_column_choices"
	"github.com/xsevy/terrapi/models/setup_column"
)

func main() {
	awsClient := aws.NewAWS()
	lambdaClient := aws.NewLambda(awsClient)
	appsyncClient := aws.NewAppSync(awsClient)
	s3Client := aws.NewS3(awsClient)
	dynamoDBClient := aws.NewDynamoDB(awsClient)

	selectColumnChoices := select_column_choices.NewSelectColumnChoicesModel()
	selectColumn := select_column.NewSelectColumnModel(selectColumnChoices, true)
	setup_column := setup_column.NewSetupColumnModel(
		lambdaClient,
		appsyncClient,
		s3Client,
		dynamoDBClient,
		false,
	)
	menu := menu.NewMenuModel(selectColumn, setup_column)
	main := main_model.NewMainModel(menu)

	p := tea.NewProgram(main)
	_, err := p.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
