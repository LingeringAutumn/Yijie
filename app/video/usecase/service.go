package usecase

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/LingeringAutumn/Yijie/app/video/domain/model"
	"github.com/LingeringAutumn/Yijie/config"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
)

func (uc *useCase) SubmitVideo(ctx context.Context, video *model.Video, videoData []byte) (videoId int64, videoUrl string, err error) {
	// 1. 验证用户身份
	uid, err := uc.svc.GetUserID(ctx)
	if err != nil {
		return 0, "", err
	}
	video.UserID = uid

	// 2. 生成视频 ID，用于命名文件和任务唯一标识
	videoId = uc.svc.GenerateVideoId()
	paymentID, err = svc.sf.NextVal()
	video.VideoID = videoId

	// 3. 构造 Kafka producer（同步发送模式）
	producer, err := sarama.NewSyncProducer([]string{config.Kafka.Broker}, nil)
	if err != nil {
		return 0, "", fmt.Errorf("failed to create kafka producer: %v", err)
	}
	defer producer.Close() // 退出前关闭 producer

	// 4. 构造消息体（上传任务）：将视频二进制数据封装成 JSON
	task := map[string]interface{}{
		"video_id": videoId, // 用于标识视频
		"user_id":  uid,
		"data":     videoData, // 视频内容（注意：太大建议换成 base64 或 URL）
	}
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return 0, "", fmt.Errorf("marshal kafka task failed: %v", err)
	}

	// 5. 创建 Kafka 消息，并发送到 Topic
	msg := &sarama.ProducerMessage{
		Topic: config.Kafka.Topic, // 你的 Kafka topic，比如 "video-upload"
		Value: sarama.ByteEncoder(taskBytes),
	}
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return 0, "", fmt.Errorf("send kafka message failed: %v", err)
	}

	// 6. 构造文件 URL 给前端预览（MinIO 地址拼接）
	videoUrl = fmt.Sprintf("%s/%s/%d", config.Minio.Addr, constants.VideoBucket, videoId)
	video.VideoURL = videoUrl

	// 7. 存储视频元数据到数据库（title、uid、videoId、videoURL）
	err = uc.db.StoreVideo(ctx, video)
	if err != nil {
		return 0, "", err
	}

	return videoId, videoUrl, nil
}

func (uc *useCase) GetVideo(ctx context.Context, videoId int64) (*model.VideoProfile, error) {

}
func (uc *useCase) SearchVideo(ctx context.Context, keyword string, tags []string, pageNum int64, pageSize int64) ([]*model.VideoProfile, error) {

}
func (uc *useCase) TrendVideo(ctx context.Context, pageNum int64, pageSize int64) ([]*model.VideoProfile, error) {

}
