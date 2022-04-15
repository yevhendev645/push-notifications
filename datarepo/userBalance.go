package datarepo

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
)

type UserBalance struct {
	ItemKey
	ETHBalance   float64 `dynamodbav:"ETH" json:"-"`
	BTCBalance   float64 `dynamodbav:"BTC" json:"-"`
	TotalBalance float64 `dynamodbav:"TotalBalance" json:"-"`
}

type UserBalancePage struct {
	Items            []UserBalance
	LastEvaluatedKey PaginationKey
	TotalCount       int64
}

// ListAllUserBalances lists all user balance with specify time
func ListAllUserBalances(svc dynamodbiface.DynamoDBAPI, kindTime time.Time) ([]UserBalance, error) {
	var page PaginationKey
	first := true
	var result []UserBalance
	for first || page != nil {
		first = false
		userBalances, err := ListUserBalances(svc, page, kindTime)
		if err != nil {
			return nil, err
		}
		result = append(result, userBalances.Items...)
		page = userBalances.LastEvaluatedKey
	}
	return result, nil
}

// ListUserBalances lists the user balance with specify time
func ListUserBalances(svc dynamodbiface.DynamoDBAPI, page PaginationKey, kindTime time.Time) (UserBalancePage, error) {
	var err error
	var result *dynamodb.QueryOutput

	result, err = queryBalanceItems(svc, page, kindTime)
	if err != nil {
		return UserBalancePage{}, err
	}

	resultPage := UserBalancePage{
		LastEvaluatedKey: result.LastEvaluatedKey,
		TotalCount:       *result.Count,
	}

	if *result.Count > 0 {
		resultPage.Items = make([]UserBalance, *result.Count)
		for index, item := range (*result).Items {
			dynamodbattribute.UnmarshalMap(item, &resultPage.Items[index])
		}
	}

	return resultPage, nil
}

// GetUserBalance gets user balance by an item key
func GetUserBalance(svc dynamodbiface.DynamoDBAPI, key ItemKey) (ub UserBalance, err error) {
	ub.ItemKey = key
	err = getItem(svc, key.ID, key.Kind, &ub)
	return ub, err
}

func queryBalanceItems(svc dynamodbiface.DynamoDBAPI, page PaginationKey, kindTime time.Time) (*dynamodb.QueryOutput, error) {
	var expr expression.Expression
	var err error

	prov := expression.BeginsWith(expression.Name("ID"), "user-")
	prov = prov.And(expression.Equal(expression.Name("Kind"), expression.Value(kindTime.UTC().Format(time.RFC3339))))
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
