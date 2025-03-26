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

struct UserProfile {
    1: string username,
    2: string email,
    3: string phone,
    4: binary avatar,
    5: string bio,
    6: i64 membershipLevel,
    7: i64 point,
    8: string team,
}
