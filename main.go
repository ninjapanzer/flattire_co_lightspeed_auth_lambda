package main

import (
	"errors"
	"lightspeed_auth_lambda/auth_reciever"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no code was provided")
)

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	// If no name is provided in the HTTP request body, throw an error
	if len(request.QueryStringParameters["code"]) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	authResp := auth_reciever.ExchangeCode(request.QueryStringParameters["code"])

	return events.APIGatewayProxyResponse{
		Body:       "Hello token: " + authResp.AccessToken + " expires in: " + strconv.FormatInt(authResp.ExpiresIn, 10),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
