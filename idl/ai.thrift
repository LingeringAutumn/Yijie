namespace go ai

struct SceneInfo {
    1: required string scene_name,                // 场景标识，这里的场景是小场景，如”浴室“”客厅“这样的，没有特殊场景就用default
    2: required string matched_component,         // 匹配的组件名
    3: optional string layout_fragment            // 该场景对应的局部 assemble JSON（可拼装）
}

struct FirstAIChatRequest{
    1: required string input_text,        // 用户输入的文字（无论是语音还是文本，转为文字）
    2: optional string language,          // 可选，语言（如 zh-CN, en-US）
    3: optional i64 timestamp             // 可选，时间戳（用于日志/追踪）
    4: required i64 uid // 用户身份唯一标识符
}


struct FirstAIChatResponse{
    1:optional string scene // 返回“智能家居”
}


struct AIChatRequest {
    1: required string input_text,        // 用户输入的文字（无论是语音还是文本，转为文字）
    2: optional string language,          // 可选，语言（如 zh-CN, en-US）
    3: optional i64 timestamp             // 可选，时间戳（用于日志/追踪）
    4: required i64 uid // 用户身份唯一标识符
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
    AIChatResponse AIChat(1: AIChatRequest req)
    // 实际的在智能家居中的聊天功能
    AIChatResponse FirstAIChat(1: FirstAIChatRequest req)
}
