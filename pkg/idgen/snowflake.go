package idgen

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// 开始时间戳 (当前时间)
	// 原来的开始时间戳，考虑将来需要与旧系统对齐迁移的问题 (2020-01-01): 1577808000000
	// twepoch uint64 = 1710413659000 // 2025-03-14 17:34:19
	// 开始时间戳 (2020-01-01): 1577808000000
	twepoch uint64 = 1577808000000

	// 数据中心ID所占的位数 (9位，最多512个数据中心)
	dataCenterIdBits uint64 = 9

	// 序列在ID中占的位数 (12位，每毫秒最多4096个序列号)
	sequenceBits uint64 = 12

	// 数据中心ID向左移位数
	dataCenterIdShift uint64 = sequenceBits

	// 时间戳向左移位数 (43位时间戳，实际使用42位，最高位保持为0)
	timestampLeftShift uint64 = sequenceBits + dataCenterIdBits

	// 生成序列的掩码
	sequenceMask uint64 = (1 << sequenceBits) - 1

	// 支持的最大数据中心ID
	maxDataCenterId uint64 = (1 << dataCenterIdBits) - 1
)

// SnowflakeIDGenerator 雪花算法ID生成器
type SnowflakeIDGenerator struct {
	mutex         sync.Mutex
	lastTimestamp uint64
	dataCenterId  int64
	sequence      uint64
}

// NewSnowflakeIDGenerator 创建一个新的雪花算法ID生成器
func NewSnowflakeIDGenerator() (*SnowflakeIDGenerator, error) {
	dataCenterId, err := getDataCenterId(maxDataCenterId)
	if err != nil {
		return nil, err
	}

	log.Printf("snowflake is init completed, my dataCenterId is: %d", dataCenterId)
	return &SnowflakeIDGenerator{
		lastTimestamp: 0,
		dataCenterId:  dataCenterId,
		sequence:      0,
	}, nil
}

// NextID 获取下一个ID (线程安全)
// 返回int64范围内的ID，确保与PG兼容
func (s *SnowflakeIDGenerator) NextID() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := timeGen()

	// 如果当前时间小于上一次ID生成的时间戳，说明系统时钟回退过
	if timestamp < s.lastTimestamp {
		return 0, errors.New("clock moved backwards, refusing to generate id")
	}

	// 如果是同一时间生成的，则进行毫秒内序列
	if s.lastTimestamp == timestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		// 毫秒内序列溢出
		if s.sequence == 0 {
			// 阻塞到下一个毫秒，获得新的时间戳
			timestamp = tilNextMillis(s.lastTimestamp)
		}
	} else {
		// 时间戳改变，毫秒内序列重置
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	// 移位并通过或运算拼到一起组成64位的ID
	// 确保最高位为0，使其在int64范围内
	uintId := ((timestamp - twepoch) << timestampLeftShift) |
		(uint64(s.dataCenterId) << dataCenterIdShift) |
		s.sequence

	return int64(uintId & 0x7FFFFFFFFFFFFFFF), nil
}

// tilNextMillis 阻塞到下一个毫秒，直到获得新的时间戳
func tilNextMillis(lastTimestamp uint64) uint64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

// timeGen 返回以毫秒为单位的当前时间
func timeGen() uint64 {
	return uint64(time.Now().UnixNano() / 1e6)
}

// getDataCenterId 获取数据中心ID
func getDataCenterId(maxDataCenterId uint64) (int64, error) {
	// 尝试从IP地址获取
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return int64(getRandomDataCenterId(maxDataCenterId)), nil
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipStr := ipnet.IP.String()
				if ipStr != "" && ipStr != "127.0.0.1" && ipStr != "0.0.0.0" && ipStr != "localhost" {
					log.Printf("snowflake use DataCenterId from ip, value is: %s", ipStr)
					ipParts := strings.Split(ipStr, ".")
					if len(ipParts) == 4 {
						// 使用IP的后两段，但进行hash压缩到9位
						part3, _ := strconv.Atoi(ipParts[2])
						part4, _ := strconv.Atoi(ipParts[3])
						// 简单的hash算法：(part3 * 256 + part4) % 512
						id := (part3*256 + part4) % int(maxDataCenterId+1)
						return int64(id), nil
					}
				}
			}
		}
	}

	// 如果无法从IP获取，则使用随机数
	return int64(getRandomDataCenterId(maxDataCenterId)), nil
}

// getRandomDataCenterId 获取随机数据中心ID
func getRandomDataCenterId(maxDataCenterId uint64) uint64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomId := uint64(r.Int63n(int64(maxDataCenterId + 1)))
	log.Printf("snowflake use DataCenterId from random, value is: %d", randomId)
	return randomId
}
