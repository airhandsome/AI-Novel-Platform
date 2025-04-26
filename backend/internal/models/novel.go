package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// OutlineItem 大纲条目结构
type OutlineItem struct {
	ID          uint          `json:"id"`
	Title       string        `json:"title"`       // 章节标题
	Description string        `json:"description"` // 章节描述
	Order       int          `json:"order"`       // 排序
	Children    []OutlineItem `json:"children"`    // 子节点
}

// Character 角色结构
type Character struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Location 地点结构
type Location struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// WorldBuilding 世界观设定结构
type WorldBuilding struct {
	Background  string      `json:"background"`
	Characters  []Character `json:"characters"`
	Locations   []Location  `json:"locations"`
}

// NovelOutline 小说大纲完整结构
type NovelOutline struct {
	Outline       []OutlineItem `json:"outline"`
	WorldBuilding WorldBuilding `json:"worldBuilding"`
}

// Value 实现 driver.Valuer 接口
func (no NovelOutline) Value() (driver.Value, error) {
	return json.Marshal(no)
}

// Scan 实现 sql.Scanner 接口
func (no *NovelOutline) Scan(value interface{}) error {
	if value == nil {
		*no = NovelOutline{
			Outline: []OutlineItem{},
			WorldBuilding: WorldBuilding{
				Background: "",
				Characters: []Character{},
				Locations:  []Location{},
			},
		}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	if err := json.Unmarshal(bytes, &no); err != nil {
		return err
	}

	// 确保切片不为 nil
	if no.Outline == nil {
		no.Outline = []OutlineItem{}
	}
	if no.WorldBuilding.Characters == nil {
		no.WorldBuilding.Characters = []Character{}
	}
	if no.WorldBuilding.Locations == nil {
		no.WorldBuilding.Locations = []Location{}
	}

	return nil
}

// StringArray 自定义类型
type StringArray []string

// Value 实现 driver.Valuer 接口
func (sa StringArray) Value() (driver.Value, error) {
	return json.Marshal(sa)
}

// Scan 实现 sql.Scanner 接口
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = StringArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, sa)
}

type Novel struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	Title         string         `json:"title" gorm:"size:100;not null"`
	Description   string         `json:"description" gorm:"type:text"`
	Category      string         `json:"category" gorm:"size:50"`
	CoverURL      string         `json:"coverUrl" gorm:"size:255"`
	Status        int            `json:"status" gorm:"default:0"` // 0: 草稿, 1: 连载中, 2: 已完结
	WordCount     int            `json:"wordCount" gorm:"default:0"`
	ReadCount     int            `json:"readCount" gorm:"default:0"`
	FavoriteCount int            `json:"favoriteCount" gorm:"default:0"`
	AuthorID      uint           `json:"authorId" gorm:"not null;type:uint REFERENCES users(id)"`
	Author        User           `json:"author" gorm:"foreignKey:AuthorID"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	Tags          StringArray    `json:"tags" gorm:"type:json"`
	NovelOutline  NovelOutline   `json:"novelOutline" gorm:"type:json"` // 小说大纲，包含世界观设定
}

func (Novel) TableName() string {
	return "novels"
}