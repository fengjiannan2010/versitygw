package notify

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

// ArcTaskRequest 表示创建任务的请求体
type ArcTaskRequest struct {
	Path       string `json:"path"`
	FerryLevel string `json:"ferryLevel"`
}

// ArcTaskResponse 表示 HTTP 响应的外层结构
type ArcTaskResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"` // 注意：这是嵌套的 JSON 字符串
}

// ArcTaskData 表示嵌套在 data 字符串中的结构
type ArcTaskData struct {
	TaskID     string `json:"taskId"`
	ErrorPaths string `json:"errorPaths"`
	DevID      string `json:"devId"`
}

// ArcTaskNotifier 用于发起 HTTP 请求
type ArcTaskNotifier struct {
	BaseURL      string
	EndpointPath string
	Client       *resty.Client
}

// NewArcTaskNotifier 创建一个新的 Notifier 实例
func NewArcTaskNotifier(baseURL, endpointPath string) *ArcTaskNotifier {
	client := resty.New()
	return &ArcTaskNotifier{
		BaseURL:      baseURL,
		EndpointPath: endpointPath,
		Client:       client,
	}
}

// NotifyFileTransfer 发起任务创建请求并解析响应
func (n *ArcTaskNotifier) NotifyFileTransfer(path, ferryLevel string) (*ArcTaskResponse, *ArcTaskData, error) {
	reqBody := ArcTaskRequest{
		Path:       path,
		FerryLevel: ferryLevel,
	}

	log.Printf("Request Json: %+v\n", reqBody)

	var respBody ArcTaskResponse

	resp, err := n.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&respBody).
		Post(n.BaseURL + n.EndpointPath)

	if err != nil {
		return nil, nil, fmt.Errorf("请求失败: %w", err)
	}

	if resp.IsError() {
		return &respBody, nil, fmt.Errorf("HTTP错误: %s", resp.Status())
	}

	// 解析 data 字段为 ArcTaskData
	var data ArcTaskData
	if err := json.Unmarshal([]byte(respBody.Data), &data); err != nil {
		return &respBody, nil, fmt.Errorf("解析 data 字段失败: %w", err)
	}

	log.Printf("Response Json: %+v\n", data)
	return &respBody, &data, nil
}
