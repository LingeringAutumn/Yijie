package model

type Image struct {
	Uid     int64  `json:"uid"`
	ImageID string `json:"image_id"`
	Url     string `json:"url"`
}
