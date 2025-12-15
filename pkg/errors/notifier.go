package errors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yb2020/odoc/pkg/logging"
)

// SlackNotifier sends error notifications to Slack
type SlackNotifier struct {
	webhookURL string
	logger     logging.Logger
	appName    string
}

// NewSlackNotifier creates a new Slack notifier
func NewSlackNotifier(webhookURL, appName string, logger logging.Logger) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		logger:     logger,
		appName:    appName,
	}
}

// Notify sends an error notification to Slack
func (n *SlackNotifier) Notify(ctx context.Context, err error) error {
	if n.webhookURL == "" {
		return fmt.Errorf("slack webhook URL is not configured")
	}

	var appErr *SystemError
	errorMessage := err.Error()
	errorType := "UNKNOWN"
	stack := ""

	if As(err, &appErr) {
		errorType = string(appErr.Type)
		stack = appErr.Stack
	}

	// Create Slack message payload
	payload := map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"fallback":    fmt.Sprintf("[%s] Error: %s", n.appName, errorMessage),
				"color":       "#FF0000",
				"title":       fmt.Sprintf("[%s] Error Notification", n.appName),
				"title_link":  "",
				"text":        errorMessage,
				"footer":      "Error Notification System",
				"footer_icon": "",
				"ts":          time.Now().Unix(),
				"fields": []map[string]interface{}{
					{
						"title": "Error Type",
						"value": errorType,
						"short": true,
					},
					{
						"title": "Time",
						"value": time.Now().Format(time.RFC3339),
						"short": true,
					},
				},
			},
		},
	}

	if stack != "" {
		payload["attachments"].([]map[string]interface{})[0]["fields"] = append(
			payload["attachments"].([]map[string]interface{})[0]["fields"].([]map[string]interface{}),
			map[string]interface{}{
				"title": "Stack Trace",
				"value": fmt.Sprintf("```%s```", stack),
				"short": false,
			},
		)
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %w", err)
	}

	// Send HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", n.webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack notification failed with status code: %d", resp.StatusCode)
	}

	n.logger.Info("msg", "Error notification sent to Slack", "error", errorMessage)
	return nil
}

// EmailNotifier sends error notifications via email
type EmailNotifier struct {
	smtpServer string
	smtpPort   int
	username   string
	password   string
	from       string
	to         []string
	logger     logging.Logger
	appName    string
}

// NewEmailNotifier creates a new email notifier
func NewEmailNotifier(smtpServer string, smtpPort int, username, password, from string, to []string, appName string, logger logging.Logger) *EmailNotifier {
	return &EmailNotifier{
		smtpServer: smtpServer,
		smtpPort:   smtpPort,
		username:   username,
		password:   password,
		from:       from,
		to:         to,
		logger:     logger,
		appName:    appName,
	}
}

// Notify sends an error notification via email
// Note: This is a placeholder implementation. In a real application, you would use a proper email library.
func (n *EmailNotifier) Notify(ctx context.Context, err error) error {
	// In a real implementation, you would use a library like "net/smtp" or a third-party package
	// to send emails. This is just a placeholder that logs the notification.
	n.logger.Info(
		"msg", "Would send email notification",
		"error", err.Error(),
		"from", n.from,
		"to", fmt.Sprintf("%v", n.to),
		"subject", fmt.Sprintf("[%s] Error Notification", n.appName),
	)
	return nil
}
