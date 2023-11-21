# terrapi

## Usage

### Configuration
This is early stage of the project.
You need to provide configuration details.
1. Create S3 bucket and DynamoDB table for storing state of the project.

## Development
1. Install [precommit](https://pre-commit.com/#install)
2. Run `pre-commit install` in the project root directory
3. Run the application using `go run main.go`

## Todo
- directory picker
- improve design
- handle edge cases 
- setup column validation
- terraform state selector
- error and info messages
- support for other appsync authorizers
- other cloud providers like gcp and azure
- multiple workspaces support
- git init on creating api
- support for js resolvers
- Mutation resolver type
- removing resources
- support for multiple graphql resolver
- apollo federation support
- scrollable columns
- translations
