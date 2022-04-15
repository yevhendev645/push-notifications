package datarepo

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
)

type ItemKind string
type PaginationKey map[string]dynamodb.AttributeValue

// ItemKey defines an item key within dynamodb
type ItemKey struct {
	ID   string   `dynamodbav:"ID" json:"-"`
	Kind ItemKind `dynamodbav:"Kind" json:"-"`
}

const (
	KindDetails ItemKind = "details"
)

func buildKey(id string, kind ItemKind) map[string]dynamodb.AttributeValue {
	return map[string]dynamodb.AttributeValue{
		"ID": {
			S: aws.String(id),
		},
		"Kind": {
			S: aws.String(string(kind)),
		},
	}
}

// getItem gets the item from dynamodb and unmarshalls it into the record
func getItem(svc dynamodbiface.DynamoDBAPI, id string, kind ItemKind, record interface{}) error {
	item, err := getItemRaw(svc, id, kind)

	if err != nil {
		return err
	}

	err = dynamodbattribute.UnmarshalMap(item, record)
	if err != nil {
		return err
	}

	return nil
}

// getItemRaw returns the raw map of the data from a dynamodb call
func getItemRaw(svc dynamodbiface.DynamoDBAPI, id string, kind ItemKind) (map[string]dynamodb.AttributeValue, error) {
	input := &dynamodb.GetItemInput{
		Key:       buildKey(id, kind),
		TableName: aws.String(DataTableName()),
	}

	req := svc.GetItemRequest(input)
	result, err := req.Send()
	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, fmt.Errorf("Unable to find item: %s, %s", id, kind)
	}
	return result.Item, nil
}

func updateItem(svc dynamodbiface.DynamoDBAPI, id string, kind ItemKind, updates interface{}) error {
	expr, err := buildSetAttributesUpdate(updates)
	if err != nil {
		return err
	}

	input := dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key:                       buildKey(id, kind),
		TableName:                 aws.String(DataTableName()),
		UpdateExpression:          expr.Update(),
	}

	req := svc.UpdateItemRequest(&input)
	_, err = req.Send()
	return err
}

func buildSetAttributesUpdate(updates interface{}) (expression.Expression, error) {
	vs := reflect.ValueOf(updates)
	update := expression.UpdateBuilder{}

	for i := 0; i < vs.NumField(); i++ {
		tag := vs.Type().Field(i).Tag
		name := vs.Type().Field(i).Name
		if v, ok := tag.Lookup("json"); ok && v != "-" {
			name = v
		}
		if v, ok := tag.Lookup("dynamodbav"); ok && v != "-" {
			name = v
		}
		update = update.Set(
			expression.Name(name),
			expression.Value(vs.Field(i).Interface()),
		)
	}

	return expression.NewBuilder().
		WithUpdate(update).
		Build()
}
