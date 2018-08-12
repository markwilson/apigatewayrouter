package apigatewayrouter

import (
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/assert"
)

var dummyHandler HandleFunc

func TestMain(m *testing.M) {
	dummyHandler = func(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{}, nil
	}

	os.Exit(m.Run())
}

func TestNewRouterHasNoCurrentRoute(t *testing.T) {
	r := NewRouter()

	assert.Equal(t, "", r.CurrentRouteName, "Expected current route to be empty string")
}

func TestNewRouterHasNoRoutes(t *testing.T) {
	r := NewRouter()

	assert.Empty(t, r.Routes, "Expected empty routes map")
}

func TestAddRouteAddsToRoutesMap(t *testing.T) {
	r := NewRouter()
	r.AddRoute("test", &Route{})

	assert.Contains(t, r.Routes, "test", "Expected test in routes")
}

func TestAddStaticRouteMatchesForValidPath(t *testing.T) {
	r := NewRouter()
	r.AddStaticRoute("test", http.MethodGet, "/test", dummyHandler)

	req := events.APIGatewayProxyRequest{
		Path:       "/test",
		HTTPMethod: http.MethodGet,
	}

	assert.True(t, r.Routes["test"].Match(req), "Expected static route matching to return true")
}

func TestAddStaticRouteDoesNotMatchForInvalidPath(t *testing.T) {
	r := NewRouter()
	r.AddStaticRoute("test", http.MethodGet, "/test", dummyHandler)

	req := events.APIGatewayProxyRequest{
		Path:       "/blah",
		HTTPMethod: http.MethodGet,
	}

	assert.False(t, r.Routes["test"].Match(req), "Expected static route matching to return false")
}

func TestAddStaticRouteDoesNotMatchForInvalidMethod(t *testing.T) {
	r := NewRouter()
	r.AddStaticRoute("test", http.MethodGet, "/test", dummyHandler)

	req := events.APIGatewayProxyRequest{
		Path:       "/test",
		HTTPMethod: http.MethodPost,
	}

	assert.False(t, r.Routes["test"].Match(req), "Expected static route matching to return false")
}

func TestAddRegExpRouteMatchesForValidPath(t *testing.T) {
	re := regexp.MustCompile("^/test$")

	r := NewRouter()
	r.AddRegExpRoute("test", http.MethodGet, re, dummyHandler)

	req := events.APIGatewayProxyRequest{
		Path:       "/test",
		HTTPMethod: http.MethodGet,
	}

	assert.True(t, r.Routes["test"].Match(req), "Expected regexp route matching to return true")
}

func TestAddRegExpRouteDoesNotMatchForInvalidPath(t *testing.T) {
	re := regexp.MustCompile("^/test$")

	r := NewRouter()
	r.AddRegExpRoute("test", http.MethodGet, re, dummyHandler)

	req := events.APIGatewayProxyRequest{
		Path:       "/blah",
		HTTPMethod: http.MethodGet,
	}

	assert.False(t, r.Routes["test"].Match(req), "Expected regexp route matching to return false")
}

func TestAddRegExpRouteDoesNotMatchForInvalidMethod(t *testing.T) {
	re := regexp.MustCompile("^/test$")

	r := NewRouter()
	r.AddRegExpRoute("test", http.MethodGet, re, dummyHandler)

	req := events.APIGatewayProxyRequest{
		Path:       "/test",
		HTTPMethod: http.MethodPost,
	}

	assert.False(t, r.Routes["test"].Match(req), "Expected regexp route matching to return false")
}

func TestHandlerThrowsNotFound(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Path:       "/test",
		HTTPMethod: http.MethodPost,
	}

	r := NewRouter()

	_, err := r.Handle(req)

	assert.NotNil(t, err)
	assert.Equal(t, "Not found", err.Error())
}

func TestHandlerSetsCurrentRouteName(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Path:       "/test",
		HTTPMethod: http.MethodPost,
	}

	expectedResp := events.APIGatewayProxyResponse{
		Body: "Test",
	}

	route := &Route{
		func(_ events.APIGatewayProxyRequest) bool {
			return true
		},
		func(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			return expectedResp, nil
		},
	}

	r := NewRouter()
	r.AddRoute("test", route)

	_, err := r.Handle(req)

	assert.Nil(t, err)
	assert.Equal(t, "test", r.CurrentRouteName)
}

func TestHandlerCallsHandleFunc(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Path:       "/test",
		HTTPMethod: http.MethodPost,
	}

	handlerCalled := false

	route := &Route{
		func(_ events.APIGatewayProxyRequest) bool {
			return true
		},
		func(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			handlerCalled = true

			return events.APIGatewayProxyResponse{}, nil
		},
	}

	r := NewRouter()
	r.AddRoute("test", route)

	_, err := r.Handle(req)

	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestHandlerReturnsCorrectValues(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Path:       "/test",
		HTTPMethod: http.MethodPost,
	}

	expectedResp := events.APIGatewayProxyResponse{
		Body: "Test",
	}

	route := &Route{
		func(_ events.APIGatewayProxyRequest) bool {
			return true
		},
		func(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			return expectedResp, nil
		},
	}

	r := NewRouter()
	r.AddRoute("test", route)

	resp, err := r.Handle(req)

	assert.Nil(t, err)
	assert.Equal(t, expectedResp, resp)
}

func TestFirstMatchGetsARoute(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Path:       "/test",
		HTTPMethod: http.MethodGet,
	}

	expectedRoute := &Route{
		func(_ events.APIGatewayProxyRequest) bool {
			return true
		},
		dummyHandler,
	}

	r := NewRouter()
	r.AddRoute("test", expectedRoute)

	name, route, err := r.firstMatch(req)

	assert.Nil(t, err)
	assert.Equal(t, "test", name)
	assert.Equal(t, expectedRoute, route)
}

func TestFirstMatchGetsNoRoute(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Path:       "/test",
		HTTPMethod: http.MethodGet,
	}

	route := &Route{
		func(_ events.APIGatewayProxyRequest) bool {
			return false
		},
		dummyHandler,
	}

	r := NewRouter()
	r.AddRoute("test", route)

	_, _, err := r.firstMatch(req)

	assert.NotNil(t, err)
}
