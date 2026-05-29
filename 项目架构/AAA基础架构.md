# Voice-Calendar

语音版的日历工具  



题目要求 

请开发一个语音日历工具，帮助用户提高日历管理的效率和便捷性。
要求：请了解用户在日历管理中的真实需求，设计并实现一个以语音交互为核心的日历管理工具，能准确、顺畅地实现通过语音添加/删除/查看事件提醒等能力。



目前来说，攻克传统日历建立相关提醒时间的麻烦 我们需要通过语音进行合理的把控，通过语音对提醒事件 进行设置、增加、删除、更改以及查询功能

介于项目的 实际性，我们目前只需要在自己本地的日历中进行相关的组件加载，不需要重新一个新的日历，只需要识别个体用户的相关举动，不用考虑不同用户之间的问题，所以此时只需要一张表即可完成当前的要求 （项目后续只需要绑定到特定的日历软件，所以此时只需要添加一张表） 

### 核心表：`events`（MySQL）

| 字段名          | MySQL 类型     | 必填 | 默认值                        | 说明与设计意图                                               |
| :-------------- | :------------- | :--- | :---------------------------- | :----------------------------------------------------------- |
| `id`            | `BIGINT`       | 是   | `AUTO_INCREMENT`              | 主键，唯一标识一条提醒/事件。                                |
| `raw_command` ⚡ | `TEXT`         | 是   | 无                            | **【核心】** 用户原始的语音转文本（如：“明天下午三点提醒我交报告”），用于查错和对账。 |
| `title`         | `VARCHAR(255)` | 是   | 无                            | LLM 提炼出的事件核心内容（如：“交报告”）。                   |
| `start_time` ⏰  | `DATETIME`     | 是   | 无                            | 提醒触发的精确时间。**必须建立索引**，这是每天查询提醒的依据。 |
| `end_time`      | `DATETIME`     | 否   | `NULL`                        | 结束时间。很多语音提醒只有单个时间点，因此允许为空。         |
| `is_all_day` ⚡  | `TINYINT(1)`   | 是   | `0`                           | **【核心】** 处理模糊时间。用户说“明天提醒我带伞”，属于全天事件，不绑定具体时分秒。 |
| `status`        | `TINYINT`      | 是   | `0`                           | 状态机：`0`=待提醒，`1`=已完成，`2`=已取消。语音删除建议软删除，将状态改为 `2`，保留数据可追溯。 |
| `created_at`    | `DATETIME`     | 是   | `CURRENT_TIMESTAMP`           | 记录创建时间。                                               |
| `updated_at`    | `DATETIME`     | 是   | `CURRENT_TIMESTAMP ON UPDATE` | 记录最后一次修改时间，支撑“语音更改事件”功能。               |

CREATE DATABASE voice_calendar DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;





其次基本框架 

 root@DESKTOP-31VCU6N:/home/gmj/homegmj/GoCode/Voice-Calendar# tree
.
├── README.md
├── cmd           main.go主程序入口
├── go.mod
├── go.sum
├── internal
│   ├── handlers     api Gin接口
│   ├── models        数据
│   └── services     业务
├── routers         路由
├── static        前端
└── 项目架构
    ├── A1表设计.md
    ├── AAA基础架构.md
    └── 图片
        └── 屏幕截图 2026-05-30 011508.png

10 directories, 6 files

构建基本框架



编写相应的路由以及接口

```
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

```



voice.go  处理语音转文本操作

```

// AudioToText 处理音频转文字 (ASR)
// 在这里调用百度语音或阿里云 ASR 的 SDK
func AudioToText(file *multipart.FileHeader) (string, error) {
	return "", fmt.Errorf("ASR 服务尚未配置，请直接传输文本指令")
}

// TextToAudio 将文字转为语音 URL (TTS)
func TextToAudio(text string) string {
	// 这里返回前端可调用的第三方免费 TTS 接口作为 MVP 替身
	return fmt.Sprintf("https://dict.youdao.com/dictvoice?audio=%s&le=zh", text)
}

```



llm.go 使用deepseek API



至此 后端所有框架基本已经全部完成，仅剩前端部分 

这里我优先考虑 结合AI进行编写 并 进行子那个赢得整合 以达到最终演示效果  （5.30 3点20）
