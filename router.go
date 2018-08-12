package apigatewayrouter

import (
	"errors"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

// Router handles custom APIGateway request routing
type Router struct {
	CurrentRouteName string
	Routes           map[string]*Route
}

// NewRouter creates a new Router
func NewRouter() *Router {
	return &Router{
		"",
		map[string]*Route{},
	}
}

// AddRoute adds a Route to a Router
func (r *Router) AddRoute(name string, route *Route) *Router {
	r.Routes[name] = route

	return r
}

// AddStaticRoute adds a static URI matcher and handler to a Router
func (r *Router) AddStaticRoute(name string, method string, uri string, handler HandleFunc) *Router {
	r.Routes[name] = &Route{
		Match: func(req events.APIGatewayProxyRequest) bool {
			return req.Path == uri && req.HTTPMethod == method
		},
		Handle: handler,
	}

	return r
}

// AddRegExpRoute adds a regular expression matcher and handler to a Router
func (r *Router) AddRegExpRoute(name string, method string, re *regexp.Regexp, handler HandleFunc) *Router {
	r.Routes[name] = &Route{
		Match: func(req events.APIGatewayProxyRequest) bool {
			return re.MatchString(req.Path) && req.HTTPMethod == method
		},
		Handle: handler,
	}

	return r
}

// AddSubRouter adds a path-prefixed Router to an existing Router
func (r *Router) AddSubRouter(name string, prefix string, subRouter *Router) *Router {
	re := regexp.MustCompile("^" + prefix)

	r.Routes[name] = &Route{
		Match: func(req events.APIGatewayProxyRequest) bool {
			if !re.MatchString(req.Path) {
				return false
			}

			_, _, err := subRouter.firstMatch(req)

			return err == nil
		},
		Handle: subRouter.Handle,
	}

	return r
}

// Handle handles an APIGateway request event
func (r *Router) Handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name, route, err := r.firstMatch(req)
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.New("Not found")
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
