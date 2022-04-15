package tasks

import (
	"encoding/json"
	"fmt"
	"notifications/datarepo"
	"notifications/notifications"
	"notifications/taskmgr"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

const taskCheckBalanceChange taskmgr.TaskType = "SearchCheckBalanceChange"

func init() {
	taskmgr.RegisterTask(taskCheckBalanceChange, &CheckBalanceChange{})
}

// CheckBalanceChange defines a task that check balance user changes
type CheckBalanceChange struct {
	taskmgr.HasRetries
}

func (st *CheckBalanceChange) Process() error {
	activeNotiUserList, err := datarepo.ListAllActiveNotiUsers(datarepo.DynamoDBService())
	if err != nil {
		return err
	}

	for _, userDetails := range activeNotiUserList {
		duration, _ := time.ParseDuration("-12h")
		previousTime := time.Now().Add(duration)

		if userDetails.BalanceNotificationTime.Before(previousTime) {
			currentUserBalance, err := datarepo.GetUserBalance(datarepo.DynamoDBService(), datarepo.ItemKey{ID: userDetails.ID, Kind: datarepo.ItemKind(time.Now().String())})
			if err != nil {
				return err
			}

			previousUserBalance, err := datarepo.GetUserBalance(datarepo.DynamoDBService(), datarepo.ItemKey{ID: userDetails.ID, Kind: datarepo.ItemKind(previousTime.String())})
			if err != nil {
				return err
			}

			precentBalanceChange := (currentUserBalance.TotalBalance - previousUserBalance.TotalBalance/previousUserBalance.TotalBalance) * 100
			if precentBalanceChange > 5 || precentBalanceChange < -5 {
				notifications.Send(userDetails.DeviceToken, &notifications.Data{
					Alert: aws.String(fmt.Sprintf("Your balance %f% change in the last 12 hours", precentBalanceChange)),
					Sound: aws.String("default"),
					Badge: aws.Int(1),
				})

				err := datarepo.SetBalanceNotificationTime(datarepo.DynamoDBService(), datarepo.ItemKey{ID: userDetails.ID, Kind: datarepo.KindDetails})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// TaskType returns the task type for CheckBalanceChange
func (st *CheckBalanceChange) TaskType() taskmgr.TaskType {
	return taskCheckBalanceChange
}

// Unmarshal task from JSON
func (st *CheckBalanceChange) Unmarshal(body []byte) (taskmgr.TaskProcessor, error) {
	result := CheckBalanceChange{}
	err := json.Unmarshal(body, &result)
	return &result, err
}
