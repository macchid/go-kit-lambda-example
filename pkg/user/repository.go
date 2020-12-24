package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/macchid/learning/userService/pkg/utils"
)

type UserData interface {
	FetchAll() ([]User, error)
	FetchOne(key string) (User, error)
	Persist(user User) error
	Delete(key string) error
}

const (
	errorFailedToUnmarshalRecord = "failed to unmarshal record"
	errorFailedToFetchRecord     = "failed to fetch record"
	errorCouldNotMarshalItem     = "could not marshal item"
	errorCouldNotDeleteItem      = "could not delete item"
	errorCouldNotDynamoPutItem   = "could not dynamo put item error"
)

type userRepo struct {
	tableName string
	client    dynamodbiface.DynamoDBAPI
	logger    *log.Logger
}

func NewRepository(dbClient dynamodbiface.DynamoDBAPI, logger *log.Logger) UserData {
	_, logEnd := utils.LogStart(logger, "User::NewRepository")
	defer logEnd(time.Now())

	return &userRepo{
		tableName: "LambdaInGoUser",
		client:    dbClient,
		logger:    logger,
	}
}

func (repo *userRepo) FetchAll() ([]User, error) {
	logger, logEnd := utils.LogStart(repo.logger, "UserRepository::FetchAll")
	defer logEnd(time.Now())

	input := &dynamodb.ScanInput{
		TableName: aws.String(repo.tableName),
	}
	result, err := dynaClient.Scan(input)
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("Unable to perform scan on table %v", repo.tableName), "err", err)
		return nil, errors.New(errorFailedToFetchRecord)
	}

	var items []User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("Unable to unmarshal fetched results from table %v", repo.tableName), "err", err)
		return nil, errors.New(errorFailedToUnmarshalRecord)
	}

	return items, nil
}

func (repo *userRepo) FetchOne(key string) (User, error) {
	logger, logEnd := utils.LogStart(repo.logger, "UserRepository::FetchOne")
	defer logEnd(time.Now())

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(key),
			},
		},
		TableName: aws.String(repo.tableName),
	}

	result, err := repo.client.GetItem(input)
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("Unable to fetch record with key: %v", key), "err", err)
		return User{}, errors.New(errorFailedToFetchRecord)
	}

	var item User
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return User{}, errors.New(errorFailedToUnmarshalRecord)
	}

	return item, nil
}

func (repo *userRepo) Persist(user User) error {
	logger, logEnd := utils.LogStart(repo.logger, "UserRepository::Persist")
	defer logEnd(time.Now())

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		level.Error(logger).Log("msg", "Unable to marshal user information in order to persist", "err", err)
		return errors.New(errorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(repo.tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		level.Error(logger).Log("msg", "Unable to persist the marshalled data", "err", err)
		return errors.New(errorCouldNotDynamoPutItem)
	}

	return nil
}

func (repo *userRepo) Delete(key string) error {
	logger, logEnd := utils.LogStart(repo.logger, "UserRepository::Delete")
	defer logEnd(time.Now())

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(key),
			},
		},
		TableName: aws.String(repo.tableName),
	}
	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("Couldn't delete the record with key %v", key), "err", err)
		return errors.New(errorCouldNotDeleteItem)
	}

	return nil
}
