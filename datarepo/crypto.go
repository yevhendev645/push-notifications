package datarepo

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
)

type CryptoInfo struct {
	ItemKey
	RateToUSD float64 `dynamodbav:"RateToUSD" json:"-"`
}

// GetCryptoRate return crypto rate with id
func GetCryptoRate(svc dynamodbiface.DynamoDBAPI, cryptoName string, kindTime time.Time) (float64, error) {
	var ci CryptoInfo
	ci.ID = cryptoName
	ci.Kind = ItemKind(kindTime.String())
	err := getItem(svc, cryptoName, ItemKind(kindTime.String()), ci)
	return ci.RateToUSD, err
}
