namespace go api.user
include "../model.thrift"

// 注册
struct RegisterRequest {
    1: required string name
    2: required string password
    3: required string email
}

struct RegisterResponse {
    1: required i64 uid;
}

// 登陆
struct LoginRequest {
    1: required string name
    2: required string password
}

struct LoginResponse {
    1: model.UserInfo user,
}

// 更新个人信息
struct UpdateUserProfileRequest{
    1: required i64 uid,
    2: required model.UserProfileReq userProfileReq,
}

struct UpdateUserProfileResponse{
    1:required model.UserProfileResp userProfileResp,
}

// 获取个人信息
struct GetUserProfileRequest{
    1: required i64 uid,
}

struct GetUserProfileResponse{
    1:required model.UserProfileResp userProfileResp,
}



service UserService {
    RegisterResponse Register(1: RegisterRequest req)(api.post = "api/v1/user/register"),
    LoginResponse Login(1: LoginRequest req)(api.post = "api/v1/user/login")
    UpdateUserProfileResponse UpdateProfile(1:UpdateUserProfileRequest req)(api.put="api/v1/user/profile/update")
    GetUserProfileResponse GetProfile(1:GetUserProfileRequest req)(api.get="api/v1/user/profile/get")
}
