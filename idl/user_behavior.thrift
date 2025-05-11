namespace go user_behavior
include "/model.thrift"

/**
 * 视频点赞请求结构
 * is_like = true 表示点赞，false 表示取消点赞
 */
struct VideoLikeRequest {
    1: required i64 video_id    // 视频ID
    2: required i64 user_id     // 用户ID
    3: required bool is_like    // 是否点赞（true=点赞，false=取消点赞）
}

/**
 * 视频点赞响应结构
 */
struct VideoLikeResponse {
    1: required model.BaseResp base_resp
}

/**
 * 点赞服务接口定义
 * 目前仅支持视频点赞/取消点赞操作
 */
service LikeService {
    VideoLikeResponse LikeVideo(1: VideoLikeRequest req)
}