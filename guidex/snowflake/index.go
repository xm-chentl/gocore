package snowflake

import (
	"strconv"
	"sync"
	"time"

	"github.com/xm-chentl/gocore/guidex"

	"github.com/golang/glog"
)

const (
	epoch             = int64(1577808000000)                           // 设置起始时间(时间戳/毫秒)：2020-01-01 00:00:00，有效期69年
	timestampBits     = uint(41)                                       // 时间戳占用位数
	dataCenterIDBits  = uint(2)                                        // 数据中心id所占位数
	workerIDBits      = uint(7)                                        // 机器id所占位数
	sequenceBits      = uint(12)                                       // 序列所占的位数
	timestampMax      = int64(-1 ^ (-1 << timestampBits))              // 时间戳最大值
	dataCenterIDMax   = int64(-1 ^ (-1 << dataCenterIDBits))           // 支持的最大数据中心id数量
	workerIDMax       = int64(-1 ^ (-1 << workerIDBits))               // 支持的最大机器id数量
	sequenceMask      = int64(-1 ^ (-1 << sequenceBits))               // 支持的最大序列id数量
	workerIDShift     = sequenceBits                                   // 机器id左移位数
	dataCenterIDShift = sequenceBits + workerIDBits                    // 数据中心id左移位数
	timestampShift    = sequenceBits + workerIDBits + dataCenterIDBits // 时间戳左移位数
)

// 参考 https://cloud.tencent.com/developer/article/1820225
type snowflake struct {
	sync.Mutex

	dataCenterID int64 // 数据中心机房ID
	sequence     int64 // 序列号
	timestamp    int64 // 时间戳
	workerID     int64 // 工作节点
}

func (s *snowflake) String() (uuid string) {
	s.Lock()
	now := time.Now().UnixNano() / 1000000 // 转毫秒
	if s.timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出12bit长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}

	t := now - epoch
	if t > timestampMax {
		s.Unlock()
		glog.Errorf("epoch must be between 0 and %d", timestampMax-1)
		return
	}

	s.timestamp = now
	r := int64((t)<<timestampShift | (s.dataCenterID << dataCenterIDShift) | (s.workerID << workerIDShift) | (s.sequence))
	s.Unlock()
	uuid = strconv.Itoa(int(r))
	return
}

func New() guidex.IGenerate {
	return &snowflake{}
}
