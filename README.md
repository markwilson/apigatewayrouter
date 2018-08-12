# AWS APIGateway Lambda Router

Library for routing APIGateway requests in a Lambda function.

## Example usage

``` go
package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/markwilson/apigatewayrouter"
)

func main() {
	r := apigatewayrouter.NewRouter()

	r.AddRoute("health", &apigatewayrouter.Route{
		Match: func(req events.APIGatewayProxyRequest) bool {
			return req.Path == "/check/health1"
		},
		Handle: func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			log.Println(r.CurrentRouteName)

			return events.APIGatewayProxyResponse{
				Body:       "OK",
				StatusCode: 200,
			}, nil
		},
	})

	r.AddStaticRoute("health2", http.MethodGet, "/check/health2", func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		log.Println(r.CurrentRouteName)

		return events.APIGatewayProxyResponse{
			Body:       "OK",
			StatusCode: 200,
		}, nil
	})

	r.AddRegExpRoute("health3", http.MethodGet, regexp.MustCompile("^/check/health3$"), func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		log.Println(r.CurrentRouteName)

		return events.APIGatewayProxyResponse{
			Body:       "OK",
			StatusCode: 200,
		}, nil
	})

	sr := apigatewayrouter.NewRouter()
	sr.AddStaticRoute("health4", http.MethodGet, "/check/health4", func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		log.Println(r.CurrentRouteName)

		return events.APIGatewayProxyResponse{
			Body:       "OK",
			StatusCode: 200,
		}, nil
	})
	r.AddSubRouter("check", "/check", sr)

	lambda.Start(r.Handle)
}
```
