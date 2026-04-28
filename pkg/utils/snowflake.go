package utils

import (
	"errors"
	"sync"
	"time"
)

const (
	snowflakeEpoch     int64 = 1704038400000 // 2024-01-01T00:00:00Z
	nodeIDBits         uint8 = 10
	sequenceBits       uint8 = 12
	maxNodeID          int64 = -1 ^ (-1 << nodeIDBits)
	maxSequence        int64 = -1 ^ (-1 << sequenceBits)
	nodeIDShift        uint8 = sequenceBits
	timestampLeftShift uint8 = nodeIDBits + sequenceBits
)

type snowflakeGenerator struct {
	mu       sync.Mutex
	nodeID   int64
	lastMs   int64
	sequence int64
}

var defaultGenerator = &snowflakeGenerator{nodeID: 1}

func NextSnowflakeID() (int64, error) {
	return defaultGenerator.NextID()
}

func (g *snowflakeGenerator) NextID() (int64, error) {
	if g.nodeID < 0 || g.nodeID > maxNodeID {
		return 0, errors.New("snowflake node id out of range")
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	nowMs := time.Now().UnixMilli()
	if nowMs < g.lastMs {
		return 0, errors.New("clock moved backwards")
	}

	if nowMs == g.lastMs {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			for nowMs <= g.lastMs {
				nowMs = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastMs = nowMs
	id := ((nowMs - snowflakeEpoch) << timestampLeftShift) | (g.nodeID << nodeIDShift) | g.sequence
	return id, nil
}
