package model

import (
	"time"
)

type Song struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Group       string    `json:"group" gorm:"column:group_name;not null"`
	Song        string    `json:"song" gorm:"column:song_name;not null"`
	ReleaseDate time.Time `json:"release_date" gorm:"column:release_date"`
	Text        string    `json:"text" gorm:"column:text;type:text"`
	Link        string    `json:"link" gorm:"column:youtube_link"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Song) TableName() string {
	return "songs"
}
