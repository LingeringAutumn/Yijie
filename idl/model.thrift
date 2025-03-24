namespace go model

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