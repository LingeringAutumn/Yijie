package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"

	"github.com/LingeringAutumn/Yijie/app/video/domain/model"
	"github.com/LingeringAutumn/Yijie/app/video/infrastructure/kafka"
	"github.com/LingeringAutumn/Yijie/config"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

func (uc *videoUseCase) SubmitVideo(ctx context.Context, video *model.Video, videoData []byte) (videoId int64, videoUrl string, err error) {
	// 1. 获取当前用户 ID
	uid, err := uc.svc.GetUserId(ctx)
	if err != nil {
		return 0, "", fmt.Errorf("get user id failed: %w", err)
	}
	video.UserID = uid

	// 2. 生成视频 ID
	videoId, err = uc.svc.GenerateVideoId()
	if err != nil {
		return 0, "", fmt.Errorf("generate video id failed: %w", err)
	}
	video.VideoID = videoId

	// 3. 上传视频文件至 MinIO
	objectKey := fmt.Sprintf("%d.mp4", videoId)
	err = utils.MinioClientGlobal.UploadFile(
		constants.VideoBucket,
		objectKey,
		"us-east-1",
		"video/mp4",
		videoData,
	)
	if err != nil {
		return 0, "", fmt.Errorf("upload to MinIO failed: %w", err)
	}

	// 4. 构建视频访问 URL
	videoUrl = fmt.Sprintf("%s/%s/%s", config.Minio.Endpoint, constants.VideoBucket, objectKey)
	video.VideoURL = videoUrl

	// 5. 构造 Kafka 配置
	kafkaCfg, err := kafka.NewProducerConfig()
	if err != nil {
		return 0, "", fmt.Errorf("build kafka config failed: %w", err)
	}
	producer, err := sarama.NewSyncProducer([]string{config.Kafka.Broker}, kafkaCfg)
	if err != nil {
		return 0, "", fmt.Errorf("create kafka producer failed: %w", err)
	}
	defer producer.Close()

	// 6. 构建 Kafka 消息
	task := map[string]interface{}{
		"video_id": videoId,
		"user_id":  uid,
		"object":   objectKey,
	}
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return 0, "", fmt.Errorf("marshal kafka task failed: %w", err)
	}
	err = producer.SendMessages([]*sarama.ProducerMessage{
		{
			Topic: config.Kafka.Topic,
			Value: sarama.ByteEncoder(taskBytes),
		},
	})
	if err != nil {
		return 0, "", fmt.Errorf("send kafka message failed: %w", err)
	}

	// 7. 存入数据库
	if err = uc.svc.StoreVideo(ctx, video); err != nil {
		return 0, "", fmt.Errorf("store video meta failed: %w", err)
	}

	// 8. 初始化 video_stats 数据（防止后续查询不到）
	stat := &model.VideoStat{
		VideoID:  videoId,
		Views:    0,
		Likes:    0,
		Comments: 0,
		HotScore: 0,
	}
	if err := uc.svc.StoreVideoStats(ctx, stat); err != nil {
		return 0, "", fmt.Errorf("store video stats failed: %w", err)
	}

	return videoId, videoUrl, nil
}

func (uc *videoUseCase) GetVideo(ctx context.Context, videoId int64) (*model.VideoProfile, error) {
	// 1. 优先从 Redis 中获取缓存数据
	videoProfile, err := uc.svc.GetVideoRedis(ctx, videoId)
	if err == nil && videoProfile != nil {
		return videoProfile, nil // 命中缓存
	}

	// 2. 缓存未命中，从数据库查询
	log.Printf("video id:%d", videoId)
	videoProfile, err = uc.svc.GetVideoDB(ctx, videoId)
	if err != nil {
		return nil, fmt.Errorf("get video from db failed: %w", err)
	}

	// 3. 回写 Redis 缓存（设置过期时间）
	if err := uc.svc.SetVideoRedis(ctx, videoProfile); err != nil {
		log.Printf("warning: failed to set video cache for id %d: %v", videoId, err)
	}

	return videoProfile, nil
}

func (uc *videoUseCase) SearchVideo(ctx context.Context, keyword string, tags []string, pageNum int64, pageSize int64) ([]*model.VideoProfile, error) {
	videoProfile, err := uc.svc.SearchVideo(ctx, keyword, tags, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return videoProfile, nil
}

func (uc *videoUseCase) TrendVideo(ctx context.Context, pageNum int64, pageSize int64) ([]*model.VideoProfile, error) {
	videoProfile, err := uc.svc.TrendVideo(ctx, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return videoProfile, nil
}
