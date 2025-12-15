package model

import "time"

// SNSMessage represents the overall structure of a message from SNS.
type SNSMessage struct {
	Type             string    `json:"Type"`
	MessageId        string    `json:"MessageId"`
	TopicArn         string    `json:"TopicArn"`
	Message          string    `json:"Message"`
	Timestamp        time.Time `json:"Timestamp"`
	SignatureVersion string    `json:"SignatureVersion"`
	Signature        string    `json:"Signature"`
	SigningCertURL   string    `json:"SigningCertURL"`
	SubscribeURL     string    `json:"SubscribeURL,omitempty"`
	UnsubscribeURL   string    `json:"UnsubscribeURL,omitempty"`
}

// S3Event represents the S3 event notification structure embedded in the SNS message.
type S3Event struct {
	Records []S3EventRecord `json:"Records"`
}

// S3EventRecord contains the details of the S3 event.
type S3EventRecord struct {
	EventVersion string    `json:"eventVersion"`
	EventSource  string    `json:"eventSource"`
	AwsRegion    string    `json:"awsRegion"`
	EventTime    time.Time `json:"eventTime"`
	EventName    string    `json:"eventName"`
	S3           S3Data    `json:"s3"`
}

// S3Data contains the bucket and object details.
type S3Data struct {
	Bucket S3Bucket `json:"bucket"`
	Object S3Object `json:"object"`
}

// S3Bucket contains the bucket name.
type S3Bucket struct {
	Name string `json:"name"`
}

// S3Object contains the object key.
type S3Object struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}
