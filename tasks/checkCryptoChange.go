package tasks

import (
	"encoding/json"
	"fmt"
	"notifications/datarepo"
	"notifications/notifications"
	"notifications/taskmgr"
	"notifications/utils"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

const taskCheckCryptoChange taskmgr.TaskType = "SearchCheckCryptoChange"

func init() {
	taskmgr.RegisterTask(taskCheckCryptoChange, &CheckCryptoChange{})
}

// CheckCryptoChange defines a task that check crypto rate changes
type CheckCryptoChange struct {
	taskmgr.HasRetries
}

func (st *CheckCryptoChange) Process() error {
	cryptoList := []string{"btc", "eth"}
	durationList := []string{"-1h", "-2h", "-3h"}

	activeNotiUserList, err := datarepo.ListAllActiveNotiUsers(datarepo.DynamoDBService())
	if err != nil {
		return err
	}

	for _, cryptoName := range cryptoList {
		currentCryptoRate, err := datarepo.GetCryptoRate(datarepo.DynamoDBService(), cryptoName, time.Now())
		if err != nil {
			return err
		}

		for _, durationStr := range durationList {
			duration, _ := time.ParseDuration(durationStr)
			previousTime := time.Now().Add(duration)
			dayDuration, _ := time.ParseDuration("-24h")
			previousDayTime := time.Now().Add(dayDuration)

			previousCryptoRate, err := datarepo.GetCryptoRate(datarepo.DynamoDBService(), cryptoName, previousTime)
			if err != nil {
				return err
			}

			precentCryptoChange := (currentCryptoRate - previousCryptoRate/previousCryptoRate) * 100
			if precentCryptoChange > 5 {
				for _, userDetails := range activeNotiUserList {
					if userDetails.DailyNotificationTime.Before(previousDayTime) {
						notifications.Send(userDetails.DeviceToken, &notifications.Data{
							Alert: aws.String(fmt.Sprintf("%s up %f% in last %s", strings.ToUpper(cryptoName), precentCryptoChange, utils.ParseDuration(durationStr))),
							Sound: aws.String("default"),
							Badge: aws.Int(1),
						})

						err := datarepo.SetDailyNotificationTime(datarepo.DynamoDBService(), datarepo.ItemKey{ID: userDetails.ID, Kind: datarepo.KindDetails})
						if err != nil {
							return err
						}
					}
				}
				return nil
			} else if precentCryptoChange < -5 {
				for _, userDetails := range activeNotiUserList {
					if userDetails.DailyNotificationTime.Before(previousDayTime) {
						notifications.Send(userDetails.DeviceToken, &notifications.Data{
							Alert: aws.String(fmt.Sprintf("%s down %f% in last %s", strings.ToUpper(cryptoName), precentCryptoChange, utils.ParseDuration(durationStr))),
							Sound: aws.String("default"),
							Badge: aws.Int(1),
						})

						err := datarepo.SetDailyNotificationTime(datarepo.DynamoDBService(), datarepo.ItemKey{ID: userDetails.ID, Kind: datarepo.KindDetails})
						if err != nil {
							return err
						}
					}
				}
				return nil
			}
		}
	}

	return nil
}

// TaskType returns the task type for CheckCryptoChange
func (st *CheckCryptoChange) TaskType() taskmgr.TaskType {
	return taskCheckCryptoChange
}

// Unmarshal task from JSON
func (st *CheckCryptoChange) Unmarshal(body []byte) (taskmgr.TaskProcessor, error) {
	result := CheckCryptoChange{}
	err := json.Unmarshal(body, &result)
	return &result, err
}
