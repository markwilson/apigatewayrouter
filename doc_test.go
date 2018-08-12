package apigatewayrouter

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Example() {
	// Initialise a new Router.
	r := NewRouter()

	// Add a single Route to it.
	r.AddRoute("test", &Route{
		Match: func(req events.APIGatewayProxyRequest) bool {
			return req.Path == "/test"
		},
		Handle: func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			fmt.Println("Current route:", r.CurrentRouteName)

			return events.APIGatewayProxyResponse{
				Body:       "OK",
				StatusCode: 200,
			}, nil
		},
	})

	// Start handling request events.
	lambda.Start(r.Handle)
}
