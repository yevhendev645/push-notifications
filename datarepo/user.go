package datarepo

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
)

type UserDetails struct {
	ItemKey
	ActiveNotification      bool      `dynamodbav:"ActiveNotification" json:"-"`
	BalanceNotificationTime time.Time `dynamodbav:"BalanceNotificationTime" json:"-"`
	DailyNotificationTime   time.Time `dynamodbav:"DailyNotificationTime" json:"-"`
	DeviceToken             string    `dynamodbav:"DeviceToken" json:"-"`
}

type UserDetailsPage struct {
	Items            []UserDetails
	LastEvaluatedKey PaginationKey
	TotalCount       int64
}

// GetUserByID returns a user with known id
func GetUserByID(svc dynamodbiface.DynamoDBAPI, id string) (details UserDetails, err error) {
	details.ID = id
	details.Kind = KindDetails
	err = getItem(svc, details.ID, details.Kind, &details)
	return details, err
}

// ListAllActiveNotiUsers lists all active notification user
func ListAllActiveNotiUsers(svc dynamodbiface.DynamoDBAPI) ([]UserDetails, error) {
	var page PaginationKey
	first := true
	var result []UserDetails
	for first || page != nil {
		first = false
		userDetails, err := ListActiveNotiUsers(svc, page)
		if err != nil {
			return nil, err
		}
		result = append(result, userDetails.Items...)
		page = userDetails.LastEvaluatedKey
	}
	return result, nil
}

// ListActiveNotiUsers lists the active notification user
func ListActiveNotiUsers(svc dynamodbiface.DynamoDBAPI, page PaginationKey) (UserDetailsPage, error) {
	var err error
	var result *dynamodb.QueryOutput

	result, err = queryUserItems(svc, page)
	if err != nil {
		return UserDetailsPage{}, err
	}

	resultPage := UserDetailsPage{
		LastEvaluatedKey: result.LastEvaluatedKey,
		TotalCount:       *result.Count,
	}

	if *result.Count > 0 {
		resultPage.Items = make([]UserDetails, *result.Count)
		for index, item := range (*result).Items {
			dynamodbattribute.UnmarshalMap(item, &resultPage.Items[index])
		}
	}

	return resultPage, nil
}

func queryUserItems(svc dynamodbiface.DynamoDBAPI, page PaginationKey) (*dynamodb.QueryOutput, error) {
	var expr expression.Expression
	var err error

	prov := expression.Equal(expression.Name("Kind"), expression.Value(KindDetails))
	prov = prov.And(expression.Equal(expression.Name("ActiveNotification"), expression.Value(true)))
	expr, err = expression.NewBuilder().WithFilter(prov).Build()

	if err != nil {
		panic("error in group expression retrieving group items")
	}

	input := &dynamodb.QueryInput{
		Select:                    "ALL_ATTRIBUTES",
		TableName:                 aws.String(DataTableName()),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	if page != nil {
		input.ExclusiveStartKey = page
	}

	return svc.QueryRequest(input).Send()
}
