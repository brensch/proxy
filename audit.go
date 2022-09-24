package proxy

import (
	"context"
	"fmt"

	functions "cloud.google.com/go/functions/apiv1"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	functionspb "google.golang.org/genproto/googleapis/cloud/functions/v1"
)

var (
	// this is the name of the handlerfunc given to the entrypoint of the cloud build
	HandlerName = "HandleProxy"
)

// AuditProxies retrieves all proxies available in the current project.
// NB currently gen2 functions do not get returned by the ListFunctions method.
// Next time i look at this repo hopefully google has made their API cater
// to my diverse requirements.
func AuditProxies(projectID string) ([]string, error) {
	ctx := context.Background()
	c, err := functions.NewCloudFunctionsClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()

	req := &functionspb.ListFunctionsRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/-", projectID),
		PageSize: 100,
	}

	var uris []string
	it := c.ListFunctions(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Error("failed to get next function", zap.Error(err))
			return nil, err
		}
		fmt.Println(resp.EntryPoint)
		fmt.Println(resp.EntryPoint == HandlerName)
		if resp.EntryPoint == HandlerName {
			resp.GetHttpsTrigger()
			uris = append(uris, resp.GetHttpsTrigger().Url)
		}
	}

	return uris, nil
}
