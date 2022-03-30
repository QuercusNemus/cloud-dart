package integration

import (
	"cloud-dart/match"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"os"
)

// DynamoDB struct
type DynamoDB struct {
	tableName *string
	client    *dynamodb.Client
}

// New init function for DynamoDB
func New(config aws.Config) (*DynamoDB, error) {
	tableName, ok := os.LookupEnv("DYNAMODB_TABLE_NAME")

	if !ok {
		return nil, errors.New("DYNAMODB_TABLE_NAME is not set")
	}

	client := dynamodb.NewFromConfig(config)

	return &DynamoDB{
		tableName: &tableName,
		client:    client,
	}, nil
}

// Write a match struct to DynamoDB
func (d *DynamoDB) Write(ctx context.Context, m *match.Match) error {
	marshalMatch, err := attributevalue.MarshalMap(m)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: d.tableName,
		Item:      marshalMatch,
	})

	if err != nil {
		return err
	}
	return nil
}

// Get a match struct from DynamoDB
func (d *DynamoDB) Get(ctx context.Context, matchId string) (match *match.Match, err error) {
	id, err := attributevalue.Marshal(matchId)
	if err != nil {
		return nil, err
	}

	marshalMatch, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: d.tableName,
		Key: map[string]types.AttributeValue{
			"ID": id,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(marshalMatch.Item) == 0 {
		return nil, errors.New("no match found")
	}

	err = attributevalue.UnmarshalMap(marshalMatch.Item, match)
	if err != nil {
		return nil, err
	}

	return match, nil
}
