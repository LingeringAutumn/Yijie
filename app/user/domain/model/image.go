package model

type Image struct {
	Uid     int64  `json:"uid"`
	ImageID int64  `json:"image_id"`
	Url     string `json:"url"`
}
