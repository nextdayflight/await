package main

import (
	"context"
	"fmt"
	"net/url"
)

type resource interface {
	fmt.Stringer
	Await(context.Context) error
}

type unavailabilityError struct {
	Reason error
}

// Error implements the error interface.
func (e *unavailabilityError) Error() string {
	return e.Reason.Error()
}

func parseResources(urlArgs []string) ([]resource, error) {
	var resources []resource
	for _, urlArg := range urlArgs {
		// Leveraging the fact the Go's URL parser matches e.g. `curl -s
		// http://example.com` as url.Path instead of throwing an error.
		u, err := url.Parse(urlArg)
		if err != nil {
			return nil, err
		}
		res, err := identifyResource(*u)
		if err != nil {
			return nil, err
		}
		resources = append(resources, res)
	}
	return resources, nil
}

func identifyResource(u url.URL) (resource, error) {
	switch u.Scheme {
	case "http", "https":
		return &httpResource{u}, nil
	case "ws", "wss":
		return &websocketResource{u}, nil
	case "tcp", "tcp4", "tcp6":
		return &tcpResource{u}, nil
	case "file":
		return &fileResource{u}, nil
	case "postgres":
		return &postgresqlResource{u}, nil
	case "mysql":
		return &mysqlResource{u}, nil
	case "amqp", "amqps":
		return &amqpResource{u}, nil
	case "":
		return &commandResource{u}, nil
	default:
		return nil, fmt.Errorf("unsupported resource scheme: %v", u.Scheme)
	}
}

func parseFragment(fragment string) url.Values {
	// Skip encountered decoding errors on invalid format for now
	v, _ := url.ParseQuery(fragment)
	// Maintain backwards-compatibility for the case that key is empty
	delete(v, "")

	return v
}
