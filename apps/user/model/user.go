package model

import(
	"time"
	"gorm.io/gorm"
)

type User struct {
	Id       int64  `gorm:"column:id;primaryKey;autoIncrement" json:"-"`                  // 内部 ID，不对外暴露
    Uuid      int64  `gorm:"column:uid;uniqueIndex;not null" json:"uid"`                 // 对外显示的“用户号”
    Name     string `gorm:"column:name;size:64;index;not null" json:"name"`             // 增加索引方便搜索
    Password string `gorm:"column:password;size:255;not null" json:"-"`
	Email    string `gorm:"column:email;size:128;uniqueIndex;not null" json:"email"`    // 唯一索引
    Phone    string `gorm:"column:phone;size:20;uniqueIndex" json:"phone"`
	Status   int8   `gorm:"column:status;default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
func (User) TableName() string {
	return "user"
}