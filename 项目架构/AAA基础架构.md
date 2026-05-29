# Voice-Calendar

语音版的日历工具  



题目要求 

请开发一个语音日历工具，帮助用户提高日历管理的效率和便捷性。
要求：请了解用户在日历管理中的真实需求，设计并实现一个以语音交互为核心的日历管理工具，能准确、顺畅地实现通过语音添加/删除/查看事件提醒等能力。



目前来说，攻克传统日历建立相关提醒时间的麻烦 我们需要通过语音进行合理的把控，通过语音对提醒事件 进行设置、增加、删除、更改以及查询功能

介于项目的 实际性，我们目前只需要在自己本地的日历中进行相关的组件加载，不需要重新一个新的日历，只需要识别个体用户的相关举动，不用考虑不同用户之间的问题，所以此时只需要一张表即可完成当前的要求 （项目后续只需要绑定到特定的日历软件，所以此时只需要添加一张表） 详细可以查看  A1表设计.md

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

  api := r.Group("/api/v1")

  {

​    //统一处理指令

​    api.POST("/command", handlers.HandleCommand)

​    // 根据选定日期查询日程 

​    api.GET("/events", handlers.GetEventsByDate)

​    // 查询近三天汇总

​    api.GET("/events/upcoming", handlers.GetUpcomingEvents)

  }
