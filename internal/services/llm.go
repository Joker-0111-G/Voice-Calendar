package services


// Intent 结构体：映射大模型返回的 JSON
type Intent struct {
	Action   string `json:"action"`     // add, delete, query
	Title    string `json:"title"`      // 事件标题
	Time     string `json:"time"`       // 格式 YYYY-MM-DD HH:MM:SS
	IsAllDay bool   `json:"is_all_day"` // 是否全天
}

// ParseIntent 调用 DeepSeek API 提取自然语言意图
func ParseIntent(text string) (*Intent, error) {
	return nil,nil
}