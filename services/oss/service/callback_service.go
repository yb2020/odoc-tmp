package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	pb "github.com/yb2020/odoc-proto/gen/go/oss"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/mq/rocketmq"
	"github.com/yb2020/odoc/pkg/mq/rocketmq/producer"
	ossConstant "github.com/yb2020/odoc/services/oss/constant"
	"github.com/yb2020/odoc/services/oss/model"
)

// CallbackService OSS回调服务
type CallbackService struct {
	logger                     logging.Logger
	config                     *config.Config
	uploadNotificationProducer *producer.RocketMQProducer
	ossService                 OssServiceInterface
}

// NewCallbackService 创建回调服务实例
func NewCallbackService(
	logger logging.Logger,
	config *config.Config,
	uploadNotificationProducer *producer.RocketMQProducer,
	ossService OssServiceInterface,
) *CallbackService {
	return &CallbackService{
		logger:                     logger,
		config:                     config,
		uploadNotificationProducer: uploadNotificationProducer,
		ossService:                 ossService,
	}
}

// HandleS3Notification 处理S3通知
func (s *CallbackService) HandleS3Notification(ctx context.Context, snsMessage model.SNSMessage) error {
	s.logger.Info("msg", "Received SNS message", "messageId", snsMessage.MessageId, "type", snsMessage.Type)

	if snsMessage.Type == "SubscriptionConfirmation" {
		s.logger.Info("msg", "Received subscription confirmation. To confirm, visit this URL", "url", snsMessage.SubscribeURL)
		return nil
	}

	var s3Event model.S3Event
	if err := json.Unmarshal([]byte(snsMessage.Message), &s3Event); err != nil {
		s.logger.Error("msg", "Failed to unmarshal S3 event from SNS message", "error", err)
		return err
	}

	if len(s3Event.Records) == 0 {
		return nil
	}

	events := make([]model.OsscallbackLog, 0, len(s3Event.Records))
	for _, record := range s3Event.Records {
		objectKey, err := url.QueryUnescape(record.S3.Object.Key)
		if err != nil {
			s.logger.Warn("msg", "Failed to decode S3 object key", "key", record.S3.Object.Key, "error", err)
			objectKey = record.S3.Object.Key
		}

		event, err := s.processUploadEvent(ctx, record.S3.Bucket.Name, objectKey, record.EventName, record.S3.Object.Size)
		if err != nil {
			s.logger.Warn("msg", "Failed to process S3 event", "objectKey", objectKey, "error", err)
			continue
		}
		if event != nil {
			events = append(events, *event)
		}
	}

	s.sendUploadNotifications(ctx, events)
	return nil
}

// HandleCallback 处理OSS回调通知
func (s *CallbackService) HandleCallback(ctx context.Context, notification *pb.MinioCallbackNotification) error {
	s.logger.Info("msg", "Received MinIO callback notification", "notification", notification)
	if len(notification.Records) == 0 {
		return nil
	}
	events := make([]model.OsscallbackLog, 0, len(notification.Records))
	for _, record := range notification.Records {
		objectKey, err := url.QueryUnescape(record.S3.Object.Key)
		if err != nil {
			s.logger.Warn("msg", "Failed to decode S3 object key", "key", record.S3.Object.Key, "error", err)
			objectKey = record.S3.Object.Key
		}

		event, err := s.processUploadEvent(ctx, record.S3.Bucket.Name, objectKey, record.EventName, record.S3.Object.Size)
		if err != nil {
			s.logger.Warn("msg", "Failed to process MinIO event", "objectKey", objectKey, "error", err)
			continue
		}
		if event != nil {
			events = append(events, *event)
		}
	}
	s.sendUploadNotifications(ctx, events)
	return nil
}

// processUploadEvent is the unified handler for processing file upload events from any source.
func (s *CallbackService) processUploadEvent(ctx context.Context, bucketName, objectKey, eventName string, size int64) (*model.OsscallbackLog, error) {
	ossRecord, err := s.ossService.GetRecordByBucketNameAndObjectKey(ctx, bucketName, objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get record by object key '%s': %w", objectKey, err)
	}
	if ossRecord == nil {
		return nil, fmt.Errorf("no record found for object key '%s'", objectKey)
	}

	if ossRecord.Status == ossConstant.FileStatusSuccess {
		s.logger.Info("msg", "Upload already processed, skipping", "objectKey", objectKey, "recordId", ossRecord.Id)
		return nil, nil
	}

	event := &model.OsscallbackLog{
		BucketName:   bucketName,
		ObjectKey:    objectKey,
		EventType:    eventName,
		Size:         size,
		RecordId:     ossRecord.Id,
		UploadUserId: ossRecord.CreatorId,
		TopicName:    ossRecord.CallbackTopic,
		UserMetadata: ossRecord.BizMetadata,
	}

	if err := s.ossService.UpdateFileStatus(ctx, ossRecord.Id, ossConstant.FileStatusSuccess, event.Size); err != nil {
		return nil, fmt.Errorf("failed to update file status for record '%d': %w", ossRecord.Id, err)
	}

	return event, nil
}

// sendUploadNotifications 发送上传通知消息
func (s *CallbackService) sendUploadNotifications(ctx context.Context, events []model.OsscallbackLog) {
	s.logger.Info("msg", "Sending upload notifications", "count", len(events))
	if s.uploadNotificationProducer == nil || len(events) == 0 {
		return
	}

	for _, event := range events {
		s.logger.Info("msg", "Sending upload notification", "event", event)
		if !(event.EventType == "ObjectCreated:Put" || event.EventType == "ObjectCreated:Post" ||
			event.EventType == "s3:ObjectCreated:Put" || event.EventType == "s3:ObjectCreated:Post") {
			continue
		}

		if event.TopicName == "" {
			continue
		}

		data, err := json.Marshal(event)
		if err != nil {
			s.logger.Error("序列化事件数据失败", "error", err.Error())
			continue
		}

		userId := event.UploadUserId
		// 创建消息接口实例
		message := &rocketmq.Message{
			Topic:  event.TopicName,
			Keys:   event.RecordId,
			UserId: userId,
			Body:   data,
		}
		s.logger.Info("msg", "Sending upload notification", "message", message)
		_, err = s.uploadNotificationProducer.SendSync(ctx, message)
		if err != nil {
			s.logger.Error("发送上传通知消息失败", "error", err.Error(), "topic", event.TopicName)
		}
	}
}
