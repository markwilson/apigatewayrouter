package apigatewayrouter

import "github.com/aws/aws-lambda-go/events"

// Route determines if it can handle an APIGateway request event and handles it if possible.
type Route struct {
	Match  MatchFunc
	Handle HandleFunc
}

// MatchFunc checks if the APIGateway request event can be handled.
type MatchFunc func(events.APIGatewayProxyRequest) bool

// HandleFunc handles a matched APIGateway request event.
type HandleFunc func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
