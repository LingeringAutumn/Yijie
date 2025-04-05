namespace go user

include "model.thrift"

//  注册
struct RegisterRequest {
    1: required string username,
    2: required string password,
    3: required string email,
    4: required string phone,
}

struct RegisterResponse {
    1: required model.BaseResp base,
    2: required i64 userID,
}

// 登陆
struct LoginRequest {
    1: string username,
    2: string password,
    3: string confirm_password,
}

struct LoginResponse {
    1: model.BaseResp base,
    2: model.UserInfo user,
}

// 更新个人信息
struct UpdateUserProfileRequest{
    1: required i64 uid,
    2: required model.UserProfileReq userProfileReq,
}

struct UpdateUserProfileResponse{
    1: model.BaseResp base,
    2:required model.UserProfileResp userProfileResp,
}

// 获取个人信息
struct GetUserProfileRequest{
    1: required i64 uid,
}

struct GetUserProfileResponse{
    1: model.BaseResp base,
    2: required model.UserProfileResp userProfileResp,
}

service UserService {
    RegisterResponse Register(1: RegisterRequest req),
    LoginResponse Login(1: LoginRequest req),
    UpdateUserProfileResponse UpdateUserProfile(1:UpdateUserProfileRequest req)
    GetUserProfileResponse GetUserProfile(1:GetUserProfileRequest req)
}
