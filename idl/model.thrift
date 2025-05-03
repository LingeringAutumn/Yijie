namespace go model

// 用户相关
struct BaseResp {
    1: i64 code,
    2: string msg,
}

struct UserInfo {
    1: i64 userId,
    2: string name,
}

struct LoginData {
    1: i64 userId,
}

// 返回的头像avatar是头像文件的url
struct UserProfileResp {
    1: string username,
    2: string email,
    3: string phone,
    4: string avatar,
    5: string bio,
    6: i64 membershipLevel,
    7: i64 point,
    8: string team,
}

struct UserProfileReq {
    1: string username,
    2: string email,
    3: string phone,
    4: string bio,
}

struct Image{
    1:required i64 imageId,
    2:required string imageUrl,
}

// 视频
struct Video {
    1: i64 video_id,              // 视频ID
    2: i64 user_id,                  // 用户ID
    3: string title,                 // 视频标题
    4: string description,           // 视频描述
    5: string cover_url,             // 封面图地址
    6: string video_url,             // 视频播放URL
    7: i64 duration_seconds,         // 视频时长
    8: i64 views,                    // 播放次数
    9: i64 likes,                    // 点赞次数
    10: i64 comments,                // 评论次数
    11: double hot_score,            // 热度分
    12: i64 created_at            // 发布时间
}


