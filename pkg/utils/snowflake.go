package utils

import (
	"fmt"
	"sync"
	"time"
)

const (
	epoch             int64 = 1577808000000 // 设置起始时间(时间戳/毫秒)：2020-01-01 00:00:00，有效期69年
	timestampBits     int64 = 41            // 时间戳占用位数
	datacenteridBits  int64 = 5             // 数据中心id所占位数
	workeridBits      int64 = 5             // 机器id所占位数
	sequenceBits      int64 = 12            // 序列所占的位数
	timestampMax      int64 = 2199023255551 // (-1 ^ (-1 << timestampBits))时间戳最大值
	datacenteridMax   int64 = 31            // (-1 ^ (-1 << datacenteridBits))支持的最大数据中心id数量
	workeridMax       int64 = 31            // (-1 ^ (-1 << workeridBits)) 支持的最大机器id数量
	sequenceMask      int64 = 4095          // (-1 ^ (-1 << sequenceBits))支持的最大序列id数量
	workeridShift     int64 = 12            // sequenceBits 机器id左移位数
	datacenteridShift int64 = 17            // sequenceBits + workeridBits 数据中心id左移位数
	timestampShift    int64 = 22            // sequenceBits + workeridBits + datacenteridBits 时间戳左移位数

	NanosecondsInAMillisecond = 1_000_000 // 每毫秒的纳秒数
	MillisecondsInASecond     = 1000      // 每秒的毫秒数
)

// Snowflake 结构体，用于生成唯一 ID
type Snowflake struct {
	// 互斥锁，确保并发安全
	sync.Mutex
	// 上一次生成 ID 的时间戳
	timestamp int64
	// 机器 ID
	workerid int64
	// 数据中心 ID
	datacenterid int64
	// 序列号
	sequence int64
}

// NewSnowflake 创建一个新的 Snowflake 实例
func NewSnowflake(datacenterid, workerid int64) (*Snowflake, error) {
	// 检查数据中心 ID 是否在有效范围内
	if datacenterid < 0 || datacenterid > datacenteridMax {
		return nil, fmt.Errorf("datacenterid must be between 0 and %d", datacenteridMax-1)
	}
	// 检查机器 ID 是否在有效范围内
	if workerid < 0 || workerid > workeridMax {
		return nil, fmt.Errorf("workerid must be between 0 and %d", workeridMax-1)
	}
	// 创建并返回 Snowflake 实例
	return &Snowflake{
		timestamp:    0,
		datacenterid: datacenterid,
		workerid:     workerid,
		sequence:     0,
	}, nil
}

// NextVal 生成下一个唯一 ID
// 生成规则：timestamp + 数据中心 id + 工作节点 id + 自旋 id
func (s *Snowflake) NextVal() (int64, error) {
	// 加锁，确保并发安全
	s.Lock()
	// 获取当前时间戳（毫秒）
	now := time.Now().UnixNano() / NanosecondsInAMillisecond
	// 如果当前时间戳和上一次生成 ID 的时间戳相同
	if s.timestamp == now {
		// 序列号加 1，并通过 sequenceMask 取模
		s.sequence = (s.sequence + 1) & sequenceMask
		// 如果序列号超出 12bit 长度
		if s.sequence == 0 {
			// 等待下一毫秒
			for now <= s.timestamp {
				now = time.Now().UnixNano() / NanosecondsInAMillisecond
			}
		}
	} else {
		// 不同时间戳下，序列号重置为 0
		s.sequence = 0
	}
	// 计算当前时间戳与起始时间的差值
	t := now - epoch
	// 检查时间戳是否超出最大值
	if t > timestampMax {
		s.Unlock()
		return 0, fmt.Errorf("epoch must be between 0 and %d", timestampMax-1)
	}
	// 更新上一次生成 ID 的时间戳
	s.timestamp = now
	// 通过位运算组合时间戳、数据中心 ID、机器 ID 和序列号生成最终 ID
	r := (t)<<timestampShift | (s.datacenterid << datacenteridShift) | (s.workerid << workeridShift) | (s.sequence)
	// 解锁
	s.Unlock()
	return r, nil
}

// GetDeviceID 获取数据中心 ID 和机器 ID
func GetDeviceID(sid int64) (datacenterid, workerid int64) {
	// 通过右移和按位与运算获取数据中心 ID
	datacenterid = (sid >> datacenteridShift) & datacenteridMax
	// 通过右移和按位与运算获取机器 ID
	workerid = (sid >> workeridShift) & workeridMax
	return
}

// GetTimestamp 获取时间戳
func GetTimestamp(sid int64) (timestamp int64) {
	// 通过右移和按位与运算获取时间戳
	timestamp = (sid >> timestampShift) & timestampMax
	return
}

// GetGenTimestamp 获取创建 ID 时的时间戳
func GetGenTimestamp(sid int64) (timestamp int64) {
	// 获取时间戳并加上起始时间
	timestamp = GetTimestamp(sid) + epoch
	return
}

// GetGenTime 获取创建 ID 时的时间字符串(精度：秒)
func GetGenTime(sid int64) (t string) {
	// 需将 GetGenTimestamp 获取的时间戳 / 1000 转换成秒
	t = time.Unix(GetGenTimestamp(sid)/MillisecondsInASecond, 0).Format("2006-01-02 15:04:05")
	return
}

// GetTimestampStatus 获取时间戳已使用的占比：范围（0.0 - 1.0）
func GetTimestampStatus() (state float64) {
	// 计算当前时间戳与起始时间的差值，并除以时间戳最大值
	state = float64((time.Now().UnixNano()/NanosecondsInAMillisecond - epoch)) / float64(timestampMax)
	return
}
