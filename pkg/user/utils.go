package user

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/macchid/learning/userService/pkg/utils"
)

var dynaClient dynamodbiface.DynamoDBAPI

func NewBusiness(region string, parent *log.Logger) (UserService, error) {
	logger, logEnd := utils.LogStart(parent, "User::NewBusiness")
	defer logEnd(time.Now())

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		level.Error(logger).Log("msg", "Unable to create a new AWS Session", "err", err)
		return nil, err
	}

	dynaClient = dynamodb.New(awsSession)

	userRepo := NewRepository(dynaClient, parent)
	userSvc := NewService(userRepo, parent)

	return userSvc, nil
}
