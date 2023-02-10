package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/projectdiscovery/retryabledns"
)

const key = "and0x00"

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (string, error) {

	if req.Headers["x"] != key {
		return "", nil
	}

	domain := req.QueryStringParameters["domain"]

	if domain == "" {
		return "", nil
	}

	dns := req.QueryStringParameters["dns"]
	var resolvers []string

	if dns == "" {
		rand.Seed(time.Now().UnixNano())
		resolvers = []string{
			"1.1.1.1:53",
			"1.0.0.1:53",
			"8.8.8.8:53",
			"8.8.4.4:53",
			"9.9.9.9:53",
			"9.9.9.10:53",
			"77.88.8.8:53",
			"77.88.8.1:53",
			"208.67.222.222:53",
			"208.67.220.220:53"}
		rand.Shuffle(len(resolvers), func(i, j int) { resolvers[i], resolvers[j] = resolvers[j], resolvers[i] })
	} else {
		resolvers = []string{fmt.Sprintf("%s:53", dns)}
	}

	dnsClient, err := retryabledns.New(resolvers, 2)
	if err != nil {
		return "", nil
	}

	recordTypes := []uint16{1, 2, 5, 6, 12, 15, 16, 28, 99}
	records, _ := dnsClient.QueryMultiple(domain, recordTypes)

	rjson, err := records.JSON()
	if err != nil {
		return "", nil
	}

	return string(rjson), nil
}

func main() {
	lambda.Start(HandleRequest)
}
