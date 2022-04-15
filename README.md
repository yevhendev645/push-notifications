# NOTIFICATION

## AWS Service Usage

### SNS (Simple Notification Service)

The project have configured Push Notification from SNS service.

### SQS (Simple Queue Service)

The project has configured 2 queues with different timeouts, one is 30 seconds and one is longer with 10 minutes.

### DynamoDB

The project uses 1 table DynamoDB with 3 primary models.

#### User Details

| Field                   | Description                                    |
| ----------------------- | ---------------------------------------------- |
| ID                      | User ID `(Ex: user-123)`                       |
| Kind                    | Use `details` as specify Kind for User Details |
| ActiveNotification      | User notification status (true or false)       |
| BalanceNotificationTime | User balance notification time                 |
| DailyNotificationTime   | User daily notification time                   |
| DeviceToken             | Device user token                              |

#### User Balance

| Field        | Description              |
| ------------ | ------------------------ |
| ID           | User ID `(Ex: user-123)` |
| Kind         | User balance update time |
| ETHBalance   | ETH balance              |
| BTCBalance   | BTC balance              |
| TotalBalance | Total balance            |

#### Crypto Rate

| Field     | Description             |
| --------- | ----------------------- |
| ID        | Crypto name `(Ex: btc)` |
| Kind      | Crypto Rate update time |
| RateToUSD | Crypto Rate to USD      |

## Filled Order Notification

`POST /api/v1/filled-notification`

This endpoint will notify the filled order to the user.

### Parameters

```json
{
  "userID": "user-123", // userID
  "message": "Your buy order is filled" // notification message
}
```

All fields are required for each destination in the list.

### Output

This returns {} if successful. Returns the errors if not successful.

## Set Notification Status

`POST /api/v1/set-notification-status`

This endpoint will set notification status for user

### Parameters

```json
{
  "userID": "user-123", // userID
  "status": false // notification status
}
```

All fields are required for each destination in the list.

### Output

This returns {} if successful. Returns the errors if not successful.

## User Balance Change Push Notification

We will run cronjob every 1 minute to check user balance change. System will get all user from DynamoDB with `true` in notification status, then will get current user balance and 12 hours ago to calculate. If all the conditions are satisfied, system will push notfication to the user.

## Crypto Rate Change Push Notification

We will run cronjob every 1 minute to check crypto rate change. System will get all user from DynamoDB with `true` in notification status, then will get current crypto rate and few hours ago to calculate. If all the conditions are satisfied, system will push notfication to the user.
