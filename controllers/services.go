package controllers

import (
	"notifications/datarepo"
)

type Services interface {
	SetNotificationStatus(userKey datarepo.ItemKey, status bool) error
	GetUserByID(userID string) (datarepo.UserDetails, error)
}

type defaultServices struct{}

// SetNotificationStatus
func (s *defaultServices) SetNotificationStatus(userKey datarepo.ItemKey, status bool) error {
	return datarepo.SetNotificationStatus(datarepo.DynamoDBService(), userKey, status)
}

// GetUserByID
func (s *defaultServices) GetUserByID(userID string) (datarepo.UserDetails, error) {
	return datarepo.GetUserByID(datarepo.DynamoDBService(), userID)
}
