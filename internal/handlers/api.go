package handlers

import (
	"fmt"
	"net/http"
	"time"
	"voice-calendar/internal/models"
	"voice-calendar/internal/services"

	"github.com/gin-gonic/gin"
)

// HandleCommand 统一处理指令 (支持纯文本和语音文件扩展)
func HandleCommand(c *gin.Context) {
	commandText := c.PostForm("text")

	// 容错：如果没传文本，看看有没有传语音文件 
	if commandText == "" {
		file, err := c.FormFile("audio")
		if err == nil {
			// 接入 ASR 解析语音文件
			commandText, _ = services.AudioToText(file)
		}
	}

	if commandText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有接收到有效指令"})
		return
	}

	// LLM 解析意图
	intent, err := services.ParseIntent(commandText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "指令解析失败: " + err.Error()})
		return
	}

	// 解析时间字符串为 Go Time 
	eventTime, _ := time.ParseInLocation("2006-01-02 15:04:05", intent.Time, time.Local)
	replyMsg := ""

	// 执行对应的数据库操作
	switch intent.Action {
	case "add":
		event := models.Event{
			RawCommand: commandText,
			Title:      intent.Title,
			StartTime:  eventTime,
			IsAllDay:   intent.IsAllDay,
			Status:     0,
		}
		models.DB.Create(&event)
		replyMsg = fmt.Sprintf("已为您添加日程：%s，时间：%s", intent.Title, intent.Time)

	case "delete":
		// 基础查询，只能删除目前状态为 0 (待办) 的日程
		dbQuery := models.DB.Model(&models.Event{}).Where("status = ?", 0)

		// 大模型提取到了标题，加入模糊匹配
		if intent.Title != "" {
			dbQuery = dbQuery.Where("title LIKE ?", "%"+intent.Title+"%")
		}

		// 放宽时间的匹配范围
		dbQuery = dbQuery.Where("start_time >= ? AND start_time <= ?", eventTime.Add(-24*time.Hour), eventTime.Add(24*time.Hour))

		// 执行更新状态（软删除）
		result := dbQuery.Update("status", 2)

		// 判断受影响的行数，给用户精准的反馈
		if result.RowsAffected == 0 {
			replyMsg = "抱歉，没有找到对应的日程。请明确说出要删除的日程名称或时间，例如：删除8点的闹钟。"
		} else {
			if intent.Title != "" {
				replyMsg = fmt.Sprintf("已为您取消日程：%s", intent.Title)
			} else {
				replyMsg = "已为您取消匹配的日程"
			}
		}

	case "query":
		replyMsg = "查询请求已受理，结果见列表"
	}

	//TTS 语音播报地址
	audioUrl := services.TextToAudio(replyMsg)

	c.JSON(http.StatusOK, gin.H{
		"message":   replyMsg,
		"audio_url": audioUrl,
		"intent":    intent,
	})
}

// GetEventsByDate 获取用户在前端日历上选定日期的日程
func GetEventsByDate(c *gin.Context) {
	dateStr := c.Query("date") // 接收前端传来的 YYYY-MM-DD

	var targetDate time.Time
	if dateStr == "" {
		targetDate = time.Now()
	} else {
		var err error
		targetDate, err = time.ParseInLocation("2006-01-02", dateStr, time.Local)
		if err != nil {
			targetDate = time.Now() // 解析失败默认看今天
		}
	}

	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var events []models.Event
	models.DB.Where("status = ? AND start_time >= ? AND start_time < ?", 0, startOfDay, endOfDay).
		Order("start_time asc").
		Find(&events)

	c.JSON(http.StatusOK, gin.H{"data": events})
}

// GetUpcomingEvents 获取近三天（含今天、明天、后天）的日程汇总
func GetUpcomingEvents(c *gin.Context) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endOfThirdDay := startOfDay.Add(3 * 24 * time.Hour) // 往后推3天

	var events []models.Event
	models.DB.Where("status = ? AND start_time >= ? AND start_time < ?", 0, startOfDay, endOfThirdDay).
		Order("start_time asc").
		Find(&events)

	c.JSON(http.StatusOK, gin.H{"data": events})
}
