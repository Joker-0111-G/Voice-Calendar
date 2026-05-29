package main

import (
	"log"
	"voice-calendar/internal/models"
	"voice-calendar/routers"
)

func main() {
	models.InitDB()
	r := routers.SetupRouter()

	log.Println("语音日历服务启动！请在浏览器访问 http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
