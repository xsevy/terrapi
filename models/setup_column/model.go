package setup_column

import (
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xsevy/terrapi/aws"
	"github.com/xsevy/terrapi/helpers"
	"github.com/xsevy/terrapi/helpers/bubbles"
	"github.com/xsevy/terrapi/helpers/models"
	"github.com/xsevy/terrapi/helpers/navigation"
	"github.com/xsevy/terrapi/messages"
	"github.com/xsevy/terrapi/styles"
)

type SetupColumnModel struct {
	id             string
	elements       []navigation.FormField
	lambdaClient   aws.Lambda
	appsyncClient  aws.AppSync
	s3Client       aws.S3
	dynamoDBClient aws.DynamoDB
	keys           helpers.KeyMap
	selected       navigation.Selected
	models.ColumnModel
}

func NewSetupColumnModel(
	lambdaClient aws.Lambda,
	appsyncClient aws.AppSync,
	s3Client aws.S3,
	dynamoDBClient aws.DynamoDB,
	focused bool,
) *SetupColumnModel {
	m := &SetupColumnModel{
		keys:           helpers.Keys,
		lambdaClient:   lambdaClient,
		appsyncClient:  appsyncClient,
		s3Client:       s3Client,
		dynamoDBClient: dynamoDBClient,
		selected:       0,
	}

	m.SetFocused(focused)
	return m
}

func (m *SetupColumnModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *SetupColumnModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Escape):
			return m, messages.SwitchColumn("select_column")
		case key.Matches(msg, m.keys.Tab):
			m.selected.Next(len(m.elements) - 1)
		case key.Matches(msg, m.keys.ShiftTab):
			m.selected.Prev()
		case key.Matches(msg, m.keys.Enter):
			if _, ok := m.elements[m.selected].(*bubbles.ButtonModel); ok {
				switch m.id {
				case helpers.ResourceIDs.CreateAppSyncAPI:
					cmd = messages.CreateResource(
						m.id,
						m.elements[0].Value(),
						messages.WithAWSRegion(m.elements[1].Value()),
						messages.WithBackendBucket(m.elements[2].Value()),
						messages.WithBackendLockTable(m.elements[3].Value()),
						messages.WithAuthorizerLambdaFunction(m.elements[4].Value()),
					)
				case helpers.ResourceIDs.CreateAppSyncDataSource:
					cmd = messages.CreateResource(
						m.id,
						m.elements[0].Value(),
						messages.WithLambdaRuntime(m.elements[1].Value()),
					)
				}

				return m, cmd
			}
		}
	}

	for i := range m.elements {
		if i == int(m.selected) {
			updatedModel, cmd := m.elements[i].Update(msg)
			cmds = append(cmds, cmd)

			if newElement, ok := updatedModel.(navigation.FormField); ok {
				m.elements[i] = newElement
			} else {
				panic("wrong type")
			}

			cmds = append(cmds, m.elements[i].Focus())
		} else {
			m.elements[i].Blur()
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *SetupColumnModel) View() string {
	views := []string{}
	for _, element := range m.elements {
		views = append(views, element.View())
	}
	content := strings.Join(views, "\n\n")

	if m.GetFocused() {
		content = styles.SetupColumnStyleFocused.Render(content)
	} else {
		content = styles.SetupColumnStyleBlured.Render("")
	}
	return content
}

func (m *SetupColumnModel) SetID(id string) {
	m.id = id
	m.setElements()
}

func (m *SetupColumnModel) setElements() {
	var wg sync.WaitGroup

	switch m.id {
	case helpers.ResourceIDs.CreateAppSyncAPI:
		var lambdaFunctions, appsyncRegions, s3Buckets, dynamoDBTables []string
		var lambdaFunctionsErr, appsyncRegionsErr, s3BucketsErr, dynamoDBTablesErr error

		wg.Add(4)

		go func() {
			defer wg.Done()
			lambdaFunctions, lambdaFunctionsErr = m.lambdaClient.ListFunctions()
			if lambdaFunctionsErr != nil {
				panic(lambdaFunctionsErr)
			}
		}()

		go func() {
			defer wg.Done()
			appsyncRegions, appsyncRegionsErr = m.appsyncClient.Regions()
			if appsyncRegionsErr != nil {
				panic(appsyncRegionsErr)
			}
		}()

		go func() {
			defer wg.Done()
			s3Buckets, s3BucketsErr = m.s3Client.ListBuckets()
			if s3BucketsErr != nil {
				panic(s3BucketsErr)
			}
		}()

		go func() {
			defer wg.Done()
			dynamoDBTables, dynamoDBTablesErr = m.dynamoDBClient.ListTables()
			if dynamoDBTablesErr != nil {
				panic(dynamoDBTablesErr)
			}
		}()

		wg.Wait()

		m.elements = []navigation.FormField{
			bubbles.NewTextInput("Name:", "name", 32),
			bubbles.NewListModel("Region:", appsyncRegions, true, 0),
			bubbles.NewListModel("Backend bucket:", s3Buckets, true, 0),
			bubbles.NewListModel("State lock:", dynamoDBTables, true, 0),
			bubbles.NewListModel("Authorizer function:", lambdaFunctions, true, 0),
			bubbles.NewButtonModel("Submit", false),
		}
	case helpers.ResourceIDs.CreateAppSyncDataSource:
		var lambdaRuntimes []string
		var lambdaRuntimesErr error

		wg.Add(1)

		go func() {
			defer wg.Done()
			lambdaRuntimes, lambdaRuntimesErr = m.lambdaClient.ListRuntimes()
			if lambdaRuntimesErr != nil {
				panic(lambdaRuntimesErr)
			}
		}()

		wg.Wait()

		m.elements = []navigation.FormField{
			bubbles.NewTextInput("Name:", "name", 32),
			bubbles.NewListModel("Runtime", lambdaRuntimes, true, 0),
			bubbles.NewButtonModel("Submit", false),
		}
	}
}
