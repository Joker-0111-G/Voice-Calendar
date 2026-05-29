package routers

import (
	"voice-calendar/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.StaticFile("/", "./static/index.html")

	api := r.Group("/api/v1")
	{
		//统一处理指令
		api.POST("/command", handlers.HandleCommand)
		// 根据选定日期查询日程
		api.GET("/events", handlers.GetEventsByDate)
		// 查询近三天汇总
		api.GET("/events/upcoming", handlers.GetUpcomingEvents)
	}

	return r
}
