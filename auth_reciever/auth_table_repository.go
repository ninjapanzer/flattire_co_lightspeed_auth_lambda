package auth_reciever

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
	"time"
)

type TokenEntity struct {
	Provider  string `json:Provider`
	Token     string `json:Token`
	TokenType string `json:TokenType`
	ExpiresIn int16  `json:ExpiresIn`
	CreatedAt int64  `json:CreatedAt`
}

var ddb *dynamodb.DynamoDB

func Save(
	token string,
	tokenType string,
	expiresIn int16,
) {
	region := os.Getenv("AWS_REGION")
	if session, err := session.NewSession(&aws.Config{ // Use aws sdk to connect to dynamoDB
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		ddb = dynamodb.New(session) // Create DynamoDB client
	}

	var tableName = aws.String("authTable")

	tokenEntity := &TokenEntity{
		Provider:  "lightspeed",
		Token:     token,
		TokenType: tokenType,
		ExpiresIn: expiresIn,
		CreatedAt: time.Now().Unix(),
	}

	item, _ := dynamodbattribute.MarshalMap(tokenEntity)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: tableName,
	}
	if _, err := ddb.PutItem(input); err != nil {
		log.Print(err)
	} else {
		log.Print("saved tokens for lightspeed")
		return
	}

	return
}
