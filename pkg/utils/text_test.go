package utils

import (
	"testing"

	"github.com/yb2020/odoc/pkg/logging"
)

type mockLogger struct {
	logging.Logger
}

func (l *mockLogger) Error(keyvals ...interface{}) error {
	return nil
}

func TestRemoveSpecialWords(t *testing.T) {
	logger := &mockLogger{}
	textUtils := NewTextUtils(logger)

	tests := []struct {
		name     string
		content  string
		list     []string
		expected string
	}{
		{
			name:     "空内容",
			content:  "",
			list:     []string{"\\p{Cntrl}", "\\u0000"},
			expected: "",
		},
		{
			name:     "空列表",
			content:  "Hello\u0000World",
			list:     []string{},
			expected: "Hello\u0000World",
		},
		{
			name:     "替换控制字符",
			content:  "Hello\u0000World\nTest",
			list:     []string{"\\p{Cntrl}"},
			expected: "Hello World Test",
		},
		{
			name:     "替换NULL字符",
			content:  "Hello\u0000World",
			list:     []string{"\\u0000"},
			expected: "Hello World",
		},
		{
			name:     "无效的正则表达式",
			content:  "Hello[World",
			list:     []string{"["},
			expected: "Hello[World",
		},
		{
			name:     "多个正则表达式",
			content:  "Hello\u0000World\nTest",
			list:     []string{"\\p{Cntrl}", "\\u0000"},
			expected: "Hello World Test",
		},
		{
			name:     "多行文本",
			content:  "Hello\nWorld\nTest",
			list:     []string{"\\n"},
			expected: "Hello World Test",
		},
		{
			name:     "连续空格",
			content:  "Hello  \u0000  World",
			list:     []string{"\\u0000", "\\s+"},
			expected: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := textUtils.RemoveSpecialWords(tt.content, tt.list)
			if result != tt.expected {
				t.Errorf("RemoveSpecialWords() = %q, want %q", result, tt.expected)
			}
		})
	}
}
