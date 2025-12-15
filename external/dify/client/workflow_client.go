package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/external/dify/callback"
	"github.com/yb2020/odoc/external/dify/constant"
	"github.com/yb2020/odoc/external/dify/proto"
	"github.com/yb2020/odoc/external/dify/transport"
	"github.com/yb2020/odoc/pkg/logging"
)

// DifyWorkflowClient Dify工作流客户端接口
type DifyWorkflowClient interface {
	DifyBaseClient

	// RunWorkflow 运行工作流
	RunWorkflow(ctx context.Context, request *proto.WorkflowRunRequest) (*proto.WorkflowRunResponse, error)

	// RunStreamWorkflow 运行流式工作流
	RunStreamWorkflow(ctx context.Context, request *proto.WorkflowRunRequest, callback callback.WorkflowStreamCallback) error

	// StopWorkflow 停止工作流
	StopWorkflow(ctx context.Context, taskID string, user string) (*proto.WorkflowStopResponse, error)

	// GetWorkflowRun 获取工作流执行情况
	GetWorkflowRun(ctx context.Context, workflowID string) (*proto.WorkflowRunStatusResponse, error)

	// GetWorkflowLogs 获取工作流日志
	GetWorkflowLogs(ctx context.Context, keyword string, status string, page int, limit int) (*proto.WorkflowLogsResponse, error)
}

// DifyWorkflowClientImpl Dify工作流客户端实现
type DifyWorkflowClientImpl struct {
	DifyBaseClientImpl
	logger logging.Logger
	tracer opentracing.Tracer
}

// NewDifyWorkflowClient 创建新的Dify工作流客户端
func NewDifyWorkflowClient(logger logging.Logger, tracer opentracing.Tracer, apiBaseUrl, apiKey string, timeout int, responseTimeout int) DifyWorkflowClient {
	return &DifyWorkflowClientImpl{
		DifyBaseClientImpl: DifyBaseClientImpl{
			apiKey:     apiKey,
			apiBaseUrl: apiBaseUrl,
			httpClient: transport.NewDefaultHttpClient(logger, timeout, responseTimeout),
			tracer:     tracer,
		},
		logger: logger,
		tracer: tracer,
	}
}

// NewDifyWorkflowClientWithHttpClient 使用自定义HTTP客户端创建新的Dify工作流客户端
func NewDifyWorkflowClientWithHttpClient(logger logging.Logger, tracer opentracing.Tracer, apiBaseUrl, apiKey string, httpClient transport.HttpClient) DifyWorkflowClient {
	return &DifyWorkflowClientImpl{
		DifyBaseClientImpl: DifyBaseClientImpl{
			apiKey:     apiKey,
			apiBaseUrl: apiBaseUrl,
			httpClient: httpClient,
			tracer:     tracer,
		},
		logger: logger,
		tracer: tracer,
	}
}

// RunWorkflow 运行工作流
func (c *DifyWorkflowClientImpl) RunWorkflow(ctx context.Context, request *proto.WorkflowRunRequest) (*proto.WorkflowRunResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyWorkflowClientImpl.RunWorkflow")
	defer span.Finish()

	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 创建响应对象
	response := &proto.WorkflowRunResponse{}

	// 发送请求并解析响应
	err := c.doJSONRequest(ctx, "POST", c.buildURL(constant.WORKFLOWS_RUN_PATH), request, response)
	if err != nil {
		return nil, fmt.Errorf("failed to run workflow: %w", err)
	}

	return response, nil
}

// RunStreamWorkflow 运行流式工作流
func (c *DifyWorkflowClientImpl) RunStreamWorkflow(ctx context.Context, request *proto.WorkflowRunRequest, callback callback.WorkflowStreamCallback) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyWorkflowClientImpl.RunStreamWorkflow")
	defer span.Finish()

	if request == nil {
		return fmt.Errorf("请求不能为空")
	}
	if callback == nil {
		return fmt.Errorf("回调不能为空")
	}

	// 设置流式响应模式
	request.ResponseMode = "streaming"

	// 构建URL
	url := transport.BuildUrl(c.apiBaseUrl, constant.WORKFLOWS_RUN_PATH)

	// 将请求转换为JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// 设置请求头
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + c.apiKey
	headers["Content-Type"] = transport.ContentTypeJson
	headers["Accept"] = "text/event-stream"

	// 发送请求
	resp, err := c.httpClient.Post(ctx, url, bytes.NewBuffer(jsonData), transport.ContentTypeJson, headers)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		errorCode, errorMessage, _ := transport.ExtractErrorInfo(string(respBody))
		return fmt.Errorf("HTTP错误 %d: %s (%s)", resp.StatusCode, errorMessage, errorCode)
	}

	// 处理流式响应
	go func() {
		defer resp.Body.Close()
		reader := bufio.NewReader(resp.Body)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					callback.OnError(err.Error())
				}
				break
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			// 提取数据部分
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				callback.OnEnd()
				break
			}

			// 解析事件类型和数据
			var event proto.BaseEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				callback.OnError(err.Error())
				continue
			}

			switch event.Event {
			case proto.EventMessage:
				callback.OnMessage(event.Data)
			case proto.EventError:
				callback.OnError(event.Data)
			case "workflow_start":
				var startData struct {
					TaskID string `json:"task_id"`
				}
				if err := json.Unmarshal([]byte(event.Data), &startData); err == nil {
					callback.OnStart(startData.TaskID)
				}
			case "workflow_finish":
				var finishData struct {
					Outputs map[string]interface{} `json:"outputs"`
				}
				if err := json.Unmarshal([]byte(event.Data), &finishData); err == nil {
					callback.OnFinish(finishData.Outputs)
				}
			case "node_start":
				var nodeStartData struct {
					NodeID   string `json:"node_id"`
					NodeName string `json:"node_name"`
					NodeType string `json:"node_type"`
				}
				if err := json.Unmarshal([]byte(event.Data), &nodeStartData); err == nil {
					callback.OnNodeStart(nodeStartData.NodeID, nodeStartData.NodeName, nodeStartData.NodeType)
				}
			case "node_finish":
				var nodeFinishData struct {
					NodeID   string                 `json:"node_id"`
					NodeName string                 `json:"node_name"`
					NodeType string                 `json:"node_type"`
					Outputs  map[string]interface{} `json:"outputs"`
				}
				if err := json.Unmarshal([]byte(event.Data), &nodeFinishData); err == nil {
					callback.OnNodeFinish(nodeFinishData.NodeID, nodeFinishData.NodeName, nodeFinishData.NodeType, nodeFinishData.Outputs)
				}
			case "node_error":
				var nodeErrorData struct {
					NodeID   string `json:"node_id"`
					NodeName string `json:"node_name"`
					NodeType string `json:"node_type"`
					Error    string `json:"error"`
				}
				if err := json.Unmarshal([]byte(event.Data), &nodeErrorData); err == nil {
					callback.OnNodeError(nodeErrorData.NodeID, nodeErrorData.NodeName, nodeErrorData.NodeType, nodeErrorData.Error)
				}
			}
		}
	}()

	return nil
}

// StopWorkflow 停止工作流
func (c *DifyWorkflowClientImpl) StopWorkflow(ctx context.Context, taskID string, user string) (*proto.WorkflowStopResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyWorkflowClientImpl.StopWorkflow")
	defer span.Finish()

	if taskID == "" {
		return nil, fmt.Errorf("任务ID不能为空")
	}

	// 准备请求体
	requestBody := map[string]string{}
	if user != "" {
		requestBody["user"] = user
	}

	// 创建响应对象
	response := &proto.WorkflowStopResponse{}

	// 发送请求并解析响应
	path := fmt.Sprintf("/workflows/run/%s/stop", taskID)
	err := c.doJSONRequest(ctx, "POST", c.buildURL(path), requestBody, response)
	if err != nil {
		return nil, fmt.Errorf("failed to stop workflow: %w", err)
	}

	return response, nil
}

// GetWorkflowRun 获取工作流执行情况
func (c *DifyWorkflowClientImpl) GetWorkflowRun(ctx context.Context, workflowID string) (*proto.WorkflowRunStatusResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyWorkflowClientImpl.GetWorkflowRun")
	defer span.Finish()

	if workflowID == "" {
		return nil, fmt.Errorf("工作流ID不能为空")
	}

	// 创建响应对象
	response := &proto.WorkflowRunStatusResponse{}

	// 发送请求并解析响应
	path := fmt.Sprintf("/workflows/run/%s", workflowID)
	err := c.doJSONRequest(ctx, "GET", c.buildURL(path), nil, response)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow run status: %w", err)
	}

	return response, nil
}

// GetWorkflowLogs 获取工作流日志
func (c *DifyWorkflowClientImpl) GetWorkflowLogs(ctx context.Context, keyword string, status string, page int, limit int) (*proto.WorkflowLogsResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyWorkflowClientImpl.GetWorkflowLogs")
	defer span.Finish()

	// 准备查询参数
	queryParams := map[string]string{}
	if keyword != "" {
		queryParams["keyword"] = keyword
	}
	if status != "" {
		queryParams["status"] = status
	}
	if page > 0 {
		queryParams["page"] = strconv.Itoa(page)
	}
	if limit > 0 {
		queryParams["limit"] = strconv.Itoa(limit)
	}

	// 创建响应对象
	response := &proto.WorkflowLogsResponse{}

	// 发送请求并解析响应
	err := c.getWithParams(ctx, "/workflows/logs", queryParams, response)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow logs: %w", err)
	}

	return response, nil
}
