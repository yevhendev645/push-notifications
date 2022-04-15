package notifications

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// Data is the data of the sending pushnotification.
type Data struct {
	Alert *string     `json:"alert,omitempty"`
	Sound *string     `json:"sound,omitempty"`
	Data  interface{} `json:"custom_data"`
	Badge *int        `json:"badge,omitempty"`
}

func snsService() *sns.SNS {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("Unable to load configuration for SNS")
	}

	return sns.New(cfg)
}

// SendTopic sends a push notification to the topic
func SendTopic(data *Data) (err error) {
	svc := snsService()

	if err != nil {
		return
	}

	m, err := newMessageJSON(data)
	if err != nil {
		return
	}

	req := svc.PublishRequest(&sns.PublishInput{
		Message:          aws.String(m),
		MessageStructure: aws.String("json"),
		TopicArn:         aws.String(snsTopic()),
	})

	_, err = req.Send()
	return
}

// Send sends a push notification
func Send(deviceToken string, data *Data) (err error) {
	svc := snsService()

	if err != nil {
		return
	}

	platformEndpointReq := svc.CreatePlatformEndpointRequest(&sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String(snsApplicationARN()),
		Token:                  aws.String(deviceToken),
	})

	platformEndpointResp, err := platformEndpointReq.Send()
	if err != nil {
		return
	}

	m, err := newMessageJSON(data)
	if err != nil {
		return
	}

	req := svc.PublishRequest(&sns.PublishInput{
		Message:          aws.String(m),
		MessageStructure: aws.String("json"),
		TargetArn:        aws.String(*platformEndpointResp.EndpointArn),
	})

	_, err = req.Send()
	return
}

func snsApplicationARN() string {
	return os.Getenv("AWS_SNS_APPLICATION_ARN")
}

func snsTopic() string {
	return os.Getenv("SNS_TOPIC")
}

type message struct {
	APNS        string `json:"APNS"`
	APNSSandbox string `json:"APNS_SANDBOX"`
	Default     string `json:"default"`
	GCM         string `json:"GCM"`
}

type iosPush struct {
	APS Data `json:"aps"`
}

type gcmPush struct {
	Message *string     `json:"message,omitempty"`
	Custom  interface{} `json:"custom"`
	Badge   *int        `json:"badge,omitempty"`
}

type gcmPushWrapper struct {
	Data gcmPush `json:"data"`
}

func newMessageJSON(data *Data) (m string, err error) {
	b, err := json.Marshal(iosPush{
		APS: *data,
	})
	if err != nil {
		return
	}
	payload := string(b)

	b, err = json.Marshal(gcmPushWrapper{
		Data: gcmPush{
			Message: data.Alert,
			Custom:  data.Data,
			Badge:   data.Badge,
		},
	})
	if err != nil {
		return
	}
	gcm := string(b)

	pushData, err := json.Marshal(message{
		Default:     *data.Alert,
		APNS:        payload,
		APNSSandbox: payload,
		GCM:         gcm,
	})
	if err != nil {
		return
	}
	m = string(pushData)
	return
}
