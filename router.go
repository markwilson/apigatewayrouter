// Package apigatewayrouter is a basic router for APIGateway-triggered Lambda
// functions. Each route defined in the router needs to be replicated into the
// APIGateway configuration.
package apigatewayrouter

import (
	"errors"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

// Router stores a map of named routes which will be looped over to find the
// first matching Route.
type Router struct {
	CurrentRouteName string
	Routes           map[string]*Route
	NotFound         HandleFunc
}

// NewRouter creates a new empty Router.
func NewRouter() *Router {
	return &Router{
		"",
		map[string]*Route{},
		nil,
	}
}

// AddRoute puts the defined Route into the Router. There is no clash detection
// for route names, if the same name string is used multiple times then only the
// most recent Route value is used.
//
// This is also used internally by the other `Add*Route` functions.
func (r *Router) AddRoute(name string, route *Route) *Router {
	r.Routes[name] = route

	return r
}

// AddStaticRoute creates a MatchFunc for an exact path match then adds it and
// the handler to the Router using AddRoute.
func (r *Router) AddStaticRoute(name string, method string, uri string, handler HandleFunc) *Router {
	r.AddRoute(name, &Route{
		Match: func(req events.APIGatewayProxyRequest) bool {
			return req.Path == uri && req.HTTPMethod == method
		},
		Handle: handler,
	})

	return r
}

// AddRegExpRoute creates a MatchFunc using a regular expression matcher then
// adds it and the handler to the Router using AddRoute.
func (r *Router) AddRegExpRoute(name string, method string, re *regexp.Regexp, handler HandleFunc) *Router {
	r.AddRoute(name, &Route{
		Match: func(req events.APIGatewayProxyRequest) bool {
			return re.MatchString(req.Path) && req.HTTPMethod == method
		},
		Handle: handler,
	})

	return r
}

// Handle is the routing part of the Router, it is responsible for finding a
// matching Route and executing it. If no matching routes are found, an error is
// triggered.
//
// This function is a valid handler for use in lambda.Start - for example:-
//
//  r := NewRouter()
//  // configure the router's routes...
//  lambda.Start(r.Handle)
func (r *Router) Handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name, route, err := r.firstMatch(req)
	if err != nil {
		if r.NotFound == nil {
			return events.APIGatewayProxyResponse{}, errors.New("Not found")
		}

		return r.NotFound(req)
	}

	r.CurrentRouteName = name

	return route.Handle(req)
}

func (r *Router) firstMatch(req events.APIGatewayProxyRequest) (string, *Route, error) {
	for name, route := range r.Routes {
		if route.Match(req) {
			return name, route, nil
		}
	}

	return "", &Route{}, errors.New("Not found")
}
