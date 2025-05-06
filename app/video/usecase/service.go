package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/LingeringAutumn/Yijie/app/video/domain/model"
	"github.com/LingeringAutumn/Yijie/app/video/infrastructure/kafka"
	"github.com/LingeringAutumn/Yijie/config"
	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

func (uc *videoUseCase) SubmitVideo(ctx context.Context, video *model.Video, videoData []byte) (videoId int64, videoUrl string, err error) {
	// 1. 验证用户身份（从上下文中提取当前登录用户的 UID）
	// 如果未登录或 Token 无效，GetUserId 会返回错误，直接中断后续操作。
	uid, err := uc.svc.GetUserId(ctx)
	if err != nil {
		return 0, "", fmt.Errorf("get user id failed: %w", err)
	}
	video.UserID = uid // 设置当前用户为视频的拥有者

	// 2. 为该视频生成唯一 ID，使用雪花算法
	videoId, err = uc.svc.GenerateVideoId()
	if err != nil {
		return 0, "", fmt.Errorf("generate video id failed: %w", err)
	}
	video.VideoID = videoId // 设置回 video 实体中，供后续使用

	// 3. 上传视频内容至 MinIO（对象存储），以对象名 videoId.mp4 命名。
	// 这是同步上传方式，适合中小文件.
	objectKey := fmt.Sprintf("%d.mp4", videoId)
	err = utils.MinioClientGlobal.UploadFile(
		constants.VideoBucket, // MinIO 中用于存视频的桶，例如 "video"
		objectKey,             // 文件名，例如 "16834545454.mp4"
		"us-east-1",           // MinIO 兼容 S3，因此区域可写固定值
		"video/mp4",           // MIME 类型，前端播放需要
		videoData,             // 视频的原始二进制字节流
	)
	if err != nil {
		return 0, "", fmt.Errorf("upload video to MinIO failed: %w", err)
	}

	// 4. 构造视频的可访问 URL，用于数据库存储和前端展示。
	// 格式为：http://{MinIO地址}/{桶}/{对象名}
	videoUrl = fmt.Sprintf("%s/%s/%s", config.Minio.Endpoint, constants.VideoBucket, objectKey)
	video.VideoURL = videoUrl

	// 5. 构建 Kafka 配置，用于生产者发送任务。
	// 包括启用 SASL 用户认证、设置协议版本、是否启用 TLS、ACK 级别等。
	// 封装在 kafka.NewProducerConfig() 函数中，避免主流程中出现大量配置代码。
	kafkaCfg, err := kafka.NewProducerConfig()
	if err != nil {
		return 0, "", fmt.Errorf("failed to build kafka config: %w", err)
	}

	// 6. 创建 Kafka 同步生产者（同步意味着能拿到是否成功发送的反馈）
	// Kafka 用于发送上传任务消息，供其他服务（如异步转码）消费处理。
	producer, err := sarama.NewSyncProducer([]string{config.Kafka.Broker}, kafkaCfg)
	if err != nil {
		return 0, "", fmt.Errorf("create kafka producer failed: %w", err)
	}
	defer producer.Close() // 退出函数前关闭 Kafka 连接，避免资源泄漏

	// 7. 构造上传任务的 Kafka 消息体，仅包含最小必要信息：
	// - video_id：唯一标识视频
	// - user_id：上传用户
	// - object：MinIO 中视频对象路径（不含完整 URL，防止暴露内部地址）
	// 注意：视频内容本身不通过 Kafka 传递，以免造成消息积压或传输失败
	task := map[string]interface{}{
		"video_id": videoId,
		"user_id":  uid,
		"object":   objectKey,
	}
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return 0, "", fmt.Errorf("marshal kafka task failed: %w", err)
	}

	// 8. 创建 Kafka 消息，指定要发送的 Topic 和消息内容
	// 通常 Topic 名字在 config 中定义，例如 "video-upload"
	msg := &sarama.ProducerMessage{
		Topic: config.Kafka.Topic,
		Value: sarama.ByteEncoder(taskBytes),
	}
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return 0, "", fmt.Errorf("send kafka message failed: %w", err)
	}

	// 9. 将视频元数据存入数据库，包括用户 ID、标题、描述、URL 等
	// 这样即便后续 Kafka 消费失败，前端也能通过数据库查询到上传记录
	err = uc.db.StoreVideo(ctx, video)
	if err != nil {
		return 0, "", fmt.Errorf("store video meta failed: %w", err)
	}

	// 10. 整个流程结束，返回成功结果：视频 ID 与视频访问地址
	return videoId, videoUrl, nil
}

func (uc *videoUseCase) GetVideo(ctx context.Context, videoId int64) (*model.VideoProfile, error) {

}
func (uc *videoUseCase) SearchVideo(ctx context.Context, keyword string, tags []string, pageNum int64, pageSize int64) ([]*model.VideoProfile, error) {

}
func (uc *videoUseCase) TrendVideo(ctx context.Context, pageNum int64, pageSize int64) ([]*model.VideoProfile, error) {

}
