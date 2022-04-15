/* Copyright, 2018-2019, Reliant Immune Diagnostics, Inc. */

package taskmgr

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

	/* #nosec */
	"math/rand"
	"os"
	"time"

	rollbar "github.com/rollbar/rollbar-go"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

var taskProcessors = map[TaskType]TaskProcessor{}

// TaskOptions defines options that a task can have
type TaskOptions struct {
	Retries     int
	Delay       int64
	LongRunning bool
}

// RegisterTask registers a task to be done.
func RegisterTask(taskType TaskType, processor TaskProcessor) {
	// TODO: Check for already registered?
	taskProcessors[taskType] = processor
}

// RegisterTasks registers a set of tasks
func RegisterTasks(tasks map[TaskType]TaskProcessor) {
	for tt, tp := range tasks {
		RegisterTask(tt, tp)
	}
}

// ProcessTask will process task based on a body
func ProcessTask(taskType string, body string, messageID string) error {
	processorBase := taskProcessors[TaskType(taskType)]
	if processorBase == nil {
		err := fmt.Errorf("unable to find processor for type %s", taskType)
		rollbar.Critical(err)
		return err
	}

	processor, err := processorBase.Unmarshal([]byte(body))
	log.Println("Processing task: ", taskType, body)
	if err != nil {
		return err
	}

	err = processTask(processor)
	return err
}

func processTask(task TaskProcessor) error {
	defer func() {
		if rvr := recover(); rvr != nil {
			rollbar.Critical(rvr)
		}
	}()
	err := task.Process()
	if err != nil {
		reportError(task, err)

		retries := task.GetRetries()
		task.SetAttempts(task.GetAttempts() + 1)
		if retries > 0 {
			err = EnqueueTask(task, TaskOptions{Retries: retries - 1,
				Delay: backoffDelayTime(task.GetAttempts(), 900)})
		}
	}
	return err
}

// backoffDelayTime returns a back off delay time in seconds.
func backoffDelayTime(attempts int, max int) int64 {
	durf := 1.0 * math.Pow(2.0, float64(attempts-1))
	/* #nosec good enough random here */
	durf = rand.Float64()*(durf-1.0) + 1.0
	if math.Round(durf) > float64(max) {
		return int64(max)
	}
	return int64(math.Round(durf))
}

func reportError(task TaskProcessor, err error) {
	params := map[string]interface{}{"taskData": task, "task": task.TaskType()}
	log.Println("Failed task: ", params, err)
	rollbar.ErrorWithExtras(rollbar.ERR, err, params)
}

// EnqueueTask enqueues a task to be performed asynchronously
func EnqueueTask(processor TaskProcessor, options TaskOptions) error {
	stringDataType := "String"
	processor.SetRetries(options.Retries)

	body, err := json.Marshal(processor)
	if err != nil {
		return err
	}

	stringBody := string(body)
	tt := string(processor.TaskType())
	var url string
	if options.LongRunning {
		url = longTaskQueueURL()
	} else {
		url = taskQueueURL()
	}

	att := sqs.MessageAttributeValue{StringValue: &tt, DataType: &stringDataType}

	input := sqs.SendMessageInput{
		DelaySeconds: &options.Delay,
		MessageBody:  &stringBody,
		MessageAttributes: map[string]sqs.MessageAttributeValue{
			"TaskType": att},
		QueueUrl: &url,
	}

	c, err := sqsConnection()
	if err != nil {
		return err
	}

	req := c.SendMessageRequest(&input)
	_, err = req.Send()
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case "InternalError":
			time.Sleep(25 * time.Millisecond) // Small time because this can happen during endpoint handling
			_, err = req.Send()
			return err
		}
	}
	return err
}

var cachedConnection *sqs.SQS

func sqsConnection() (*sqs.SQS, error) {
	if cachedConnection == nil {
		cfg, err := external.LoadDefaultAWSConfig()
		if err != nil {
			return nil, err
		}

		cachedConnection = sqs.New(cfg)
	}
	return cachedConnection, nil
}

func taskQueueURL() string {
	return os.Getenv("TASK_QUEUE_URL")
}

func longTaskQueueURL() string {
	return os.Getenv("LONG_TASK_QUEUE_URL")
}
