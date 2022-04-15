package datarepo

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
)

// SetNotificationStatus set notification status
func SetNotificationStatus(svc dynamodbiface.DynamoDBAPI, userKey ItemKey, status bool) error {
	u := struct {
		ActiveNotification bool
	}{
		ActiveNotification: status,
	}

	return updateItem(svc, userKey.ID, userKey.Kind, u)
}

// SetBalanceNotificationTime set BalanceNotification time
func SetBalanceNotificationTime(svc dynamodbiface.DynamoDBAPI, userKey ItemKey) error {
	u := struct {
		BalanceNotificationTime time.Time
	}{
		BalanceNotificationTime: time.Now(),
	}

	return updateItem(svc, userKey.ID, userKey.Kind, u)
}

// SetDailyNotificationTime set DailyNotification time
func SetDailyNotificationTime(svc dynamodbiface.DynamoDBAPI, userKey ItemKey) error {
	u := struct {
		DailyNotificationTime time.Time
	}{
		DailyNotificationTime: time.Now(),
	}

	return updateItem(svc, userKey.ID, userKey.Kind, u)
}
