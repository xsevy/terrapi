package messages

import tea "github.com/charmbracelet/bubbletea"

type CreateResourceMsg struct {
	ID               string
	ProjectName      string
	LambdaRuntime    string
	AWSRegion        string
	BackendBucket    string
	BackendLockTable string
}

type createResourceOption func(*CreateResourceMsg)

func CreateResource(id string, ProjectName string, options ...createResourceOption) tea.Cmd {
	return func() tea.Msg {
		msg := &CreateResourceMsg{
			ID:          id,
			ProjectName: ProjectName,
		}
		for _, opt := range options {
			opt(msg)
		}
		return *msg
	}
}

func WithLambdaRuntime(runtime string) createResourceOption {
	return func(msg *CreateResourceMsg) {
		msg.LambdaRuntime = runtime
	}
}

func WithAWSRegion(region string) createResourceOption {
	return func(msg *CreateResourceMsg) {
		msg.AWSRegion = region
	}
}

func WithBackendBucket(bucket string) createResourceOption {
	return func(msg *CreateResourceMsg) {
		msg.BackendBucket = bucket
	}
}

func WithBackendLockTable(table string) createResourceOption {
	return func(msg *CreateResourceMsg) {
		msg.BackendLockTable = table
	}
}
