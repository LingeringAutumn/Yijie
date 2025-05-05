namespace go api.video
include "../model.thrift"

/**
 * 视频投稿请求结构
 * 用于提交视频信息（视频文件二进制不直接走 RPC，而是通过 MinIO 存储后上传 URL）
 */
struct VideoSubmissionRequest {
    1: required i64 user_id                // 投稿者用户ID
    2: required string title              // 视频标题
    3: required string description        // 视频描述，可为空
    4: optional string cover_url          // 封面图URL（可选）
    5: required i64 duration_seconds      // 视频时长（单位：秒）
}

/**
 * 视频投稿响应结构
 * 返回视频ID及存储URL（可回显）
 */
struct VideoSubmissionResponse {
    1: required model.BaseResp base_resp  // 通用响应结构
    2: required i64 video_id           // 成功创建后的视频ID
    3: required string video_url          // 视频可访问URL
}

/**
 * 获取视频详情请求结构
 * 用于获取某一条视频详情数据
 */
struct VideoDetailRequest {
    1: required i64 video_id           // 视频ID
}

/**
 * 视频详情响应结构
 */
struct VideoDetailResponse {
    1: required model.BaseResp base_resp
    2: required model.Video video         // 视频详细数据
}

/**
 * 视频搜索请求结构
 * 支持关键词模糊匹配，分页获取
 */
struct VideoSearchRequest {
    1: required string keyword            // 搜索关键词（匹配标题/描述）
    2: optional list<string> tags         // 可选标签（预留扩展）
    3: required i64 page_num              // 第几页，从1开始
    4: required i64 page_size             // 每页多少条数据
}

/**
 * 视频搜索响应结构
 */
struct VideoSearchResponse {
    1: required model.BaseResp base_resp
    2: required list<model.Video> videos  // 匹配的视频列表
}

/**
 * 热榜视频请求结构
 * 按照热度评分返回热门视频
 */
struct VideoTrendingRequest {
    1: required i64 page_num              // 第几页
    2: required i64 page_size             // 每页条数
}

/**
 * 热榜视频响应结构
 */
struct VideoTrendingResponse {
    1: required model.BaseResp base_resp
    2: required list<model.Video> videos  // 热门视频列表
}

/**
 * 视频服务接口定义
 * 支持视频投稿、查询详情、关键词搜索、热榜获取等
 */
service VideoService {
    VideoSubmissionResponse SubmitVideo(1: VideoSubmissionRequest req)(api.post = "api/v1/video/submit"),
    VideoDetailResponse GetVideo(1: VideoDetailRequest req)(api.get = "api/v1/video/get"),
    VideoSearchResponse SearchVideo(1: VideoSearchRequest req)(api.get = "api/v1/video/search"),
    VideoTrendingResponse TrendingVideo(1: VideoTrendingRequest req)(api.get = "api/v1/video/trending")
}
