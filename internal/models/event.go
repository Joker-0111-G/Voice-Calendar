package models

import (
	"time"
	"gorm.io/gorm"
)

// Event 语音日历事件表
type Event struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	RawCommand string         `gorm:"type:text;not null" json:"raw_command"`
	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
	StartTime  time.Time      `gorm:"type:datetime;index;not null" json:"start_time"`
	EndTime    *time.Time     `gorm:"type:datetime;default:null" json:"end_time"`
	IsAllDay   bool           `gorm:"type:tinyint(1);default:0" json:"is_all_day"`
	Status     int8           `gorm:"type:tinyint;default:0" json:"status"` // 0:待办, 1:完成, 2:删除
	CreatedAt  time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}