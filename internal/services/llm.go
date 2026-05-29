package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Intent 结构体：映射大模型返回的 JSON
type Intent struct {
	Action   string `json:"action"`     // add, delete, query
	Title    string `json:"title"`      // 事件标题
	Time     string `json:"time"`       // 格式 YYYY-MM-DD HH:MM:SS
	IsAllDay bool   `json:"is_all_day"` // 是否全天
}

// ParseIntent 调用 DeepSeek API 提取自然语言意图
func ParseIntent(text string) (*Intent, error) {
	apiKey := "******"  //不方便透漏
	apiUrl := "https://api.deepseek.com/chat/completions"

	// Prompt
	systemPrompt := `你是一个日历助手。当前时间是 ` + time.Now().Format("2006-01-02 15:04:05") + `。
请分析用户的输入，提取日历操作意图。
你必须且只能输出合法的 JSON 格式。
字段要求：
- "action": 必须是 "add", "delete", 或 "query"
- "title": 事件名称
- "time": 标准时间戳 "YYYY-MM-DD HH:MM:SS"。如果用户使用"明天"等相对时间，请基于当前时间准确推算。如果没有具体时分秒，默认 08:00:00。
- "is_all_day": 布尔值 true/false`

	// 构建请求体 (启用 DeepSeek 的 JSON 输出强制模式)
	requestBody := map[string]interface{}{
		"model": "deepseek-chat", // DeepSeek模型
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": text},
		},
		"temperature": 0.1, // 低温度保证输出像机器一样稳定
		"response_format": map[string]string{
			"type": "json_object", // 强制 DeepSeek 只返回 JSON，杜绝多余的废话
		},
	}

	jsonData, _ := json.Marshal(requestBody)

	// 发起 HTTP 请求
	req, _ := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 15 * time.Second} // 15 秒超时防止卡死
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 DeepSeek 失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// 拦截非 200 的报错
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DeepSeek API 报错 (状态码 %d): %s", resp.StatusCode, string(body))
	}
	
	// 解析 DeepSeek 返回的原始结构
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 API 响应失败: %v", err)
	}

	// 提取内容
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("DeepSeek 没有返回有效的 choices")
	}
	
	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	content := message["content"].(string)

	// 清理可能残留的 Markdown 标记 
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	//将提取出的纯 JSON 字符串映射到 Intent 结构体
	var intent Intent
	err = json.Unmarshal([]byte(content), &intent)
	if err != nil {
		return nil, fmt.Errorf("JSON 映射到结构体失败: %v, 大模型原文: %s", err, content)
	}

	return &intent, nil
}