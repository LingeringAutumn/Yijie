## 1. IDL

下面是我们的idl代码，idl的就是定义接口的数据类型的文件

这个东西放这里主要是说如果你要问豆包的话就有一个参考，实际没什么用

这个是model,定义了一些user里面用到的结构体
```thrift
# idl/model.thrift
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

// 传进来的头像avatar是二进制文件本身
struct UserProfileReq {
    1: string username,
    2: string email,
    3: string phone,
    4: binary avatar,
    5: string bio,
}

struct Image{
    1:required i64 imageId,
    2:required string imageUrl,
}

```

user模块

```thrift
# idl/api/user.thrift
namespace go api.user
include "../model.thrift"

// 注册
struct RegisterRequest {
    1: required string name
    2: required string password
    3: required string email
    4: required string phone,
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

```

ai模块的
```thrift
# idl/api/ai.thrift
namespace go api.ai

struct SceneInfo {
    1: required string scene_name,                // 场景标识，这里的场景是小场景，如”浴室“”客厅“这样的，没有特殊场景就用default
    2: required string matched_component,         // 匹配的组件名
    3: optional string layout_fragment            // 该场景对应的局部 assemble JSON（可拼装）
}

struct FirstAIChatRequest{
    1: required string input_text,        // 用户输入的文字（无论是语音还是文本，转为文字）
    2: optional string language,          // 可选，语言（如 zh-CN, en-US）
    3: optional i64 timestamp             // 可选，时间戳（用于日志/追踪）
}

struct FirstAIChatResponse{
    1:optional string scene // 返回“智能家居”
}

struct AIChatRequest {
    1: required string input_text,        // 用户输入的文字（无论是语音还是文本，转为文字）
    2: optional string language,          // 可选，语言（如 zh-CN, en-US）
    3: optional i64 timestamp             // 可选，时间戳（用于日志/追踪）
}

struct AIChatResponse {
    1: required string reply_text,                // 系统基础回复，例如“好的，请稍等”
    // 这里为什么用list,因为我们可能要根据用户的个人模型来一次性返回多个场景，例如，用户总是在开完空调之后开电视
    // 那我们就要在开完空调之后提供给电视的组件，组装为list一次性返回
    2: required list<SceneInfo> scenes,           // 多个场景及其详细信息
    3: required string assemble_layout            // 最终完整布局 JSON，没有的话就返回一个默认的字符串，比如"default_assemble_layout"
}

service ChatService {
    // 在“翌界”界面发送第一次请求：“智能家居”
    AIChatResponse AIChat(1: AIChatRequest req)(api.post="api/v1/ai/chat/default")
    // 实际的在智能家居中的聊天功能
    AIChatResponse FirstAIChat(1: AIChatRequest req)(api.post="api/v1/ai/chat/first")
}

```

# 2. 具体接口

**注意，所有的url前面的端口都是http://127.0.0.1:20001，也就是说，比如注册端口**
**我给你的url是/api/v1/user/register**
**实际url则是二者的组合，即http://127.0.0.1:20001/api/v1/user/register**

**然后，我不知道对于前端来说，返回值是否需要提前全部定义好。返回值基本上就是一个状态吗（200,401,402,404,500这样的）加一个message+一个封装里实际数据data**

比如下面这个

**响应数据**
```
{
    "code": "10000", // 状态码，所有的接口都有
    "message": "Success", // 响应信息，所有的接口都有
    "data": {  // 我们需要的实际数据（封装好的），这个部分不同的接口不一样
        "user": {
            "userId": 1,
            "name": "yijie"
        }
    }
}
```

## 1. 用户

### 1. 注册
- **服务名称**：UserService.Register
- **URL**：http://127.0.0.1:20001/api/v1/user/register
- **HTTP请求类型**：POST
- **传输数据及格式**：
    - **请求数据**：
      ```json
      {
          "name": "string",
          "password": "string",
          "email": "string",
          "phone": "string"
      }
      ```
    - **响应数据**：
      ```json
      {
          "uid": 0
      }
      ```

**我的实际的参数**

**请求数据**
````
{
    "name":"yijie"
    "password":"yijiepwd"
    "email":"yijie@qq.com"
    "phone:18105992997"
}

````

**响应数据**

```
{
    "code": "10000",
    "message": "Success",
    "data": {
        "uid": 2
    }
}
```



### 2. 登陆
- **服务名称**：UserService.Login
- **URL**：http://127.0.0.1:20001/api/v1/user/login
- **HTTP请求类型**：POST
  - **传输数据及格式**：
      - **请求数据**：
        ```json
        {
            "name": "yijie",
            "password": "yijiepwd"
        }
        ```
        实际相应数据是底下这样的
        - **响应数据**：
          ```json
          {
            "code": "10000",
            "message": "Success",
                "data": {
                    "user": {
                    "userId": 1,
                    "name": "yijie"
                }
            }
          }
        ```

**注意鉴权的问题，我们的jwt令牌（token）是从Response Header 里的Access-Token里提取出来的，在登陆完成后，后续的所有操作，你都要在实际请求**
**的header里的Authorization里把这个jwt令牌放进去，我们必须通过这个进行用户身份的识别。**

**个人信息部分的响应数据不用管，先别做，先做ai**

### 3. 更新个人信息
- **服务名称**：UserService.UpdateProfile
- **URL**：http://127.0.0.1:20001/api/v1/user/profile/update
- **HTTP请求类型**：PUT
- **传输数据及格式**：
    - **请求数据**：
      ```json
      {
          "uid": 0,
          "userProfileReq": {
              "username": "string",
              "email": "string",
              "phone": "string",
              "avatar": "binary",
              "bio": "string"
          }
      }
      ```
    - **响应数据**：

      ```

### 4. 获取个人信息
- **服务名称**：UserService.GetProfile
- **URL**：http://127.0.0.1:20001/api/v1/user/profile/get
- **HTTP请求类型**：GET
- **传输数据及格式**：
    - **请求数据**：
      ```json
      {
          "uid": 0
      }
      ```
    - **响应数据**：
      ```json

      ```

## 2. AI

### 1. 实际的在智能家居中的聊天功能
- **服务名称**：ChatService.AIChat
- **URL**：http://127.0.0.1:20001/api/v1/ai/chat/default
- **HTTP请求类型**：POST
- **传输数据及格式**：
    - **请求数据**：
      ```json
      {
          "input_text": "string",
          "language": "string",
          "timestamp": 0
      }
      ```
      - **响应数据**：
      ```json
            {
                "code": "10000",
                "message": "Success",
                "data":{
                    "reply_text": "string",
                    "scenes": [
                        {
                            "scene_name": "string",
                            "matched_component": "string",
                            "layout_fragment": "string"
                        }
                    ],
                    "assemble_layout": "string"
                }
            }
      
        ```

### 2. 在“翌界”界面发送第一次请求
- **服务名称**：ChatService.FirstAIChat
- **URL**：http://127.0.0.1:20001/api/v1/ai/chat/first
- **HTTP请求类型**：POST
- **传输数据及格式**：
    - **请求数据**：
      ```json
      {
          "input_text": "string",
          "language": "string",
          "timestamp": 0
      }
      ```
    - **响应数据**：
      ```json
      {
        "code": "10000",
        "message": "Success",
        "data":{
          "scene":"string",
        },
      }
      ``` 