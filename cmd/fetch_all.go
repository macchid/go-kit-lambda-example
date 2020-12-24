package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport/awslambda"
	"github.com/macchid/learning/userService/pkg/transport"
	"github.com/macchid/learning/userService/pkg/user"
	"github.com/macchid/learning/userService/pkg/utils"
)

func main() {
	region := os.Getenv("AWS_REGION")

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "Lambda::FetchOne",
			"ts", log.DefaultTimestamp,
			"caller", utils.Caller(3),
		)
	}

	level.Info(logger).Log("msg", "Service started")
	defer level.Info(logger).Log("msg", "Service ended")

	ctx := context.Background()
	svc, err := user.NewBusiness(region, &logger)
	if err != nil {
		level.Error(logger).Log("msg", "Critical Error - Could no initialize the service", "err", err)
		os.Exit(-1)
	}

	options := []awslambda.HandlerOption{
		awslambda.HandlerErrorLogger(logger),
		awslambda.HandlerErrorEncoder(encodeAGProxyError),
	}

	lambda.StartHandlerWithContext(ctx, awslambda.NewHandler(
		transport.MakeLoggerMiddleware(&logger, "Endpoints::FetchAll")(transport.MakeFetchAllEndpoint(svc)),
		decodeAGProxyRequest,
		encodeAGProxyResponse,
		options...,
	))
}

func decodeAGProxyRequest(ctx context.Context, b []byte) (interface{}, error) {
	return nil, nil
}

func encodeAGProxyResponse(ctx context.Context, resp interface{}) ([]byte, error) {
	result := resp.(transport.FetchAllResponse)

	body, err := json.Marshal(result)
	if err != nil {
		return encodeAGProxyError(ctx, err)
	}

	response, err := json.Marshal(events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	})

	return response, err
}

func encodeAGProxyError(ctx context.Context, err error) ([]byte, error) {
	response, err := json.Marshal(events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf("{\"error\":\"%s\"}", err),
	})

	return response, err
}