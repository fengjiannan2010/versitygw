package notify

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
)

// ArcTaskRequest 表示创建任务的请求体
type ArcTaskRequest struct {
	Path       string `json:"path"`
	FerryLevel string `json:"ferryLevel"`
	LoginName  string `json:"loginName"`
}

// ArcTaskResponse 表示 HTTP 响应的外层结构
type ArcTaskResponse struct {
	Code  string `json:"code"`
	State string `json:"state"`
	Msg   string `json:"msg"`
	Data  string `json:"data,omitempty"` // 如果 `Data` 不可用，可以忽略
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

// NewArcTaskNotifier 根据 URL 判断是否启用 TLS
func NewArcTaskNotifier(baseURL, endpointPath string) *ArcTaskNotifier {
	client := resty.New()
	if strings.HasPrefix(baseURL, "https://") {
		// 启用 TLS 配置
		client.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true, // 跳过证书校验（生产环境建议设为 false 并提供 CA）
		})
		log.Println("[INFO] HTTPS detected, TLS config enabled.")
	} else {
		log.Println("[INFO] HTTP detected, no TLS config.")
	}

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
		LoginName:  "admin", // 可以根据需要调整
	}

	log.Printf("[REQUEST] Task creation request: %+v\n", reqBody)

	// 发起 POST 请求
	resp, err := n.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "*/*").
		SetBody(reqBody).
		Post(n.BaseURL + n.EndpointPath)

	// 请求错误处理
	if err != nil {
		log.Printf("[ERROR] Request failed: %v\n", err)
		return nil, nil, fmt.Errorf("请求失败: %w", err)
	}

	// HTTP 错误处理
	if resp.IsError() {
		log.Printf("[ERROR] HTTP error: %s\n", resp.Status())
		return nil, nil, fmt.Errorf("HTTP错误: %s", resp.Status())
	}

	log.Printf("[RESPONSE] Task creation response : %s\n", string(resp.Body()))

	var respBody ArcTaskResponse
	if err := json.Unmarshal(resp.Body(), &respBody); err != nil {
		log.Printf("[ERROR] Failed to parse 'respBody' field: %v\n", err)
		return nil, nil, fmt.Errorf("解析 respBody 字段失败: %w", err)
	}

	// 如果响应数据包含 Data 字段，则解析为 ArcTaskData
	if respBody.Data != "" {
		var data ArcTaskData
		if err := json.Unmarshal([]byte(respBody.Data), &data); err != nil {
			log.Printf("[ERROR] Failed to parse 'data' field: %v\n", err)
			return &respBody, nil, fmt.Errorf("解析 data 字段失败: %w", err)
		}
		log.Printf("[RESPONSE] Task creation response: %+v\n", data)
		return &respBody, &data, nil
	}

	return &respBody, nil, nil
}

// runArcTask 封装主调用逻辑
func runArcTask(baseUrl, endpointPath, path, ferryLevel string) error {
	notifier := NewArcTaskNotifier(baseUrl, endpointPath)
	path = filepath.Clean(filepath.Join("/", path))
	resp, data, err := notifier.NotifyFileTransfer(path, ferryLevel)
	if err != nil {
		log.Printf("[runArcTaskNotify] 创建任务失败: %v", err)
		return fmt.Errorf("创建任务失败: %w", err)
	}

	if resp != nil && resp.Code == "0" {
		if data != nil {
			// 输出响应信息
			log.Printf("[runArcTaskNotify] 任务创建成功！响应消息: %s, 任务ID: %s, 设备ID: %s, 错误路径: %s",
				resp.Msg, data.TaskID, data.DevID, data.ErrorPaths)
		} else {
			log.Printf("[runArcTaskNotify] 任务创建成功！响应消息: %s, 编码: %s, 状态: %s",
				resp.Msg, resp.Code, resp.State)
		}
	}
	return nil
}
