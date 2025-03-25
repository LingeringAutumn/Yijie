package constants

// 第一组常量定义了默认的数据中心 ID 和工作节点 ID
// 这两个常量用于在没有明确指定数据中心 ID 和工作节点 ID 时，为 Snowflake 算法提供默认值
const (
	// DefaultDataCenterID 是默认的数据中心 ID，设置为 0
	// 在使用 Snowflake 算法生成唯一 ID 时，如果没有指定数据中心 ID，会使用这个默认值
	DefaultDataCenterID = int64(0)
	// DefaultWorkerID 是默认的工作节点 ID，设置为 0
	// 在使用 Snowflake 算法生成唯一 ID 时，如果没有指定工作节点 ID，会使用这个默认值
	DefaultWorkerID = int64(0)
)

// 第二组常量定义了不同服务对应的工作节点 ID
// 使用 iota 关键字进行自增赋值，确保每个服务的工作节点 ID 是唯一的
const (
	// WorkerOfUserService 是用户服务对应的工作节点 ID
	// 它的值从 1 开始，因为 iota 从 0 开始，这里加了 1
	WorkerOfUserService = 1 + int64(iota)
	// WorkerOfOrderService 是订单服务对应的工作节点 ID
	// 由于使用了 iota，它的值会自动递增，这里是 2
	WorkerOfOrderService
)
