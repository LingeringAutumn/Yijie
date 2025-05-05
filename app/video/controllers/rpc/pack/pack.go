package pack

import (
	dmodel "github.com/LingeringAutumn/Yijie/app/video/domain/model"
	kmodel "github.com/LingeringAutumn/Yijie/kitex_gen/model"
)

// BuildVideo 将 entities 定义的 Video 实体转换成 idl 定义的 RPC 交流实体，类似 dto
func BuildVideo(video *dmodel.VideoProfile) *kmodel.Video {
	return &kmodel.Video{
		VideoId:         video.VideoID,
		UserId:          video.UserID,
		Title:           video.Title,
		Description:     video.Description,
		CoverUrl:        video.CoverURL,
		VideoUrl:        video.VideoURL,
		DurationSeconds: video.DurationSeconds,
		Views:           video.Views,
		Likes:           video.Likes,
		Comments:        video.Comments,
		HotScore:        video.HotScore,
		CreatedAt:       video.CreatedAt,
	}
}

// BuildVideoList 构建视频列表
func BuildVideoList(videoList []*dmodel.VideoProfile) []*kmodel.Video {
	var resp []*kmodel.Video
	for _, v := range videoList {
		resp = append(resp, BuildVideo(v))
	}
	return resp
}
