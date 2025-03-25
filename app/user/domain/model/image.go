package model

type Image struct {
	ImageID int64  `json:"image_id" gorm:"column:image_id;primaryKey;autoIncrement"`
	Uid     int64  `json:"uid" gorm:"column:uid;not null;index"`
	Url     string `json:"url" gorm:"column:url;type:varchar(255);not null"`
}
