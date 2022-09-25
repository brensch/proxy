package client

import (
	"context"
	"fmt"

	"google.golang.org/api/run/v1"
)

var (
	// this is the name of the handlerfunc given to the entrypoint of the cloud build
	HandlerName = "proxyrequest"
)

// AuditProxies retrieves all proxies available in the current project.
func AuditProxies(projectID string) ([]string, error) {
	ctx := context.Background()

	c, err := run.NewService(ctx)
	if err != nil {
		return nil, err
	}

	list := c.Projects.Locations.Services.List(fmt.Sprintf("projects/%s/locations/-", projectID))

	it, err := list.Do()
	if err != nil {
		return nil, err
	}

	var proxies []string
	for _, service := range it.Items {
		if service.Metadata.Name == HandlerName {
			proxies = append(proxies, service.Status.Address.Url)
		}
	}

	return proxies, nil
}
