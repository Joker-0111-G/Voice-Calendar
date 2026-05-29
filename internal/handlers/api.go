package handlers

import (
	"net/http"
	"time"
	"voice-calendar/internal/models"

	"github.com/gin-gonic/gin"
)

// HandleCommand 统一处理指令 (支持纯文本和语音文件扩展)
func HandleCommand(c *gin.Context) {

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
